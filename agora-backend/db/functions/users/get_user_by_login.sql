CREATE OR REPLACE FUNCTION agora.get_user_by_login(
    p_login TEXT
)
RETURNS TABLE(
    id            UUID,
    login         TEXT,
    display_name  TEXT,
    password_hash TEXT,
    is_verified   BOOLEAN,
    created_at    TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ
)
LANGUAGE sql
AS $$
    SELECT id, login, display_name, password_hash, is_verified, created_at, updated_at
    FROM agora.users
    WHERE login = p_login;
$$;