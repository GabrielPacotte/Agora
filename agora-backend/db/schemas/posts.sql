CREATE TABLE IF NOT EXISTS agora.posts (
    id           UUID PRIMARY KEY,
    author_id    UUID NOT NULL REFERENCES agora.users(id) ON DELETE CASCADE,
    title        TEXT NOT NULL,
    content      TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL,
    is_invisible BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS agora.post_stances (
    id          UUID PRIMARY KEY,
    post_id     UUID NOT NULL REFERENCES agora.posts(id) ON DELETE CASCADE,
    label       TEXT NOT NULL,
    color       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS agora.post_tags (
    post_id UUID NOT NULL REFERENCES agora.posts(id) ON DELETE CASCADE,
    tag     TEXT NOT NULL REFERENCES agora.tags(value) ON DELETE RESTRICT,
    PRIMARY KEY (post_id, tag)
);

CREATE INDEX IF NOT EXISTS idx_posts_author ON agora.posts(author_id);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON agora.posts(created_at);

CREATE TYPE agora.post_stance_input AS (
    id          UUID,
    label       TEXT,
    color       TEXT,
    description TEXT
);