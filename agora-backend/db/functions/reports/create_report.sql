CREATE OR REPLACE FUNCTION agora.create_report(
    p_id UUID,
    p_author_id UUID,
    p_subject_id UUID,
    p_subject_type agora.content_type,
    p_description TEXT,
    p_created_at TIMESTAMPTZ,
    p_updated_at TIMESTAMPTZ
)
RETURNS VOID
LANGUAGE sql
AS $$
    -- TODO
$$;