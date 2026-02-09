CREATE OR REPLACE FUNCTION agora.create_user(
    p_id            UUID,
    p_login         TEXT,
    p_display_name  TEXT,
    p_password_hash TEXT,
    p_is_verified   BOOLEAN,
    p_created_at    TIMESTAMPTZ,
    p_updated_at    TIMESTAMPTZ
)
RETURNS VOID
LANGUAGE sql
AS $$
    INSERT INTO agora.users (
        id, login, display_name, password_hash, is_verified, created_at, updated_at
    )
    VALUES (
        p_id, p_login, p_display_name, p_password_hash, p_is_verified, p_created_at, p_updated_at
    );
$$;