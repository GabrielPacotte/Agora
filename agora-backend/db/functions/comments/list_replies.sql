CREATE OR REPLACE FUNCTION agora.list_replies(
    p_comment_id UUID,
    p_limit      INT,
    p_offset     INT
)
RETURNS TABLE (
    id UUID,
    post_id UUID,
    author_id UUID,
    stance_id UUID,
    stance_label TEXT,
    stance_color TEXT,
    stance_description TEXT,
    reply_to_id UUID,
    content TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
)
LANGUAGE sql
AS $$
    SELECT
        c.id,
        c.post_id,
        c.author_id,
        s.id AS stance_id,
        s.label AS stance_label,
        s.color AS stance_color,
        s.description AS stance_description,
        c.reply_to_id,
        c.content,
        c.created_at,
        c.updated_at
    FROM agora.comments c
    JOIN agora.post_stances s ON c.stance_id = s.id
    WHERE c.reply_to_id = p_comment_id
      AND c.is_invisible = FALSE
    ORDER BY c.created_at ASC
    LIMIT p_limit
    OFFSET p_offset;
$$;