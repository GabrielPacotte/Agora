CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS agora.sessions (
    id                  UUID PRIMARY KEY,
    user_id             UUID NOT NULL REFERENCES agora.users(id) ON DELETE CASCADE,

    refresh_token_hash  BYTEA NOT NULL UNIQUE,

    created_at          TIMESTAMPTZ NOT NULL,
    expires_at          TIMESTAMPTZ NOT NULL,

    revoked_at          TIMESTAMPTZ NULL,
    last_used_at        TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON agora.sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON agora.sessions(expires_at);
