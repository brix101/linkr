-- name: GetLinkByCode :one
SELECT id, code, url, created_at FROM links WHERE code = $1 LIMIT 1;

-- name: GetLinkByURL :one
SELECT id, code, url, created_at FROM LINKS WHERE url = $1 LIMIT 1;

-- name: CreateLink :one
INSERT INTO links (
    code,
    url,
    expires_at,
    user_id
)
VALUES ($1, $2, $3, $4)
RETURNING id, code, url, created_at;


-- name: CreateClick :exec
INSERT INTO clicks (link_id, ip_address, user_agent, referrer)
VALUES ($1, $2, $3, $4);
