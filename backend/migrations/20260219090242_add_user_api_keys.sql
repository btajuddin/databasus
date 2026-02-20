-- +goose Up
-- +goose StatementBegin

CREATE TABLE user_api_keys (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL,
    hashed_key   TEXT NOT NULL,
    key_prefix   TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ
);

ALTER TABLE user_api_keys
    ADD CONSTRAINT fk_user_api_keys_user_id
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE;

ALTER TABLE user_api_keys
    ADD CONSTRAINT uq_user_api_keys_user_id UNIQUE (user_id);

CREATE INDEX idx_user_api_keys_key_prefix ON user_api_keys (key_prefix);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_user_api_keys_key_prefix;
DROP TABLE IF EXISTS user_api_keys;

-- +goose StatementEnd
