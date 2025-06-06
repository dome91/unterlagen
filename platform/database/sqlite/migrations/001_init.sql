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

-- +goose Down
DROP TABLE users;

DROP TABLE folders;

DROP TABLE documents;
