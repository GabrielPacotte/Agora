CREATE TABLE IF NOT EXISTS agora.comments (
    id           UUID PRIMARY KEY,
    post_id      UUID NOT NULL REFERENCES agora.posts(id) ON DELETE CASCADE,
    author_id    UUID NOT NULL REFERENCES agora.users(id) ON DELETE CASCADE,
    stance_id    UUID NOT NULL REFERENCES agora.post_stances(id) ON DELETE RESTRICT,
    reply_to_id  UUID NULL,
    content      TEXT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL,
    is_invisible BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT comments_id_post_unique UNIQUE (id, post_id),
    CONSTRAINT comments_reply_same_post_fk
        FOREIGN KEY (reply_to_id, post_id)
        REFERENCES agora.comments (id, post_id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_comments_post ON agora.comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_reply_to ON agora.comments(reply_to_id);
CREATE INDEX IF NOT EXISTS idx_comments_author ON agora.comments(author_id);