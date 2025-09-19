package llm

import (
	"time"
	"unterlagen/features/archive"
	"unterlagen/features/assistant"
)

var (
	_ assistant.Embedder        = &DumbAI{}
	_ assistant.Answerer        = &DumbAI{}
	_ archive.DocumentSummarizer = &DumbAI{}
)

// DumbAI does nothing and is used when no AI is provided via configuration
type DumbAI struct {
}

// Answer implements assistant.LLM.
func (ai *DumbAI) Answer(question string, nodes []assistant.Node) (string, error) {
	time.Sleep(5 * time.Second)
	return "I have no idea...", nil
}

// Generate implements assistant.Embedder.
func (ai *DumbAI) Generate(text string) (assistant.Embeddings, error) {
	return assistant.Embeddings{0.0, 0.0, 0.0, 0.0}, nil
}

// SummarizeText implements archive.DocumentSummarizer.
func (ai *DumbAI) SummarizeText(text string) (archive.DocumentSummary, error) {
	return archive.DocumentSummary{}, nil // Return empty summary when no LLM is available
}

func NewDumbAI() *DumbAI {
	return &DumbAI{}
}
