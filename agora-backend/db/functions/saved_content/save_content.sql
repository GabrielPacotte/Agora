CREATE OR REPLACE FUNCTION agora.save_content(
    p_user_id      UUID,
    p_content_id   UUID,
    p_content_type agora.content_type
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    IF p_content_type = 'post' THEN
        IF NOT EXISTS (
            SELECT 1 FROM agora.posts WHERE id = p_content_id AND is_invisible = FALSE
        ) THEN
            RAISE EXCEPTION 'Post not found or not visible'
                USING ERRCODE = '23503';
        END IF;
    ELSIF p_content_type = 'comment' THEN
        IF NOT EXISTS (
            SELECT 1 FROM agora.comments WHERE id = p_content_id AND is_invisible = FALSE
        ) THEN
            RAISE EXCEPTION 'Comment not found or not visible'
                USING ERRCODE = '23503';
        END IF;
    ELSE
        RAISE EXCEPTION 'Invalid content type %', p_content_type;
    END IF;

    INSERT INTO agora.saved_content (user_id, content_id, content_type)
    VALUES (p_user_id, p_content_id, p_content_type)
    ON CONFLICT DO NOTHING;
END;
$$;