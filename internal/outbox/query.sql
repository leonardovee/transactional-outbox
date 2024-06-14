-- name: GetOutbox :one
SELECT * FROM outbox 
WHERE id = $1 LIMIT 1;

-- name: ListOutbox :many
SELECT * FROM outbox
ORDER BY created_at ASC
FOR UPDATE SKIP LOCKED
LIMIT $1;

-- name: CreateOutbox :one
INSERT INTO outbox (
    id, aggregate_id, aggregate_type, type, payload
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: DeleteOutbox :one
DELETE FROM outbox
WHERE id = $1
RETURNING *;

