package archive

import "log/slog"

var _ DocumentAnalyzer = &PDFAnalyzer{}

type PDFAnalyzer struct {
	documentStorage DocumentStorage
	previewStorage  DocumentPreviewStorage
}

// ExtractText implements DocumentAnalyzer.
func (p *PDFAnalyzer) ExtractText(document Document) (string, error) {
	slog.Info("PDFAnalyzer ExtractText not implemented yet.")
	return "", nil
}

// GeneratePreviews implements DocumentAnalyzer.
func (p *PDFAnalyzer) GeneratePreviews(document Document) ([]string, error) {
	slog.Info("PDFAnalyzer GeneratePreviews not implemented yet.")
	return []string{}, nil
}

func NewPDFAnalyzer(documentStorage DocumentStorage, previewStorage DocumentPreviewStorage) *PDFAnalyzer {
	return &PDFAnalyzer{
		documentStorage: documentStorage,
		previewStorage:  previewStorage,
	}
}
