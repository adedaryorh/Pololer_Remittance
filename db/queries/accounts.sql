-- name: CreateAccount :one
INSERT INTO accounts (
        customer_id,
        account_type,
        account_status,
        balance,
        currency
) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetAccountByID :one
SELECT * FROM accounts WHERE id= $1;

-- name: GetAccountByCustomerId :many
SELECT * FROM accounts WHERE customer_id= $1;

-- name: ListAccounts :many
SELECT * FROM accounts ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateAccountBalance :one
UPDATE accounts SET balance = $1 WHERE id = $2 RETURNING *;

-- name: UpdateAccountNumber :one
UPDATE accounts SET account_number = $1 WHERE id = $2 RETURNING *;

-- name: UpdateAccountBalanceManual :one
UPDATE accounts SET balance = balance + sqlc.arg(amount) WHERE id = sqlc.arg(id) RETURNING *;

-- name: UpdateAccountType :one
UPDATE accounts SET account_type = $1 WHERE id = $2 RETURNING *;

-- name: UpdateAccountStatus :one
UPDATE accounts SET account_status = $1 WHERE id = $2 RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;

-- name: DeleteAllAccount :exec
DELETE FROM accounts;