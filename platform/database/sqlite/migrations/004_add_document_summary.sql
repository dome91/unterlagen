-- +goose Up
-- Add summary column to documents table
ALTER TABLE documents ADD COLUMN summary JSON DEFAULT '{}';

-- Drop and recreate FTS table to include summary
DROP TABLE documents_fts;

CREATE VIRTUAL TABLE documents_fts USING fts5(
    document_id UNINDEXED,
    title,
    filename,
    text,
    summary,
    owner UNINDEXED
);

-- Populate FTS table with existing documents (excluding trashed ones)
INSERT INTO documents_fts(document_id, title, filename, text, summary, owner)
SELECT id, title, filename, text, COALESCE(summary, '') as summary, owner
FROM documents
WHERE trashed_at IS NULL;

-- +goose Down
-- Drop FTS table
DROP TABLE documents_fts;

-- Recreate original FTS table without summary
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

-- Remove summary column from documents table
ALTER TABLE documents DROP COLUMN summary;