CREATE OR REPLACE FUNCTION agora.get_user_by_id(
    p_id UUID
)
RETURNS TABLE (
    id           UUID,
    login        TEXT,
    display_name TEXT,
    created_at   TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ
)
LANGUAGE sql
AS $$
    SELECT id, login, display_name, created_at, updated_at
    FROM agora.users
    WHERE id = p_id;
$$;