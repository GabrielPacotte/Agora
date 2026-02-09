CREATE TYPE agora.moderation_decision AS ENUM (
    'pending',
    'report_reviewed',
    'content_deleted',
    'user_suspended',
    'user_banished'
);

CREATE TABLE IF NOT EXISTS agora.moderation_actions (
    id          UUID PRIMARY KEY,
    moderator_id UUID NOT NULL REFERENCES agora.users(id) ON DELETE RESTRICT,
    decision    agora.moderation_decision NOT NULL,
    description TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS agora.moderation_action_reports (
    action_id UUID NOT NULL REFERENCES agora.moderation_actions(id) ON DELETE CASCADE,
    report_id UUID NOT NULL REFERENCES agora.reports(id) ON DELETE CASCADE,
    PRIMARY KEY (action_id, report_id)
);

CREATE INDEX IF NOT EXISTS idx_moderation_actions_mod ON agora.moderation_actions(moderator_id);