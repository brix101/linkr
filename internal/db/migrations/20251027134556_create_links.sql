-- +goose Up
-- +goose StatementBegin
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    code VARCHAR(16) UNIQUE NOT NULL,
    url TEXT NOT NULL,
    user_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ DEFAULT NULL
);

-- Create index for URL lookups
CREATE INDEX idx_links_url ON links (url);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_links_url;
DROP TABLE IF EXISTS links;
-- +goose StatementEnd
