-- +goose Up
-- +goose StatementBegin
-- Create articles table to store scraped news articles
CREATE TABLE articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_name TEXT NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    content TEXT NOT NULL,
    image_url TEXT,
    language TEXT NOT NULL CHECK (language IN ('en', 'si', 'ta')),
    published_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd
-- +goose StatementBegin
-- Create indexes for better query performance
CREATE INDEX idx_articles_source_name ON articles(source_name);
CREATE INDEX idx_articles_language ON articles(language);
CREATE INDEX idx_articles_published_at ON articles(published_at);
CREATE INDEX idx_articles_url ON articles(url);

-- +goose StatementEnd
-- +goose StatementBegin
-- Create FTS (Full Text Search) virtual table for article search
CREATE VIRTUAL TABLE articles_fts USING fts5(
    title,
    content,
    content='articles',
    content_rowid='id'
);

-- +goose StatementEnd
-- +goose StatementBegin
-- Create triggers to maintain FTS index
CREATE TRIGGER articles_fts_insert AFTER INSERT ON articles BEGIN
    INSERT INTO articles_fts(rowid, title, content) VALUES (new.id, new.title, new.content);
END;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE TRIGGER articles_fts_delete AFTER DELETE ON articles BEGIN
    INSERT INTO articles_fts(articles_fts, rowid, title, content) VALUES('delete', old.id, old.title, old.content);
END;

-- +goose StatementEnd
-- +goose StatementBegin
CREATE TRIGGER articles_fts_update AFTER UPDATE ON articles BEGIN
    INSERT INTO articles_fts(articles_fts, rowid, title, content) VALUES('delete', old.id, old.title, old.content);
    INSERT INTO articles_fts(rowid, title, content) VALUES (new.id, new.title, new.content);
END;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Drop triggers first
DROP TRIGGER IF EXISTS articles_fts_update;
DROP TRIGGER IF EXISTS articles_fts_delete;
DROP TRIGGER IF EXISTS articles_fts_insert;

-- +goose StatementEnd
-- +goose StatementBegin
-- Drop FTS virtual table
DROP TABLE IF EXISTS articles_fts;

-- +goose StatementEnd
-- +goose StatementBegin
-- Drop indexes
DROP INDEX IF EXISTS idx_articles_url;
DROP INDEX IF EXISTS idx_articles_published_at;
DROP INDEX IF EXISTS idx_articles_language;
DROP INDEX IF EXISTS idx_articles_source_name;

-- +goose StatementEnd
-- +goose StatementBegin
-- Drop tables
DROP TABLE IF EXISTS articles;

-- +goose StatementEnd