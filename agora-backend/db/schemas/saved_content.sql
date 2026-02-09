CREATE TYPE agora.content_type AS ENUM ('post', 'comment');

CREATE TABLE IF NOT EXISTS agora.saved_content (
    user_id      UUID NOT NULL REFERENCES agora.users(id) ON DELETE CASCADE,
    content_id   UUID NOT NULL,
    content_type agora.content_type NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, content_id, content_type)
);

CREATE INDEX IF NOT EXISTS idx_saved_content_user ON agora.saved_content(user_id);