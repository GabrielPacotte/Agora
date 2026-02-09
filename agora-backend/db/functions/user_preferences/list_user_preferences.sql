CREATE OR REPLACE FUNCTION agora.list_user_preferences(
    p_user_id UUID
)
RETURNS TABLE(tag TEXT)
LANGUAGE sql
AS $$
    SELECT tag
    FROM agora.user_preferences
    WHERE user_id = p_user_id;
$$;