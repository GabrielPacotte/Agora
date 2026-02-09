CREATE OR REPLACE FUNCTION agora.update_comment(
    p_id          UUID,
    p_stance_id   UUID,
    p_reply_to_id UUID,
    p_content     TEXT,
    p_updated_at  TIMESTAMPTZ
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE agora.comments
    SET stance_id   = p_stance_id,
        reply_to_id = p_reply_to_id,
        content     = p_content,
        updated_at  = p_updated_at
    WHERE id = p_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'comment not found'
            USING ERRCODE = 'P0002';
    END IF;
END;
$$;