-- name: CreateAccount :one
INSERT INTO accounts (
     owner,balance,currency
) VALUES (
          $1, $2, $3
)
RETURNING *;

-- name: GetAccount :one
SELECT * from accounts
WHERE id = $1 AND currency = $2 AND owner = $3 LIMIT 1;

-- name: GetAccountByIdAndCurrency :one
SELECT * from accounts
WHERE id = $1 AND currency = $2 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * from accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * from accounts
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateAccount :one
UPDATE accounts SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE from accounts where id = $1;