-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    chat_id INT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    text VARCHAR(5000) NOT NULL 
        CHECK (LENGTH(text) BETWEEN 1 AND 5000),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages CASCADE;
-- +goose StatementEnd
