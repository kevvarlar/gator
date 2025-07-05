-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*
FROM posts
INNER JOIN feed_follows f ON f.feed_id = posts.feed_id
INNER JOIN users u ON u.id = f.user_id
WHERE u.name = $1
ORDER BY published_at DESC
LIMIT $2;