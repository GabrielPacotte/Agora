CREATE OR REPLACE FUNCTION agora.update_user_public_fields(
    p_id            UUID,
    p_login         TEXT,
    p_display_name  TEXT,
    p_is_verified   BOOLEAN,
    p_updated_at    TIMESTAMPTZ
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE agora.users
    SET
        login         = p_login,
        display_name  = p_display_name,
        is_verified   = p_is_verified,
        updated_at    = p_updated_at
    WHERE id = p_id;
    
    IF NOT FOUND THEN
        RAISE EXCEPTION 'user not found'
            USING ERRCODE = 'P0002'; -- no_data_found
    END IF;
END
$$;