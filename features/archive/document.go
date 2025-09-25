package archive

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"path"
	"path/filepath"
	"strings"
	"time"
	"unterlagen/features/common"

	"github.com/h2non/filetype"
)

var (
	ErrUnsupportedFiletype = errors.New("unsupported filetype")
	ErrNotAllowed          = errors.New("not allowed")
)

const (
	PDF        Filetype = "pdf"
	Unknown    Filetype = "unknown"
	ThirtyDays          = 30 * 24 * time.Hour
)

type Filetype string

type Document struct {
	ID               string
	Title            string
	Filename         string
	Filetype         Filetype
	Filesize         uint64
	Text             string
	Summary          DocumentSummary
	PreviewFilepaths []string
	Owner            string
	FolderID         string
	TrashedAt        sql.NullTime
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func newDocument(filename string, filetype Filetype, filesize uint64, owner string, folderID string) Document {
	title := strings.TrimSuffix(filename, filepath.Ext(filename))
	return Document{
		ID:       common.GenerateID(),
		Title:    title,
		Filename: filename,
		Filetype: filetype,
		Filesize: filesize,
		Summary: DocumentSummary{
			IsGenerating: true,
		},
		Owner:    owner,
		FolderID: folderID,
		TrashedAt: sql.NullTime{
			Valid: false,
		},
	}
}

func (document Document) Name() string {
	return strings.TrimSuffix(document.Filename, filepath.Ext(document.Filename))
}

func (document Document) Filepath() string {
	return path.Join(document.Owner, document.ID, document.Filename)
}

func (document Document) PreviewPrefix() string {
	return path.Join(document.Owner, document.ID, "previews")
}

func (document Document) ShouldBeDeleted() bool {
	if !document.TrashedAt.Valid {
		return false
	}

	return time.Since(document.TrashedAt.Time) >= ThirtyDays
}

func (document Document) IsTrashed() bool {
	return document.TrashedAt.Valid
}

type DocumentSummary struct {
	Overview     string   `json:"overview"`
	KeyPoints    []string `json:"key_points"`
	IsGenerating bool     `json:"is_generating"`
}

type DocumentRepository interface {
	Save(document Document) error
	FindByID(id string) (Document, error)
	FindAllByIDIn(ids []string) ([]Document, error)
	FindAllByOwner(owner string) ([]Document, error)
	FindAllByFolderID(folderID string) ([]Document, error)
	FindAllTrashed() ([]Document, error)
	DeleteByID(id string) error
}

type DocumentConsumer func(r io.Reader) error
type DocumentStorage interface {
	Store(filepath string, r io.Reader) error
	Retrieve(filepath string, consumer DocumentConsumer) error
	Delete(filepath string) error
	Size(filepath string) (int64, error)
}

type DocumentPreviewStorage interface {
	Store(path string, r io.Reader) error
	Retrieve(path string, consumer func(r io.Reader) error) error
	Delete(preview string) error
}

type DocumentAnalyzer interface {
	GeneratePreviews(document Document) ([]string, error)
	ExtractText(document Document) (string, error)
}

type DocumentSummarizer interface {
	SummarizeText(text string) (DocumentSummary, error)
}

type DocumentMessages interface {
	PublishDocumentUpserted(document Document) error
	SubscribeDocumentUpserted(subscriber func(document Document) error) error
	PublishDocumentTextExtracted(document Document) error
	SubscribeDocumentTextExtracted(subscriber func(document Document) error) error
	PublishDocumentDeleted(document Document) error
	SubscribeDocumentDeleted(subscriber func(document Document) error) error
}

type documents struct {
	repository     DocumentRepository
	storage        DocumentStorage
	previewStorage DocumentPreviewStorage
	messages       DocumentMessages
	taskScheduler  *common.TaskScheduler
}

func (d *documents) UploadDocument(filename string, filesize uint64, folderID string, owner string, r io.Reader) error {
	document := newDocument(filename, Unknown, filesize, owner, folderID)
	err := d.storage.Store(document.Filepath(), r)
	if err != nil {
		return err
	}

	err = d.storage.Retrieve(document.Filepath(), func(r io.Reader) error {
		filetype, err := d.determineFiletype(r)
		if err != nil {
			return err
		}
		document.Filetype = filetype
		return nil
	})
	if err != nil {
		return err
	}

	err = d.repository.Save(document)
	if err != nil {
		return err
	}

	err = d.messages.PublishDocumentUpserted(document)
	if err != nil {
		return err
	}

	return d.scheduleDocumentProcessing(document)
}

func (d *documents) determineFiletype(r io.Reader) (Filetype, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	header := make([]byte, 261)
	if _, err := io.ReadFull(tee, header); err != nil {
		return "", err
	}

	kind, _ := filetype.Match(header)
	if kind.Extension == "pdf" {
		return PDF, nil
	}

	return "", ErrUnsupportedFiletype
}

func (d *documents) GetDocumentsInFolder(folderID string, owner string) ([]Document, error) {
	documents, err := d.repository.FindAllByFolderID(folderID)
	if err != nil {
		return nil, err
	}

	// Filter documents to only those owned by the user
	var userDocuments []Document
	for _, document := range documents {
		if document.Owner == owner {
			userDocuments = append(userDocuments, document)
		}
	}

	return userDocuments, nil
}

func (d *documents) GetDocument(id string, owner string) (Document, error) {
	document, err := d.repository.FindByID(id)
	if err != nil {
		return Document{}, err
	}

	if document.Owner != owner {
		return Document{}, ErrNotAllowed
	}

	return document, nil
}

func (d *documents) GetDocuments(ids []string, owner string) ([]Document, error) {
	documents, err := d.repository.FindAllByIDIn(ids)
	if err != nil {
		return nil, err
	}

	for _, document := range documents {
		if document.Owner != owner {
			return nil, ErrNotAllowed
		}
	}

	return documents, nil
}

func (d *documents) GetDocumentPreview(id string, owner string, pageNumber int, consumer func(r io.Reader) error) error {
	document, err := d.repository.FindByID(id)
	if err != nil {
		return err
	}

	if document.Owner != owner {
		return ErrNotAllowed
	}

	previewFilepath := document.PreviewFilepaths[pageNumber]
	return d.previewStorage.Retrieve(previewFilepath, consumer)
}

func (d *documents) DownloadDocument(documentID string, owner string, consumer DocumentConsumer) error {
	document, err := d.repository.FindByID(documentID)
	if err != nil {
		return err
	}

	if document.Owner != owner {
		return errors.New("unauthorized")
	}

	return d.storage.Retrieve(document.Filepath(), consumer)
}

func (d *documents) ExportAllDocuments(owner string, writer io.Writer) error {
	documents, err := d.repository.FindAllByOwner(owner)
	if err != nil {
		return err
	}

	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	for _, document := range documents {
		if document.IsTrashed() {
			continue
		}

		err := d.storage.Retrieve(document.Filepath(), func(r io.Reader) error {
			fileWriter, err := zipWriter.Create(document.Filename)
			if err != nil {
				return err
			}

			_, err = io.Copy(fileWriter, r)
			return err
		})

		if err != nil {
			slog.Error("failed to add document to zip",
				slog.String("documentID", document.ID),
				slog.String("error", err.Error()))
			continue
		}
	}

	return nil
}

func (d *documents) TrashDocument(documentID string, owner string) error {
	document, err := d.repository.FindByID(documentID)
	if err != nil {
		return err
	}

	if document.Owner != owner {
		return errors.New("unauthorized")
	}

	document.TrashedAt.Valid = true
	document.TrashedAt.Time = time.Now()
	return d.repository.Save(document)
}

func (d *documents) RestoreDocument(documentID string, owner string) error {
	document, err := d.repository.FindByID(documentID)
	if err != nil {
		return err
	}

	if document.Owner != owner {
		return errors.New("unauthorized")
	}

	document.TrashedAt.Valid = false
	return d.repository.Save(document)
}

func (d *documents) UpdateDocumentTitle(documentID string, owner string, newTitle string) error {
	document, err := d.repository.FindByID(documentID)
	if err != nil {
		return err
	}

	if document.Owner != owner {
		return errors.New("unauthorized")
	}

	document.Title = newTitle
	document.UpdatedAt = time.Now()

	err = d.repository.Save(document)
	if err != nil {
		return err
	}

	return d.messages.PublishDocumentUpserted(document)
}

func (d *documents) rescheduleAllDocumentTasks(owner string) error {
	documents, err := d.repository.FindAllByOwner(owner)
	if err != nil {
		return err
	}
	for _, document := range documents {
		err := d.scheduleDocumentProcessing(document)
		if err != nil {
			slog.Error("failed to schedule document for processing", "error", err.Error(), "documentID", document.ID)
		}
	}

	return nil
}

func (d *documents) emptyTrash(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-ticker.C:
			documents, err := d.repository.FindAllTrashed()
			if err != nil {
				slog.Error("failed to find all trashed documents", "error", err)
				continue
			}

			for _, document := range documents {
				if document.ShouldBeDeleted() {
					err = d.storage.Delete(document.Filepath())
					if err != nil {
						slog.Error("failed to delete document file", "error", err)
						continue
					}

					for _, path := range document.PreviewFilepaths {
						err := d.previewStorage.Delete(path)
						if err != nil {
							slog.Error("failed to delete preview file", "error", err)
							continue
						}
					}

					err := d.repository.DeleteByID(document.ID)
					if err != nil {
						slog.Error("failed to delete document", "error", err)
						continue
					}

					err = d.messages.PublishDocumentDeleted(document)
					if err != nil {
						slog.Error("failed to publish document deleted event", "error", err)
						continue
					}
				}
			}
		case <-ctx.Done():
			slog.Info("documents trash emptying stopped")
			return
		}
	}
}

