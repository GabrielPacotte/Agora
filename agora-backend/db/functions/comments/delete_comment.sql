CREATE OR REPLACE FUNCTION agora.delete_comment(
    p_id UUID
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM agora.comments
    WHERE id = p_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'comment not found'
            USING ERRCODE = 'P0002';
    END IF;
END;
$$;