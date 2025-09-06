-- name: SearchArticles :many
SELECT a.* FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ?
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: SearchArticlesByLanguage :many
SELECT a.* FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.language = ?
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: SearchArticlesBySource :many
SELECT a.* FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.source_name = ?
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: SearchArticlesByLanguageAndSource :many
SELECT a.* FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.language = ? AND a.source_name = ?
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: SearchArticlesByLanguages :many
SELECT a.* FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.language IN (sqlc.slice('languages'))
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: SearchArticlesBySources :many
SELECT a.* FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.source_name IN (sqlc.slice('sources'))
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: SearchArticlesByLanguagesAndSources :many
SELECT a.* FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.language IN (sqlc.slice('languages')) AND a.source_name IN (sqlc.slice('sources'))
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: SearchArticlesHighlight :many
SELECT 
    a.*,
    snippet(articles_fts, 0, '<mark>', '</mark>', '...', 32) as title_snippet,
    snippet(articles_fts, 1, '<mark>', '</mark>', '...', 64) as content_snippet
FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ?
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: SearchArticlesHighlightByLanguage :many
SELECT 
    a.*,
    snippet(articles_fts, 0, '<mark>', '</mark>', '...', 32) as title_snippet,
    snippet(articles_fts, 1, '<mark>', '</mark>', '...', 64) as content_snippet
FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.language = ?
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: SearchArticlesHighlightBySource :many
SELECT 
    a.*,
    snippet(articles_fts, 0, '<mark>', '</mark>', '...', 32) as title_snippet,
    snippet(articles_fts, 1, '<mark>', '</mark>', '...', 64) as content_snippet
FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.source_name = ?
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: SearchArticlesHighlightByLanguageAndSource :many
SELECT 
    a.*,
    snippet(articles_fts, 0, '<mark>', '</mark>', '...', 32) as title_snippet,
    snippet(articles_fts, 1, '<mark>', '</mark>', '...', 64) as content_snippet
FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.language = ? AND a.source_name = ?
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: SearchArticlesHighlightByLanguages :many
SELECT 
    a.*,
    snippet(articles_fts, 0, '<mark>', '</mark>', '...', 32) as title_snippet,
    snippet(articles_fts, 1, '<mark>', '</mark>', '...', 64) as content_snippet
FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.language IN (sqlc.slice('languages'))
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: SearchArticlesHighlightBySources :many
SELECT 
    a.*,
    snippet(articles_fts, 0, '<mark>', '</mark>', '...', 32) as title_snippet,
    snippet(articles_fts, 1, '<mark>', '</mark>', '...', 64) as content_snippet
FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.source_name IN (sqlc.slice('sources'))
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: SearchArticlesHighlightByLanguagesAndSources :many
SELECT 
    a.*,
    snippet(articles_fts, 0, '<mark>', '</mark>', '...', 32) as title_snippet,
    snippet(articles_fts, 1, '<mark>', '</mark>', '...', 64) as content_snippet
FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.language IN (sqlc.slice('languages')) AND a.source_name IN (sqlc.slice('sources'))
ORDER BY published_at DESC
LIMIT ? OFFSET ?;

-- name: GetSearchResultCount :one
SELECT COUNT(*) FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ?;

-- name: GetSearchResultCountByLanguage :one
SELECT COUNT(*) FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.language = ?;

-- name: GetSearchResultCountBySource :one
SELECT COUNT(*) FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.source_name = ?;

-- name: GetSearchResultCountByLanguageAndSource :one
SELECT COUNT(*) FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.language = ? AND a.source_name = ?;

-- name: GetSearchResultCountByLanguages :one
SELECT COUNT(*) FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.language IN (sqlc.slice('languages'));

-- name: GetSearchResultCountBySources :one
SELECT COUNT(*) FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.source_name IN (sqlc.slice('sources'));

-- name: GetSearchResultCountByLanguagesAndSources :one
SELECT COUNT(*) FROM articles a
JOIN articles_fts fts ON a.id = fts.rowid
WHERE articles_fts MATCH ? AND a.language IN (sqlc.slice('languages')) AND a.source_name IN (sqlc.slice('sources'));
