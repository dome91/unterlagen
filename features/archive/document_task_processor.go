package archive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io"
	"log/slog"
	"path"
	"strings"
	"time"
	"unterlagen/features/common"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/webassembly"
)

type DocumentTaskProcessor struct {
	repository     DocumentRepository
	storage        DocumentStorage
	previewStorage DocumentPreviewStorage
	messages       DocumentMessages
	analyzers      map[Filetype]DocumentAnalyzer
}

func (p *DocumentTaskProcessor) Name() string {
	return "DocumentTaskProcessor"
}

func (p *DocumentTaskProcessor) ProcessTask(task common.Task) error {
	switch task.Type {
	case common.TaskTypeExtractText:
		return p.processTextExtraction(task)
	case common.TaskTypeGeneratePreviews:
		return p.processPreviewGeneration(task)
	default:
		return nil
	}
}

func (p *DocumentTaskProcessor) ResponsibleFor() []common.TaskType {
	return []common.TaskType{
		common.TaskTypeExtractText,
		common.TaskTypeGeneratePreviews,
	}
}

func (p *DocumentTaskProcessor) processTextExtraction(task common.Task) error {
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

func (p *DocumentTaskProcessor) processPreviewGeneration(task common.Task) error {
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

	return p.messages.PublishDocumentTextExtracted(document)
}

func newDocumentProcessor(
	repository DocumentRepository,
	storage DocumentStorage,
	previewStorage DocumentPreviewStorage,
	messages DocumentMessages,
	shutdown *common.Shutdown,
) *DocumentTaskProcessor {
	pdfAnalyzer := NewPDFAnalyzer(storage, previewStorage, shutdown)
	analyzers := make(map[Filetype]DocumentAnalyzer)
	analyzers[PDF] = pdfAnalyzer

	return &DocumentTaskProcessor{
		repository:     repository,
		storage:        storage,
		previewStorage: previewStorage,
		messages:       messages,
		analyzers:      analyzers,
	}
}

var _ DocumentAnalyzer = &PDFAnalyzer{}

type PDFAnalyzer struct {
	documentStorage DocumentStorage
	previewStorage  DocumentPreviewStorage
	pool            pdfium.Pool
}

func (p *PDFAnalyzer) withInstance(document Document, block func(instance pdfium.Pdfium, pdfDocument *responses.OpenDocument) error) error {
	instance, err := p.pool.GetInstance(time.Second * 30)
	if err != nil {
		return err
	}
	defer instance.Close()

	var pdfDocument *responses.OpenDocument
	err = p.documentStorage.Retrieve(document.Filepath(), func(r io.Reader) error {
		data, err := io.ReadAll(r)
		if err != nil {
			return err
		}

		pdfDocument, err = instance.OpenDocument(&requests.OpenDocument{
			File: &data,
		})
		return err
	})
	if err != nil {
		return err
	}
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{
		Document: pdfDocument.Document,
	})

	return block(instance, pdfDocument)
}

// ExtractText implements DocumentAnalyzer.
func (p *PDFAnalyzer) ExtractText(document Document) (string, error) {
	var text string
	p.withInstance(document, func(instance pdfium.Pdfium, pdfDocument *responses.OpenDocument) error {
		pageCount, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
			Document: pdfDocument.Document,
		})
		if err != nil {
			return err
		}

		var textBuilder strings.Builder
		for page := range pageCount.PageCount {
			pageText, err := instance.GetPageText(&requests.GetPageText{
				Page: requests.Page{ByIndex: &requests.PageByIndex{
					Document: pdfDocument.Document,
					Index:    page,
				}},
			})
			if err != nil {
				slog.Warn("Failed to extract text from page", "error", err, "document", document.ID, "page", page)
				continue
			}

			if pageText.Text != "" {
				textBuilder.WriteString(pageText.Text)
				textBuilder.WriteString("\n\n")
			}
		}

		text = strings.TrimSpace(textBuilder.String())
		return nil
	})

	return text, nil
}

// GeneratePreviews implements DocumentAnalyzer.
func (p *PDFAnalyzer) GeneratePreviews(document Document) ([]string, error) {
	var filepaths []string
	err := p.withInstance(document, func(instance pdfium.Pdfium, pdfDocument *responses.OpenDocument) error {
		pageCount, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{
			Document: pdfDocument.Document,
		})
		if err != nil {
			return err
		}

		for page := range pageCount.PageCount {
			filepath := path.Join(document.PreviewPrefix(), fmt.Sprintf("page%d.jpeg", page))
			err := p.generatePreview(instance, pdfDocument, filepath, page)
			if err != nil {
				continue
			}
			filepaths = append(filepaths, filepath)
		}

		return nil
	})

	return filepaths, err
}

func (p *PDFAnalyzer) generatePreview(instance pdfium.Pdfium, pdfDocument *responses.OpenDocument, filepath string, page int) error {
	pageRender, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
		DPI: 80, // The DPI to render the page in.
		Page: requests.Page{
			ByIndex: &requests.PageByIndex{
				Document: pdfDocument.Document,
				Index:    page,
			},
		},
	})
	if err != nil {
		slog.Warn("failed to render page", "err", err.Error())
		return err
	}
	defer pageRender.Cleanup()

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, pageRender.Result.Image, &jpeg.Options{
		Quality: 90,
	})
	if err != nil {
		slog.Warn("failed to encode render to jpeg", "err", err.Error())
		return err
	}

	err = p.previewStorage.Store(filepath, &buf)
	if err != nil {
		slog.Warn("failed to store preview", "err", err.Error())
	}
	return err
}

func NewPDFAnalyzer(documentStorage DocumentStorage, previewStorage DocumentPreviewStorage, shutdown *common.Shutdown) *PDFAnalyzer {
	pool, err := webassembly.Init(webassembly.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 4, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
	})
	if err != nil {
		panic(err)
	}

	shutdown.AddCallback(func() {
		err := pool.Close()
		if err != nil {
			slog.Error("failed to close wasm pool", "error", err.Error())
		}
	})

	return &PDFAnalyzer{
		documentStorage: documentStorage,
		previewStorage:  previewStorage,
		pool:            pool,
	}
}