func (d *documents) scheduleDocumentProcessing(document Document) error {
	payload := DocumentProcessingPayload{DocumentID: document.ID}

	err := d.taskScheduler.ScheduleTask(common.TaskTypeExtractText, payload, 3)
	if err != nil {
		return err
	}

	err = d.taskScheduler.ScheduleTask(common.TaskTypeGeneratePreviews, payload, 3)
	if err != nil {
		return err
	}

	slog.Info("scheduled async processing for document", "document_id", document.ID)
	return nil
}

type DocumentProcessingPayload struct {
	DocumentID string `json:"document_id"`
}

func newDocuments(
	repository DocumentRepository,
	storage DocumentStorage,
	previewStorage DocumentPreviewStorage,
	messages DocumentMessages,
	summarizer DocumentSummarizer,
	jobScheduler *common.JobScheduler,
	taskScheduler *common.TaskScheduler,
	shutdown *common.Shutdown) *documents {

	documents := &documents{
		repository:     repository,
		storage:        storage,
		previewStorage: previewStorage,
		messages:       messages,
		taskScheduler:  taskScheduler,
	}

	documentProcessor := newDocumentProcessor(repository, storage, previewStorage, messages, summarizer, shutdown)
	jobScheduler.Schedule(documents.emptyTrash)
	taskScheduler.Register(documentProcessor)

	// Schedule summarization after text extraction completes
	err := messages.SubscribeDocumentTextExtracted(func(document Document) error {
		payload := DocumentProcessingPayload{DocumentID: document.ID}
		return taskScheduler.ScheduleTask(common.TaskTypeSummarizeDocument, payload, 3)
	})
	if err != nil {
		panic(err)
	}

	return documents
}
