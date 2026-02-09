CREATE OR REPLACE FUNCTION agora.update_user_password(
    p_id UUID,
    p_password_hash TEXT,
    p_updated_at TIMESTAMPTZ
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE agora.users
    SET password_hash = p_password_hash,
        updated_at = p_updated_at
    WHERE id = p_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'user not found'
            USING ERRCODE = 'P0002';
    END IF;
END
$$;