CREATE TABLE IF NOT EXISTS agora.reports (
    id           UUID PRIMARY KEY,
    author_id    UUID NOT NULL REFERENCES agora.users(id) ON DELETE CASCADE,
    post_subject_id UUID REFERENCES agora.posts(id) ON DELETE CASCADE,
    comment_subject_id UUID REFERENCES agora.comments(id) ON DELETE CASCADE,
    subject_type agora.content_type NOT NULL,
    description  TEXT NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL,
    updated_at   TIMESTAMPTZ NOT NULL,

    CHECK (
        (post_subject_id IS NOT NULL AND comment_subject_id IS NULL) OR
        (post_subject_id IS NULL AND comment_subject_id IS NOT NULL)
    )
);
