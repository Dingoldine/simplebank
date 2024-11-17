-- name: CreateEntry :one
INSERT INTO Entries (
    account_id,
    amount
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM Entries
WHERE id = $1;

-- name: GetEntriesByAccountId :many
SELECT * FROM Entries
WHERE account_id = $1;

-- name: ListEntrys :many
SELECT * FROM Entries
ORDER BY account_id
LIMIT $1
OFFSET $2;

-- name: UpdateEntry :exec
UPDATE Entries
SET amount = $2
WHERE id = $1;

-- name: DeleteEntry :exec
DELETE FROM Entries WHERE id = $1;