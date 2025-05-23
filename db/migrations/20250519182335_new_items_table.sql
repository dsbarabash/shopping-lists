-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items (
    id text primary key,
    title text not null,
    comment text,
    is_done bool,
    user_id text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    shopping_list_id text not null
    --FOREIGN KEY (shopping_list_id) REFERENCES lists(id)
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
-- +goose StatementEnd
