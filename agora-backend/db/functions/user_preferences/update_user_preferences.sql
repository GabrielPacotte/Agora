CREATE OR REPLACE FUNCTION agora.update_user_preferences(
    p_user_id UUID,
    p_tags    TEXT[]
)
RETURNS void
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM agora.user_preferences
    WHERE user_id = p_user_id;
    
    IF p_tags IS NOT NULL THEN
        INSERT INTO agora.tags(value)
        SELECT DISTINCT t FROM unnest(p_tags) AS t
        ON CONFLICT (value) DO NOTHING;

        INSERT INTO agora.user_preferences (user_id, tag)
        SELECT DISTINCT p_user_id, t FROM unnest(p_tags) AS t;
    END IF;
END;
$$;