CREATE OR REPLACE FUNCTION agora.create_session(
    p_id                UUID,
    p_user_id           UUID,
    p_refresh_token     TEXT,
    p_created_at        TIMESTAMPTZ,
    p_expires_at        TIMESTAMPTZ
)
RETURNS VOID
LANGUAGE sql
AS $$
    INSERT INTO agora.sessions (
        id, user_id, refresh_token_hash, created_at, expires_at, revoked_at, last_used_at
    )
    VALUES (
        p_id,
        p_user_id,
        digest(p_refresh_token, 'sha256'),
        p_created_at,
        p_expires_at,
        NULL,
        NULL
    );
$$;