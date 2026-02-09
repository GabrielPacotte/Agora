CREATE OR REPLACE FUNCTION agora.create_post(
    p_id          UUID,
    p_author_id   UUID,
    p_title       TEXT,
    p_content     TEXT,
    p_created_at  TIMESTAMPTZ,
    p_updated_at  TIMESTAMPTZ,
    p_tags        TEXT[],
    p_stances     JSONB
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

    INSERT INTO agora.posts (
        id, author_id, title, content, created_at, updated_at, is_invisible
    )
    VALUES (
        p_id, p_author_id, p_title, p_content, p_created_at, p_updated_at, FALSE
    );

    IF p_tags IS NOT NULL THEN
        INSERT INTO agora.tags(value)
        SELECT DISTINCT t FROM unnest(p_tags) AS t
        ON CONFLICT (value) DO NOTHING;

        INSERT INTO agora.post_tags (post_id, tag)
        SELECT DISTINCT p_id, t FROM unnest(p_tags) AS t;
    END IF;

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