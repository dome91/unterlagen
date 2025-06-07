package archive

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"log/slog"
	"path"
	"strings"
	"time"

	"github.com/klippa-app/go-pdfium"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/responses"
	"github.com/klippa-app/go-pdfium/webassembly"
)

var _ DocumentAnalyzer = &PDFAnalyzer{}

type PDFAnalyzer struct {
	documentStorage DocumentStorage
	previewStorage  DocumentPreviewStorage
}

func getPool() (pdfium.Pool, error) {
	return webassembly.Init(webassembly.Config{
		MinIdle:  1, // Makes sure that at least x workers are always available
		MaxIdle:  1, // Makes sure that at most x workers are ever available
		MaxTotal: 1, // Maxium amount of workers in total, allows the amount of workers to grow when needed, items between total max and idle max are automatically cleaned up, while idle workers are kept alive so they can be used directly.
	})
}

func (p *PDFAnalyzer) withInstance(document Document, block func(instance pdfium.Pdfium, pdfDocument *responses.OpenDocument) error) error {
	pool, err := getPool()
	if err != nil {
		return err
	}
	defer pool.Close()
	instance, err := pool.GetInstance(time.Second * 30)
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
			pageRender, err := instance.RenderPageInDPI(&requests.RenderPageInDPI{
				DPI: 200, // The DPI to render the page in.
				Page: requests.Page{
					ByIndex: &requests.PageByIndex{
						Document: pdfDocument.Document,
						Index:    page,
					},
				},
			})
			if err != nil {
				slog.Warn("failed to render page", "err", err.Error())
				continue
			}
			defer pageRender.Cleanup()

			filepath := path.Join(document.PreviewPrefix(), fmt.Sprintf("page%d.jpeg", page))
			var buf bytes.Buffer
			err = jpeg.Encode(&buf, pageRender.Result.Image, &jpeg.Options{
				Quality: 90,
			})
			if err != nil {
				slog.Warn("failed to encode render to jpeg", "err", err.Error())
				continue
			}

			err = p.previewStorage.Store(filepath, &buf)
			if err != nil {
				slog.Warn("failed to store preview", "err", err.Error())
			}

			filepaths = append(filepaths, filepath)
		}

		return nil
	})

	return filepaths, err
}

func NewPDFAnalyzer(documentStorage DocumentStorage, previewStorage DocumentPreviewStorage) *PDFAnalyzer {
	return &PDFAnalyzer{
		documentStorage: documentStorage,
		previewStorage:  previewStorage,
	}
}
