package bleve

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"unterlagen/features/archive"
	"unterlagen/features/search"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
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
func createIndexMapping() mapping.IndexMapping {
	// Create a mapping for the document
	documentMapping := bleve.NewDocumentMapping()

	// ID field - stored but not analyzed
	idFieldMapping := bleve.NewTextFieldMapping()
	idFieldMapping.Store = true
	idFieldMapping.Index = false
	documentMapping.AddFieldMappingsAt("id", idFieldMapping)

	// Name field - analyzed and stored for searching and highlighting
	nameFieldMapping := bleve.NewTextFieldMapping()
	nameFieldMapping.Store = true
	nameFieldMapping.Index = true
	nameFieldMapping.Analyzer = "en"
	documentMapping.AddFieldMappingsAt("name", nameFieldMapping)

	// Text field - analyzed for full-text search
	textFieldMapping := bleve.NewTextFieldMapping()
	textFieldMapping.Store = true
	textFieldMapping.Index = true
	textFieldMapping.Analyzer = "en"
	textFieldMapping.IncludeTermVectors = true
	documentMapping.AddFieldMappingsAt("text", textFieldMapping)

	// Owner field - keyword field for filtering
	ownerFieldMapping := bleve.NewKeywordFieldMapping()
	ownerFieldMapping.Store = true
	ownerFieldMapping.Index = true
	documentMapping.AddFieldMappingsAt("owner", ownerFieldMapping)

	// FolderID field - keyword field for filtering
	folderFieldMapping := bleve.NewKeywordFieldMapping()
	folderFieldMapping.Store = true
	folderFieldMapping.Index = true
	documentMapping.AddFieldMappingsAt("folder_id", folderFieldMapping)

	// Filename field - stored for display
	filenameFieldMapping := bleve.NewTextFieldMapping()
	filenameFieldMapping.Store = true
	filenameFieldMapping.Index = false
	documentMapping.AddFieldMappingsAt("filename", filenameFieldMapping)

	// IsTrashed field - boolean field for filtering
	trashedFieldMapping := bleve.NewBooleanFieldMapping()
	trashedFieldMapping.Store = true
	trashedFieldMapping.Index = true
	documentMapping.AddFieldMappingsAt("is_trashed", trashedFieldMapping)

	// Date fields
	dateFieldMapping := bleve.NewDateTimeFieldMapping()
	dateFieldMapping.Store = true
	dateFieldMapping.Index = true
	documentMapping.AddFieldMappingsAt("created_at", dateFieldMapping)
	documentMapping.AddFieldMappingsAt("updated_at", dateFieldMapping)

	// Create the index mapping
	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("document", documentMapping)
	indexMapping.DefaultMapping = documentMapping

	return indexMapping
}

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

		// Create snippet from highlights or fallback to text
		result.Snippet = s.createSnippet(hit)

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
func NewSearchRepository(indexPath string) (*SearchRepository, error) {
	// Check if index already exists
	var index bleve.Index

	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		// Create new index
		indexMapping := createIndexMapping()
		index, err = bleve.New(indexPath, indexMapping)
		if err != nil {
			return nil, fmt.Errorf("failed to create new search index: %w", err)
		}
		slog.Info("created new search index", "path", indexPath)
	} else {
		// Open existing index
		index, err = bleve.Open(indexPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open existing search index: %w", err)
		}
		slog.Info("opened existing search index", "path", indexPath)
	}

	return &SearchRepository{
		index: index,
	}, nil
}
