-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS lists (
    id text primary key,
    title text not null,
    user_id text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    state int not null
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS lists;
-- +goose StatementEnd
