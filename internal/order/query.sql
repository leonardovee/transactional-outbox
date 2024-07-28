-- name: GetOrder :one
SELECT * FROM orders 
WHERE id = $1 LIMIT 1;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY created_at;

-- name: CreateOrder :one
INSERT INTO orders (
    id, status, total
) VALUES (
    $1, $2, $3
)
RETURNING *;
