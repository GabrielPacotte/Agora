CREATE OR REPLACE FUNCTION agora.get_session_by_refresh_token(
    p_refresh_token TEXT
)
RETURNS TABLE (
    id UUID,
    user_id UUID,
    created_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ
)
LANGUAGE sql
AS $$
    SELECT s.id, s.user_id, s.created_at, s.expires_at, s.revoked_at, s.last_used_at
    FROM agora.sessions s
    WHERE s.refresh_token_hash = digest(p_refresh_token, 'sha256')
      AND s.revoked_at IS NULL
      AND s.expires_at > NOW();
$$;