-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  username varchar(150) UNIQUE,
  email varchar(100) UNIQUE,
  passwordHash varchar(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users
-- +goose StatementEnd
