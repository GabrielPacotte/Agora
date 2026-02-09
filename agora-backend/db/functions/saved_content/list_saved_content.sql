CREATE OR REPLACE FUNCTION agora.list_saved_content(
    p_user_id UUID,
    p_limit   INT,
    p_offset  INT
)
RETURNS TABLE (
    user_id      UUID,
    content_id   UUID,
    content_type agora.content_type,
    created_at   TIMESTAMPTZ
)
LANGUAGE sql
AS $$
    SELECT user_id, content_id, content_type, created_at
    FROM agora.saved_content
    WHERE user_id = p_user_id
    ORDER BY created_at DESC
    LIMIT p_limit
    OFFSET p_offset;
$$;