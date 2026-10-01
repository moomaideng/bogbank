-- +goose Up
CREATE TABLE categories (
    id         uuid PRIMARY KEY,
    user_id    uuid NOT NULL,
    name       text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (user_id, name)
);

-- +goose Down
DROP TABLE IF EXISTS categories;
