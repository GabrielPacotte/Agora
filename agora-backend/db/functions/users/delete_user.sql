CREATE OR REPLACE FUNCTION agora.delete_user(p_id UUID)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM agora.users
    WHERE id = p_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'user not found'
            USING ERRCODE = 'P0002';
    END IF;
END;
$$;