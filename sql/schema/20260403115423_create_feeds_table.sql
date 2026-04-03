-- +goose Up
-- +goose StatementBegin
CREATE TABLE feeds(
  id uuid primary key default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null,
  name text unique,
  url text unique,
  user_id uuid references users(id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feeds;
-- +goose StatementEnd
