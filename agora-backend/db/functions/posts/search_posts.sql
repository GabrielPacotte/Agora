CREATE OR REPLACE FUNCTION agora.search_posts(
    p_tags     TEXT[],
    p_from     TIMESTAMPTZ,
    p_to       TIMESTAMPTZ,
    p_limit    INT,
    p_offset   INT
)
RETURNS TABLE (
    id         UUID,
    author_id  UUID,
    title      TEXT,
    content    TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    is_invisible BOOLEAN,
    tags       TEXT[],
    stances    JSONB
)
LANGUAGE sql
AS $$
    SELECT
        p.id,
        p.author_id,
        p.title,
        p.content,
        p.created_at,
        p.updated_at,
        p.is_invisible,
        COALESCE((
            SELECT array_agg(pt.tag ORDER BY pt.tag)
            FROM agora.post_tags pt
            WHERE pt.post_id = p.id
        ), ARRAY[]::TEXT[]) AS tags,
        COALESCE((
            SELECT jsonb_agg(
                jsonb_build_object(
                    'id', s.id,
                    'label', s.label,
                    'color', s.color,
                    'description', s.description
                )
                ORDER BY s.id
            )
            FROM agora.post_stances s
            WHERE s.post_id = p.id
        ), '[]'::jsonb) AS stances
    FROM agora.posts p
    WHERE p.is_invisible = FALSE
      AND (p_from IS NULL OR p.created_at >= p_from)
      AND (p_to   IS NULL OR p.created_at <= p_to)
      AND (
            p_tags IS NULL
            OR cardinality(p_tags) = 0
            OR EXISTS (
                SELECT 1
                FROM agora.post_tags pt
                WHERE pt.post_id = p.id
                  AND pt.tag = ANY(p_tags)
            )
      )
    ORDER BY p.created_at DESC
    LIMIT p_limit
    OFFSET p_offset;
$$;