-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "user" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tg_id BIGINT UNIQUE NOT NULL,
    tg_username VARCHAR(255),
    chat_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE(chat_id)
);

CREATE TABLE IF NOT EXISTS "event" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id VARCHAR(255) NOT NULL,
    owner_id VARCHAR(255) NOT NULL,
    text TEXT NOT NULL,
    notify_at VARCHAR(255),
    schedule VARCHAR(255),
    delta VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "event";
DROP TABLE IF EXISTS "user";
-- +goose StatementEnd
