-- name: CreateCustomer :one
INSERT INTO customer (
        email,
        hashed_password,
        username,
        firstname,
        lastname,
        gender,
        state_of_origin
) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetCustomerByID :one
SELECT * FROM customer WHERE id= $1;

-- name: GetCustomerByEmail :one
SELECT * FROM customer WHERE email= $1;

-- name: ListCustomer :many
SELECT * FROM customer ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateCustomerPassword :one
UPDATE customer SET hashed_password = $1, updated_at = $2
WHERE id = $3 RETURNING *;

-- name: DeleteCustomer :exec
DELETE FROM customer WHERE id = $1;

-- name: DeleteAllCustomer :exec
DELETE FROM customer;