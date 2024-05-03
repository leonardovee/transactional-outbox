-- name: GetOrder :one
SELECT * FROM orders 
WHERE id = $1 LIMIT 1;

-- name: GetAggregate :one
SELECT * FROM orders 
WHERE aggregate_id = $1 LIMIT 1;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY created_at;

-- name: CreateOrder :one
INSERT INTO orders (
    id, aggregate_id, status, total
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;
