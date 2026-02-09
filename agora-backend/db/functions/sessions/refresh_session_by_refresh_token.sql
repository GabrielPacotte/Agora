CREATE OR REPLACE FUNCTION agora.revoke_session_by_refresh_token(
    p_refresh_token TEXT,
    p_revoked_at TIMESTAMPTZ
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE agora.sessions
    SET revoked_at = p_revoked_at
    WHERE refresh_token_hash = digest(p_refresh_token, 'sha256')
      AND revoked_at IS NULL;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'session not found' USING ERRCODE = 'P0002';
    END IF;
END;
$$;