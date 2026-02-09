CREATE OR REPLACE FUNCTION agora.touch_session_last_used(
    p_id UUID,
    p_last_used_at TIMESTAMPTZ
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE agora.sessions
    SET last_used_at = p_last_used_at
    WHERE id = p_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'session not found' USING ERRCODE = 'P0002';
    END IF;
END;
$$;