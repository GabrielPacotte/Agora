CREATE OR REPLACE FUNCTION agora.create_comment(
    p_id          UUID,
    p_post_id     UUID,
    p_author_id   UUID,
    p_stance_id   UUID,
    p_reply_to_id UUID,
    p_content     TEXT,
    p_created_at  TIMESTAMPTZ,
    p_updated_at  TIMESTAMPTZ
)
RETURNS VOID
LANGUAGE sql
AS $$
    INSERT INTO agora.comments (
        id, post_id, author_id, stance_id, reply_to_id,
        content, created_at, updated_at, is_invisible
    )
    VALUES (
        p_id, p_post_id, p_author_id, p_stance_id, p_reply_to_id,
        p_content, p_created_at, p_updated_at, FALSE
    );
$$;