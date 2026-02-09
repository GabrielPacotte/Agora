CREATE OR REPLACE FUNCTION agora.get_post_by_id(
    p_id UUID
)
RETURNS TABLE (
    id UUID,
    author_id UUID,
    title TEXT,
    content TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    is_invisible BOOLEAN,
    tags TEXT[],
    stances JSONB
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
    (
        SELECT COALESCE(array_agg(pt.tag), '{}')
        FROM agora.post_tags pt
        WHERE pt.post_id = p.id
    ) AS tags,
    (
        SELECT COALESCE(
            jsonb_agg(
                jsonb_build_object(
                    'id', ps.id,
                    'label', ps.label,
                    'color', ps.color,
                    'description', ps.description
                )
            ),
            '[]'::jsonb
        )
        FROM agora.post_stances ps
        WHERE ps.post_id = p.id
    ) AS stances
FROM agora.posts p
WHERE p.id = p_id;
$$;