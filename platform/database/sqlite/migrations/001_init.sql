-- +goose Up
CREATE TABLE users (
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL,
    password_change_necessary INTEGER NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    PRIMARY KEY (username)
);

CREATE TABLE folders (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    parent_id TEXT,
    owner TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (parent_id) REFERENCES folders (id) ON DELETE CASCADE,
    FOREIGN KEY (owner) REFERENCES users (username) ON DELETE CASCADE
);

CREATE TABLE documents (
    id TEXT NOT NULL,
    filename TEXT NOT NULL,
    filetype TEXT NOT NULL,
    filesize INTEGER NOT NULL,
    text TEXT,
    trashed_at DATETIME,
    folder_id TEXT NOT NULL,
    owner TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (folder_id) REFERENCES folders (id) ON DELETE CASCADE,
    FOREIGN KEY (owner) REFERENCES users (username) ON DELETE CASCADE
);

CREATE TABLE documents_previews (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    document_id TEXT NOT NULL,
    filepath TEXT NOT NULL,
    page_number INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (document_id) REFERENCES documents (id) ON DELETE CASCADE,
    UNIQUE(document_id, page_number)
);

CREATE INDEX idx_documents_previews_document_id ON documents_previews(document_id);

CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    status TEXT NOT NULL,
    payload TEXT NOT NULL,
    error TEXT,
    attempts INTEGER NOT NULL DEFAULT 0,
    max_attempts INTEGER NOT NULL DEFAULT 3,
    next_run_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);

CREATE INDEX idx_tasks_status_next_run ON tasks(status, next_run_at);

-- +goose Down
DROP INDEX idx_tasks_status_next_run;

DROP TABLE tasks;

DROP INDEX idx_documents_previews_document_id;

DROP TABLE documents_previews;

DROP TABLE documents;

DROP TABLE folders;

DROP TABLE users;
