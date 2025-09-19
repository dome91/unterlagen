-- +goose Up
-- Create FTS5 virtual table for document search
CREATE VIRTUAL TABLE documents_fts USING fts5(
    document_id UNINDEXED,
    title,
    filename,
    text,
    owner UNINDEXED
);

-- Populate FTS table with existing documents (excluding trashed ones)
INSERT INTO documents_fts(document_id, title, filename, text, owner)
SELECT id, title, filename, text, owner
FROM documents
WHERE trashed_at IS NULL;

-- +goose Down
DROP TABLE documents_fts;