CREATE OR REPLACE FUNCTION agora.update_post(
    p_id           UUID,
    p_title        TEXT,
    p_content      TEXT,
    p_is_invisible BOOLEAN,
    p_updated_at   TIMESTAMPTZ,
    p_tags         TEXT[],
    p_stances      JSONB
)
RETURNS VOID
LANGUAGE plpgsql
AS $$
DECLARE
    v_count INT;
BEGIN
    IF p_stances IS NULL OR jsonb_typeof(p_stances) <> 'array' THEN
        RAISE EXCEPTION 'p_stances must be a JSON array'
        USING ERRCODE = 'AG001';
    END IF;

    SELECT COUNT(*)
    INTO v_count
    FROM jsonb_array_elements(p_stances);

    IF v_count < 3 THEN
        RAISE EXCEPTION 'A post must have at least 3 stances'
            USING ERRCODE = 'AG001';
    END IF;

    UPDATE agora.posts 
    SET
        title        = p_title,
        content      = p_content,
        updated_at   = p_updated_at,
        is_invisible = p_is_invisible
    WHERE id = p_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'post not found'
            USING ERRCODE = 'P0002';
    END IF;

    DELETE FROM agora.post_tags
    WHERE post_id = p_id;

    IF p_tags IS NOT NULL THEN
        INSERT INTO agora.tags(value)
        SELECT DISTINCT t FROM unnest(p_tags) AS t
        ON CONFLICT (value) DO NOTHING;

        INSERT INTO agora.post_tags (post_id, tag)
        SELECT DISTINCT p_id, t FROM unnest(p_tags) AS t;
    END IF;

    DELETE FROM agora.post_stances
    WHERE post_id = p_id;

    INSERT INTO agora.post_stances (id, post_id, label, color, description)
    SELECT
        (elem->>'id')::UUID,
        p_id,
        elem->>'label',
        elem->>'color',
        COALESCE(elem->>'description', '')
    FROM jsonb_array_elements(p_stances) AS elem;

END;
$$;