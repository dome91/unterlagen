package archive

import (
	"encoding/json"
	"errors"
	"log/slog"
	"unterlagen/features/common"
)

type AsyncDocumentProcessor struct {
	repository     DocumentRepository
	storage        DocumentStorage
	previewStorage DocumentPreviewStorage
	messages       DocumentMessages
	analyzers      map[Filetype]DocumentAnalyzer
}


func (p *AsyncDocumentProcessor) ProcessTask(task common.Task) error {
	switch task.Type {
	case common.TaskTypeExtractText:
		return p.processTextExtraction(task)
	case common.TaskTypeGeneratePreviews:
		return p.processPreviewGeneration(task)
	default:
		return errors.New("unknown task type")
	}
}

func (p *AsyncDocumentProcessor) processTextExtraction(task common.Task) error {
	var payload DocumentProcessingPayload
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		return err
	}

	document, err := p.repository.FindByID(payload.DocumentID)
	if err != nil {
		return err
	}

	analyzer, ok := p.analyzers[document.Filetype]
	if !ok {
		return ErrUnsupportedFiletype
	}

	text, err := analyzer.ExtractText(document)
	if err != nil {
		return err
	}

	document.Text = text
	if err := p.repository.Save(document); err != nil {
		return err
	}

	slog.Info("text extracted for document", "document_id", document.ID)
	return nil
}

func (p *AsyncDocumentProcessor) processPreviewGeneration(task common.Task) error {
	var payload DocumentProcessingPayload
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		return err
	}

	document, err := p.repository.FindByID(payload.DocumentID)
	if err != nil {
		return err
	}

	analyzer, ok := p.analyzers[document.Filetype]
	if !ok {
		return ErrUnsupportedFiletype
	}

	previewFilepaths, err := analyzer.GeneratePreviews(document)
	if err != nil {
		return err
	}

	document.PreviewFilepaths = previewFilepaths
	if err := p.repository.Save(document); err != nil {
		return err
	}

	slog.Info("previews generated for document", "document_id", document.ID, "count", len(previewFilepaths))
	
	return p.messages.PublishDocumentAnalyzed(document)
}


func NewAsyncDocumentProcessor(
	repository DocumentRepository,
	storage DocumentStorage,
	previewStorage DocumentPreviewStorage,
	messages DocumentMessages,
) *AsyncDocumentProcessor {
	pdfAnalyzer := NewPDFAnalyzer(storage, previewStorage)
	analyzers := make(map[Filetype]DocumentAnalyzer)
	analyzers[PDF] = pdfAnalyzer

	return &AsyncDocumentProcessor{
		repository:     repository,
		storage:        storage,
		previewStorage: previewStorage,
		messages:       messages,
		analyzers:      analyzers,
	}
}