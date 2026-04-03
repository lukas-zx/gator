-- +goose Up
-- +goose StatementBegin
CREATE TABLE feed_follows(
  id uuid primary key default gen_random_uuid(),
  created_at timestamp not null default now(),
  updated_at timestamp not null,
  user_id uuid references users(id) on delete cascade,
  feed_id uuid references feeds(id) on delete cascade,
  unique(user_id, feed_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feed_follows;
-- +goose StatementEnd
