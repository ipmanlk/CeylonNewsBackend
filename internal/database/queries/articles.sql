-- name: CreateArticle :one
INSERT INTO articles (
    source_name,
    title,
    url,
    content,
    image_url,
    language,
    published_at
) VALUES (
    ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetArticleByID :one
SELECT * FROM articles WHERE id = ?;

-- name: GetArticleByURL :one
SELECT * FROM articles WHERE url = ?;

-- name: GetArticlesByLanguage :many
SELECT * FROM articles 
WHERE language = ? 
ORDER BY published_at DESC 
LIMIT ? OFFSET ?;

-- name: GetArticlesBySource :many
SELECT * FROM articles 
WHERE source_name = ? 
ORDER BY published_at DESC 
LIMIT ? OFFSET ?;

-- name: GetArticlesByLanguageAndSource :many
SELECT * FROM articles 
WHERE language = ? AND source_name = ? 
ORDER BY published_at DESC 
LIMIT ? OFFSET ?;

-- name: GetArticlesByLanguages :many
SELECT * FROM articles 
WHERE language IN (sqlc.slice('languages')) 
ORDER BY published_at DESC 
LIMIT ? OFFSET ?;

-- name: GetArticlesBySources :many
SELECT * FROM articles 
WHERE source_name IN (sqlc.slice('sources')) 
ORDER BY published_at DESC 
LIMIT ? OFFSET ?;

-- name: GetArticlesByLanguagesAndSources :many
SELECT * FROM articles 
WHERE language IN (sqlc.slice('languages')) 
AND source_name IN (sqlc.slice('sources')) 
ORDER BY published_at DESC 
LIMIT ? OFFSET ?;

-- name: GetRecentArticles :many
SELECT * FROM articles 
ORDER BY published_at DESC 
LIMIT ? OFFSET ?;

-- name: GetRecentArticlesByLanguage :many
SELECT * FROM articles 
WHERE language = ? 
ORDER BY published_at DESC 
LIMIT ? OFFSET ?;

-- name: GetRecentArticlesBySource :many
SELECT * FROM articles 
WHERE source_name = ? 
ORDER BY published_at DESC 
LIMIT ? OFFSET ?;

-- name: UpdateArticle :one
UPDATE articles 
SET 
    title = ?,
    content = ?,
    image_url = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ? 
RETURNING *;

-- name: DeleteArticle :exec
DELETE FROM articles WHERE id = ?;

-- name: GetDistinctSources :many
SELECT DISTINCT source_name FROM articles ORDER BY source_name;

-- name: GetArticleCount :one
SELECT COUNT(*) FROM articles;

-- name: GetArticleCountByLanguage :one
SELECT COUNT(*) FROM articles WHERE language = ?;

-- name: GetArticleCountBySource :one
SELECT COUNT(*) FROM articles WHERE source_name = ?;

-- name: GetArticleCountByLanguageAndSource :one
SELECT COUNT(*) FROM articles WHERE language = ? AND source_name = ?;