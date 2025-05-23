-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id text primary key,
    name text not null,
    UNIQUE (name),
    password text not null,
    state int not null
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
