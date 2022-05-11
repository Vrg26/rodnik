-- +goose Up
CREATE TABLE friendships(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_on timestamp DEFAULT now(),
    friend_from uuid NOT NULL REFERENCES users ("id"),
    friend_to uuid NOT NULL REFERENCES users ("id")
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE freindships;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
