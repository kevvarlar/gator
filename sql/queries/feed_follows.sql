-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT
    i.*,
    f.name AS feed_name,
    u.name AS user_name
FROM inserted_feed_follow i
INNER JOIN feeds f ON f.id = i.feed_id
INNER JOIN users u ON u.id = i.user_id;

-- name: GetFeedFollowsForUser :many
SELECT
    ff.*,
    f.name AS feed_name,
    u.name AS user_name
FROM feed_follows ff
INNER JOIN feeds f ON f.id = ff.feed_id
INNER JOIN users u ON u.id = ff.user_id
WHERE u.name = $1;