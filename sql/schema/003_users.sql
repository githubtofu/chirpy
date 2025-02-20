-- +goose Up
alter table users
add hashed_password TEXT NOT NULL DEFAULT 'unset';

-- +goose Down
alter table users
drop hashed_password;

