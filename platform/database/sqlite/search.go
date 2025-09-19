package sqlite

import (
	"fmt"
	"log/slog"
	"strings"
	"unterlagen/features/archive"
	"unterlagen/features/search"

	"github.com/jmoiron/sqlx"
)

var _ search.SearchRepository = &SearchRepository{}

type SearchRepository struct {
	*sqlx.DB
}

// IndexDocument implements search.SearchRepository.
func (s *SearchRepository) IndexDocument(document archive.Document) error {
	// Delete existing entry if it exists
	_, err := s.Exec(`
		DELETE FROM documents_fts
		WHERE document_id = ?
	`, document.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing document from search index %s: %w", document.ID, err)
	}

	// Insert/update the document in FTS table
	_, err = s.Exec(`
		INSERT INTO documents_fts(document_id, title, filename, text, owner)
		VALUES (?, ?, ?, ?, ?)
	`, document.ID, document.Title, document.Filename, document.Text, document.Owner)
	if err != nil {
		return fmt.Errorf("failed to index document %s: %w", document.ID, err)
	}

	slog.Debug("indexed document", "id", document.ID, "title", document.Title)
	return nil
}

// SearchDocuments implements search.SearchRepository.
func (s *SearchRepository) SearchDocuments(query string, owner string, limit int) ([]search.SearchResult, error) {
	// Simple approach: use FTS for text search and regular WHERE for owner filter
	ftsQuery := s.buildFTSQuery(query)

	sqlQuery := `
		SELECT
			document_id,
			COALESCE(title, filename) as name,
			bm25(documents_fts) as rank,
			snippet(documents_fts, 3, '<mark>', '</mark>', '...', 32) as snippet
		FROM documents_fts
		WHERE documents_fts MATCH ?
		AND owner = ?
		ORDER BY bm25(documents_fts)
		LIMIT ?
	`

	var results []search.SearchResult
	err := s.Select(&results, sqlQuery, ftsQuery, owner, limit)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	return results, nil
}

// buildFTSQuery converts a user query into a proper FTS5 query
func (s *SearchRepository) buildFTSQuery(query string) string {
	// Split query into terms and escape them
	terms := strings.Fields(strings.TrimSpace(query))
	if len(terms) == 0 {
		return ""
	}

	// Use prefix matching with * to support partial matches
	var ftsTerms []string
	for _, term := range terms {
		// Escape special FTS5 characters
		escapedTerm := strings.ReplaceAll(term, `"`, `""`)
		// Add * for prefix matching to support partial matches
		ftsTerms = append(ftsTerms, escapedTerm+"*")
	}

	// Join terms with AND - all terms must match (as prefixes)
	return strings.Join(ftsTerms, " AND ")
}

// NewSearchRepository creates a new SQLite FTS search repository
func NewSearchRepository(db *sqlx.DB) *SearchRepository {
	return &SearchRepository{
		DB: db,
	}
}
