package llm

import (
	"slices"
	"strings"
	"unicode"
	"unterlagen/features/assistant"
	"unterlagen/platform/configuration"
)

var _ assistant.Chunker = &FixedSizeChunker{}
var _ assistant.Chunker = &RecursiveChunker{}

type FixedSizeChunker struct {
	maxChunkSize int
}

func (f *FixedSizeChunker) Chunk(text string) ([]string, error) {
	runes := []rune(text)
	chunks := slices.Chunk(runes, f.maxChunkSize)
	var result []string
	for chunk := range chunks {
		result = append(result, string(chunk))
	}

	return result, nil
}

func NewFixedSizeChunker(config configuration.Configuration) *FixedSizeChunker {
	return &FixedSizeChunker{
		maxChunkSize: config.Assistant.Chunker.MaxChunkSize,
	}
}

type RecursiveChunker struct {
	maxChunkSize int
	chunkOverlap int
}

// Chunk implements assistant.Chunker.
func (r *RecursiveChunker) Chunk(text string) ([]string, error) {
	// Clean and normalize the input text
	cleanedText := strings.TrimSpace(string(text))
	if cleanedText == "" {
		return []string{}, nil
	}

	// If text is smaller than max chunk size, return it as a single chunk
	if len(cleanedText) <= r.maxChunkSize {
		return []string{cleanedText}, nil
	}

	return r.chunk(cleanedText), nil
}

func (r *RecursiveChunker) chunk(text string) []string {
	var chunks []string
	if len(text) <= r.maxChunkSize {
		return []string{text}
	}

	splitIndex := r.findSplitPoint(text)

	firstHalf := strings.TrimSpace(text[:splitIndex])
	secondHalf := strings.TrimSpace(text[splitIndex:])

	if len(firstHalf) > 0 {
		chunks = append(chunks, r.chunk(firstHalf)...)
	}
	if len(secondHalf) > 0 {
		chunks = append(chunks, r.chunk(secondHalf)...)
	}

	// Add overlap if needed
	if r.chunkOverlap > 0 {
	}

	return chunks
}

func (r *RecursiveChunker) findSplitPoint(text string) int {
	targetIndex := r.maxChunkSize
	for i := targetIndex; i >= 0 && i < len(text); i-- {
		if r.isSentenceBoundary(text, i) {
			return i + 1
		}
	}

	for i := targetIndex; i >= 0 && i < len(text); i-- {
		if unicode.IsSpace(rune(text[i])) {
			return i
		}
	}

	return targetIndex
}

func (r *RecursiveChunker) isSentenceBoundary(text string, index int) bool {
	if index >= len(text) || index < 0 {
		return false
	}

	// Check for common sentence endings
	sentenceEnders := []string{". ", "? ", "! ", ".\n", "?\n", "!\n"}
	for _, ender := range sentenceEnders {
		if index+len(ender) <= len(text) {
			if text[index:index+len(ender)] == ender {
				return true
			}
		}
	}

	return false
}

func NewRecursiveChunker(configuration configuration.Configuration) *RecursiveChunker {
	return &RecursiveChunker{
		maxChunkSize: configuration.Assistant.Chunker.MaxChunkSize,
		chunkOverlap: configuration.Assistant.Chunker.ChunkOverlap,
	}
}
