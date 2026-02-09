CREATE OR REPLACE FUNCTION agora.unsave_content(
    p_user_id      UUID,
    p_content_id   UUID,
    p_content_type agora.content_type
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM agora.saved_content
    WHERE user_id = p_user_id
      AND content_id = p_content_id
      AND content_type = p_content_type;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'saved content not found'
            USING ERRCODE = 'P0002';
    END IF;
END;
$$;