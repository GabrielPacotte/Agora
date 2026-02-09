CREATE OR REPLACE FUNCTION agora.delete_post(p_id UUID)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM agora.posts
    WHERE id = p_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'post not found'
        USING ERRCODE = 'P0002';
    END IF;
END;
$$;