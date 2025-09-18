-- +goose Up
ALTER TABLE documents ADD COLUMN title TEXT;

-- Initialize title with filename without extension for existing documents
UPDATE documents SET title = REPLACE(filename, SUBSTR(filename, INSTR(filename, '.')), '') WHERE title IS NULL;

-- Make title required for new documents
-- Note: SQLite doesn't support adding NOT NULL constraints with ALTER TABLE,
-- but we'll enforce this in the application code

-- +goose Down
ALTER TABLE documents DROP COLUMN title;