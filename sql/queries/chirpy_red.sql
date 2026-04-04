-- name: UpgradeToChirpyRed :exec
UPDATE users
SET chirpy_red = TRUE
WHERE id = $1;

-- name: DowngradeFromChirpyRed :exec
UPDATE users
SET chirpy_red = FALSE
WHERE id = $1;