-- +goose Up
CREATE TABLE users(
  id uuid primary key default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null,
  name text unique
);

-- +goose Down
DROP TABLE users;
