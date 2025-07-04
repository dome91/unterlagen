package bleve

import (
	"fmt"
	"log/slog"
	"strings"
	"unterlagen/features/archive"
	"unterlagen/features/search"

	"github.com/blevesearch/bleve/v2"
	bleveSearch "github.com/blevesearch/bleve/v2/search"
)

var _ search.SearchRepository = &SearchRepository{}

// DocumentIndex represents a document in the search index
type DocumentIndex struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Text  string `json:"text"`
	Owner string `json:"owner"`
}

type SearchRepository struct {
	index bleve.Index
}

// createIndexMapping creates the mapping for document indexing
// IndexDocument implements search.SearchRepository.
func (s *SearchRepository) IndexDocument(document archive.Document) error {
	// Convert archive.Document to DocumentIndex
	docIndex := DocumentIndex{
		ID:    document.ID,
		Name:  document.Name(),
		Text:  document.Text,
		Owner: document.Owner,
	}

	// Index the document
	err := s.index.Index(document.ID, docIndex)
	if err != nil {
		return fmt.Errorf("failed to index document %s: %w", document.ID, err)
	}

	slog.Debug("indexed document", "id", document.ID, "name", document.Name())
	return nil
}

// SearchDocuments implements search.SearchRepository.
func (s *SearchRepository) SearchDocuments(query string, owner string, limit int) ([]search.SearchResult, error) {
	// Create the search query
	searchQuery := bleve.NewConjunctionQuery()

	// Add the text search query
	if query != "" {
		// Create a multi-field query searching both name and text
		nameQuery := bleve.NewMatchQuery(query)
		nameQuery.SetField("name")
		nameQuery.SetBoost(2.0) // Boost name matches

		textQuery := bleve.NewMatchQuery(query)
		textQuery.SetField("text")

		// Combine name and text queries with OR
		contentQuery := bleve.NewDisjunctionQuery(nameQuery, textQuery)
		searchQuery.AddQuery(contentQuery)
	}

	// Filter by owner
	ownerQuery := bleve.NewTermQuery(owner)
	ownerQuery.SetField("owner")
	searchQuery.AddQuery(ownerQuery)

	// Create search request
	searchRequest := bleve.NewSearchRequest(searchQuery)
	searchRequest.Size = limit
	searchRequest.Highlight = bleve.NewHighlight()
	searchRequest.Highlight.AddField("name")
	searchRequest.Highlight.AddField("text")
	searchRequest.Fields = []string{"id", "name", "filename", "text"}

	// Execute search
	searchResult, err := s.index.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Convert results
	results := make([]search.SearchResult, 0, len(searchResult.Hits))
	for _, hit := range searchResult.Hits {
		result := search.SearchResult{
			DocumentID: hit.ID,
			Rank:       hit.Score,
		}

		// Extract name from fields
		if nameField, exists := hit.Fields["name"]; exists {
			if name, ok := nameField.(string); ok {
				result.Name = name
			}
		}

		results = append(results, result)
	}

	return results, nil
}

// createSnippet creates a snippet from search highlights or document text
func (s *SearchRepository) createSnippet(hit *bleveSearch.DocumentMatch) string {
	const maxSnippetLength = 200

	// Try to use highlights first
	if hit.Fragments != nil {
		// Prefer name highlights
		if nameFragments, exists := hit.Fragments["name"]; exists && len(nameFragments) > 0 {
			return strings.Join(nameFragments, " ... ")
		}

		// Use text highlights
		if textFragments, exists := hit.Fragments["text"]; exists && len(textFragments) > 0 {
			return strings.Join(textFragments, " ... ")
		}
	}

	// Fallback to truncated text content
	if textField, exists := hit.Fields["text"]; exists {
		if text, ok := textField.(string); ok && text != "" {
			if len(text) > maxSnippetLength {
				// Find a good break point near the limit
				breakPoint := maxSnippetLength
				for i := maxSnippetLength; i > maxSnippetLength-50 && i > 0; i-- {
					if text[i] == ' ' || text[i] == '.' || text[i] == '\n' {
						breakPoint = i
						break
					}
				}
				return text[:breakPoint] + "..."
			}
			return text
		}
	}

	// Final fallback to filename
	if nameField, exists := hit.Fields["name"]; exists {
		if name, ok := nameField.(string); ok {
			return name
		}
	}

	return "No preview available"
}

// NewSearchRepository creates a new Bleve search repository
func NewSearchRepository(index bleve.Index) *SearchRepository {
	return &SearchRepository{
		index: index,
	}
}
