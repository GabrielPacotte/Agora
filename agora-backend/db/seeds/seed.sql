-- ================================================
-- SEED DATA FOR KAINE AGORA (development only)
-- ================================================

SET search_path TO agora;

-- ============================
-- 1. TAGS
-- ============================

INSERT INTO tags (value) VALUES
    ('politics'),
    ('economy'),
    ('videogames'),
    ('france'),
    ('science'),
    ('opinion');


-- ============================
-- 2. USERS
-- ============================

INSERT INTO users (id, login, display_name, password_hash, is_verified, created_at, updated_at)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'alice', 'Alice Dupont', 'HASH1', true, NOW(), NOW()),
    ('22222222-2222-2222-2222-222222222222', 'bob',   'Bob Martin',   'HASH2', true, NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333333', 'charlie', 'Charlie Doe', 'HASH3', false, NOW(), NOW());


-- ============================
-- 3. USER PREFERENCES
-- ============================

INSERT INTO user_preferences (user_id, tag)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'politics'),
    ('11111111-1111-1111-1111-111111111111', 'france'),
    ('22222222-2222-2222-2222-222222222222', 'videogames'),
    ('22222222-2222-2222-2222-222222222222', 'science');


-- ============================
-- 4. POSTS
-- ============================

INSERT INTO posts (id, author_id, title, content, created_at, updated_at)
VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
        '11111111-1111-1111-1111-111111111111',
        'What do you think of the 2025 French budget?',
        'Here is a summary of the proposed budget...',
        NOW() - INTERVAL '2 days',
        NOW() - INTERVAL '1 day'
    ),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
        '22222222-2222-2222-2222-222222222222',
        'GTA 7 just announced!',
        'Rockstar finally dropped the official trailer...',
        NOW() - INTERVAL '5 hours',
        NOW() - INTERVAL '3 hours'
    );


-- ============================
-- 5. POST TAGS
-- ============================

INSERT INTO post_tags (post_id, tag)
VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'politics'),
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'france'),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'videogames'),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'opinion');


-- ============================
-- 6. POST STANCES
-- ============================

INSERT INTO post_stances (id, post_id, label, color, description) VALUES
    ('aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
        'Agree', '#00ff00', ''),
    ('aaaa2222-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
        'Disagree', '#ff0000', ''),
    ('aaaa3333-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
        'Need more info', '#0000ff', ''),

    ('bbbb1111-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
        'Hyped', '#00ff00', 'Very excited'),
    ('bbbb2222-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
        'Not convinced', '#ff0000', ''),
    ('bbbb3333-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
        'Waiting for gameplay', '#0000ff', '');


-- ============================
-- 7. COMMENTS
-- ============================

INSERT INTO comments (id, post_id, author_id, stance_id, reply_to_id, content, created_at, updated_at)
VALUES
    ('cccc1111-cccc-cccc-cccc-cccccccccccc',
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
        '22222222-2222-2222-2222-222222222222',
        'aaaa2222-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
        NULL,
        'I strongly disagree with this budget.',
        NOW() - INTERVAL '1 day',
        NOW() - INTERVAL '1 day'
    ),
    ('cccc2222-cccc-cccc-cccc-cccccccccccc',
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
        '11111111-1111-1111-1111-111111111111',
        'aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
        'cccc1111-cccc-cccc-cccc-cccccccccccc',
        'I understand your concerns, but here’s some context...',
        NOW() - INTERVAL '20 hours',
        NOW() - INTERVAL '20 hours'
    );


-- ============================
-- 8. SAVED CONTENT
-- ============================

INSERT INTO saved_content (user_id, content_id, content_type)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'post'),
    ('22222222-2222-2222-2222-222222222222', 'cccc1111-cccc-cccc-cccc-cccccccccccc', 'comment');


-- ============================
-- 9. REPORTS
-- ============================

INSERT INTO reports (id, author_id, post_subject_id, comment_subject_id, subject_type, description, created_at, updated_at)
VALUES
    ('dddd1111-dddd-dddd-dddd-dddddddddddd',
        '33333333-3333-3333-3333-333333333333',
        NULL,
        'cccc1111-cccc-cccc-cccc-cccccccccccc',
        'comment',
        'This comment is inappropriate.',
        NOW() - INTERVAL '10 hours',
        NOW() - INTERVAL '10 hours'
    ), ('eeee1111-eeee-eeee-eeee-eeeeeeeeeeee',
        '33333333-3333-3333-3333-333333333333',
        'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
        NULL,
        'post',
        'This post is inappropriate.',
        NOW() - INTERVAL '8 hours',
        NOW() - INTERVAL '8 hours'
    );;


-- ============================
-- 10. MODERATION ACTIONS
-- ============================

INSERT INTO moderation_actions (id, moderator_id, decision, description, created_at, updated_at)
VALUES
    ('eeee1111-eeee-eeee-eeee-eeeeeeeeeeee',
        '11111111-1111-1111-1111-111111111111',
        'report_reviewed',
        'Reviewed and accepted.',
        NOW() - INTERVAL '5 hours',
        NOW() - INTERVAL '5 hours'
    );

INSERT INTO moderation_action_reports (action_id, report_id)
VALUES
    ('eeee1111-eeee-eeee-eeee-eeeeeeeeeeee', 'dddd1111-dddd-dddd-dddd-dddddddddddd');
