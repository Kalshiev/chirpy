-- +goose up
ALTER TABLE users
ADD COLUMN chirpy_red BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose down
ALTER TABLE users
DROP COLUMN chirpy_red;