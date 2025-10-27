-- +goose Up
-- +goose StatementBegin
CREATE TABLE clicks (
    id BIGSERIAL PRIMARY KEY,
    link_id INT NOT NULL REFERENCES links (id) ON DELETE CASCADE,
    clicked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address TEXT,
    user_agent TEXT,
    referrer TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS clicks;
-- +goose StatementEnd
