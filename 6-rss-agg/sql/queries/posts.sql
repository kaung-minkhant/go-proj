-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, description, published_at, url, feed_id)
VALUES ($1, now() AT TIME ZONE 'utc', now() AT TIME ZONE 'utc', $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.* FROM posts
JOIN feed_follows on posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;

