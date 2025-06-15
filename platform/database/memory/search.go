package memory

import (
	"strings"
	"unterlagen/features/archive"
	"unterlagen/features/search"
)

var _ search.SearchRepository = &SearchRepository{}

type IndexEntry struct {
	DocumentID string
	Name       string
	Text       string
}

type SearchRepository struct {
	index []IndexEntry
}

// IndexDocument implements search.SearchRepository.
func (s *SearchRepository) IndexDocument(document archive.Document) error {
	s.index = append(s.index, IndexEntry{
		DocumentID: document.ID,
		Name:       document.Name(),
		Text:       document.Text,
	})
	return nil
}

// SearchDocuments implements search.SearchRepository.
func (s *SearchRepository) SearchDocuments(query string, owner string, limit int) ([]search.SearchResult, error) {
	var results []search.SearchResult
	queryLower := strings.ToLower(query)

	for _, entry := range s.index {
		titleContains := strings.Contains(strings.ToLower(entry.Name), queryLower)
		textContains := strings.Contains(strings.ToLower(entry.Text), queryLower)
		if titleContains || textContains {
			rank := 0.0
			if titleContains {
				rank += 1.0
			}
			if textContains {
				rank += 0.5
			}

			// Generate snippet from text content
			snippet := generateSnippet(entry.Text, query, 150)

			results = append(results, search.SearchResult{
				DocumentID: entry.DocumentID,
				Name:       entry.Name,
				Snippet:    snippet,
				Rank:       rank,
			})
		}
	}
	return results, nil
}

// generateSnippet creates a text snippet around the search query
func generateSnippet(text, query string, maxLength int) string {
	if text == "" {
		return ""
	}

	textLower := strings.ToLower(text)
	queryLower := strings.ToLower(query)

	// Find the first occurrence of the query
	index := strings.Index(textLower, queryLower)
	if index == -1 {
		// If query not found, return beginning of text
		if len(text) <= maxLength {
			return text
		}
		return text[:maxLength] + "..."
	}

	// Calculate snippet bounds
	start := index - 50
	if start < 0 {
		start = 0
	}

	end := start + maxLength
	if end > len(text) {
		end = len(text)
	}

	snippet := text[start:end]

	// Add ellipsis if needed
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(text) {
		snippet = snippet + "..."
	}

	return snippet
}

func NewSearchRepository() *SearchRepository {
	return &SearchRepository{}
}
