-- name: CreateBook :exec
INSERT INTO books (title, author, genre, price)
VALUES ($1, $2, $3, $4);

-- name: GetBookByID :one
SELECT id, title, author, genre, price, created_at, updated_at
FROM books
WHERE id = $1;

-- name: ListBooks :many
SELECT id, title, author, genre, price, created_at, updated_at
FROM books
ORDER BY created_at DESC;

-- name: UpdateBook :exec
UPDATE books
SET title = $2, author = $3, genre = $4, price = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteBook :exec
DELETE FROM books
WHERE id = $1;
