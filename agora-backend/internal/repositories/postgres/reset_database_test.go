package postgres_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestMain(m *testing.M) {
	dsn := os.Getenv("KAINE_TEST_DATABASE_URL")
	if dsn == "" {
		os.Exit(m.Run())
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	testPool = pool

	if err := resetTestDatabase(ctx, testPool); err != nil {
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}

func resetTestDatabase(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `DROP SCHEMA IF EXISTS agora CASCADE; CREATE SCHEMA agora;`)
	if err != nil {
		return err
	}

	schemaFiles := []string{
		"db/schemas/common.sql",
		"db/schemas/users.sql",
		"db/schemas/user_preferences.sql",
		"db/schemas/posts.sql",
		"db/schemas/comments.sql",
		"db/schemas/saved_content.sql",
		"db/schemas/reports.sql",
		"db/schemas/moderation_actions.sql",
		"db/schemas/sessions.sql",
		"db/functions/comments/create_comment.sql",
		"db/functions/comments/delete_comment.sql",
		"db/functions/comments/get_comment_by_id.sql",
		"db/functions/comments/list_replies.sql",
		"db/functions/comments/list_root_comments.sql",
		"db/functions/comments/create_comment.sql",
		"db/functions/comments/update_comment.sql",
		"db/functions/moderation_actions/create_moderation_action.sql",
		"db/functions/moderation_actions/create_moderation_actions_against_user.sql",
		"db/functions/moderation_actions/list_moderation_actions.sql",
		"db/functions/moderation_actions/update_moderation_action.sql",
		"db/functions/posts/create_post.sql",
		"db/functions/posts/delete_post.sql",
		"db/functions/posts/get_feed_posts.sql",
		"db/functions/posts/get_post_by_id.sql",
		"db/functions/posts/list_posts_by_user.sql",
		"db/functions/posts/search_posts.sql",
		"db/functions/posts/update_post.sql",
		"db/functions/reports/create_report.sql",
		"db/functions/reports/get_report_by_id.sql",
		"db/functions/reports/list_reports_against_user.sql",
		"db/functions/reports/list_reports.sql",
		"db/functions/saved_content/list_saved_content.sql",
		"db/functions/saved_content/save_content.sql",
		"db/functions/saved_content/unsave_content.sql",
		"db/functions/sessions/create_session.sql",
		"db/functions/sessions/get_session_by_refresh_token.sql",
		"db/functions/sessions/refresh_session_by_refresh_token.sql",
		"db/functions/sessions/touch_session_last_used.sql",
		"db/functions/user_preferences/list_user_preferences.sql",
		"db/functions/user_preferences/update_user_preferences.sql",
		"db/functions/users/create_user.sql",
		"db/functions/users/delete_user.sql",
		"db/functions/users/get_user_by_id.sql",
		"db/functions/users/get_user_by_login.sql",
		"db/functions/users/update_user_password.sql",
		"db/functions/users/update_user_public_fields.sql",
	}

	for _, path := range schemaFiles {
		if err := applySQLFile(ctx, pool, path); err != nil {
			return err
		}
	}

	funcFiles, err := filepath.Glob("db/functions/*/*.sql")
	if err != nil {
		return err
	}
	for _, path := range funcFiles {
		if err := applySQLFile(ctx, pool, path); err != nil {
			return err
		}
	}

	if err := applySQLFile(ctx, pool, "db/seeds/seed.sql"); err != nil {
		return err
	}

	return nil
}
