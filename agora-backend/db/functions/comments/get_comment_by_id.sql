CREATE OR REPLACE FUNCTION agora.get_comment_by_id(
    p_id UUID
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
        s.id,
        s.label,
        s.color,
        s.description,
        c.reply_to_id,
        c.content,
        c.created_at,
        c.updated_at
    FROM agora.comments c
    JOIN agora.post_stances s
        ON s.id = c.stance_id
    WHERE c.id = p_id;
$$;