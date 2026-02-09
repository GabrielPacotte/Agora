CREATE TABLE IF NOT EXISTS agora.user_preferences (
    user_id UUID NOT NULL REFERENCES agora.users(id) ON DELETE CASCADE,
    tag     TEXT NOT NULL REFERENCES agora.tags(value) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, tag)
);

CREATE INDEX IF NOT EXISTS idx_user_prefs_user ON agora.user_preferences(user_id);