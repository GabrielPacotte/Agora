package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepositoryPG struct {
	db *pgxpool.Pool
}

func NewCommentRepositoryPG(db *pgxpool.Pool) *CommentRepositoryPG {
	return &CommentRepositoryPG{db: db}
}

func (r *CommentRepositoryPG) Create(ctx context.Context, comment *domain.Comment) error {
	cmd := `
		SELECT agora.create_comment(
			$1, $2, $3, $4, $5, $6, $7, $8
		);
	`

	var replyTo any
	if comment.ReplyToID() != nil {
		replyTo = *comment.ReplyToID()
	} else {
		replyTo = nil
	}

	_, err := r.db.Exec(ctx, cmd,
		comment.ID(),
		comment.PostID(),
		comment.AuthorID(),
		comment.Stance().ID(),
		replyTo,
		comment.Content(),
		comment.CreatedAt(),
		comment.UpdatedAt(),
	)
	return wrapPGError(err)
}

func (r *CommentRepositoryPG) Update(ctx context.Context, comment *domain.Comment) error {
	cmd := `
		SELECT agora.update_comment(
			$1, $2, $3, $4, $5
		);
	`

	var replyTo any
	if comment.ReplyToID() != nil {
		replyTo = *comment.ReplyToID()
	} else {
		replyTo = nil
	}

	_, err := r.db.Exec(ctx, cmd,
		comment.ID(),
		comment.Stance().ID(),
		replyTo,
		comment.Content(),
		comment.UpdatedAt(),
	)
	return wrapPGError(err)
}

func (r *CommentRepositoryPG) Delete(ctx context.Context, commentID uuid.UUID) error {
	if commentID == uuid.Nil {
		return domain.ErrInvalidID
	}

	cmd := `SELECT agora.delete_comment($1);`
	_, err := r.db.Exec(ctx, cmd, commentID)
	return wrapPGError(err)
}

func (r *CommentRepositoryPG) GetByID(
	ctx context.Context,
	commentID uuid.UUID,
) (*domain.Comment, error) {
	cmd := `
        SELECT
            id,
            post_id,
            author_id,
            stance_id,
            stance_label,
            stance_color,
            stance_description,
            reply_to_id,
            content,
            created_at,
            updated_at
        FROM agora.get_comment_by_id($1);
    `
	row := r.db.QueryRow(ctx, cmd, commentID)
	return scanComment(row)
}

func (r *CommentRepositoryPG) ListRootComments(
	ctx context.Context,
	postID uuid.UUID,
	page filters.Page,
) ([]domain.Comment, error) {
	cmd := `
		SELECT
			id,
			post_id,
			author_id,
			stance_id,
			stance_label,
			stance_color,
			stance_description,
			reply_to_id,
			content,
			created_at,
			updated_at
		FROM agora.list_root_comments($1, $2, $3);
	`

	rows, err := r.db.Query(ctx, cmd, postID, page.Limit(), page.Offset())
	if err != nil {
		return nil, wrapPGError(err)
	}
	defer rows.Close()

	var out []domain.Comment

	for rows.Next() {
		c, err := scanComment(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}

	if rows.Err() != nil {
		return nil, wrapPGError(rows.Err())
	}

	return out, nil
}

func (r *CommentRepositoryPG) ListReplies(
	ctx context.Context,
	parentCommentID uuid.UUID,
	page filters.Page,
) ([]domain.Comment, error) {
	cmd := `
		SELECT
			id,
			post_id,
			author_id,
			stance_id,
			stance_label,
			stance_color,
			stance_description,
			reply_to_id,
			content,
			created_at,
			updated_at
		FROM agora.list_replies($1, $2, $3);
	`

	rows, err := r.db.Query(ctx, cmd, parentCommentID, page.Limit(), page.Offset())
	if err != nil {
		return nil, wrapPGError(err)
	}
	defer rows.Close()

	var out []domain.Comment

	for rows.Next() {
		c, err := scanComment(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}

	if rows.Err() != nil {
		return nil, wrapPGError(rows.Err())
	}

	return out, nil
}

func scanComment(row pgx.Row) (*domain.Comment, error) {
	var (
		id          uuid.UUID
		postID      uuid.UUID
		authorID    uuid.UUID
		stanceID    uuid.UUID
		stanceLabel string
		stanceColor string
		stanceDesc  string
		replyToID   *uuid.UUID
		content     string
		createdAt   time.Time
		updatedAt   time.Time
	)

	if err := row.Scan(
		&id,
		&postID,
		&authorID,
		&stanceID,
		&stanceLabel,
		&stanceColor,
		&stanceDesc,
		&replyToID,
		&content,
		&createdAt,
		&updatedAt,
	); err != nil {
		return nil, wrapPGError(err)
	}

	color, err := domain.NewColor(stanceColor)
	if err != nil {
		return nil, fmt.Errorf("invalid stance color from DB: %w", err)
	}

	stance, err := domain.NewStance(stanceID, stanceLabel, color, stanceDesc)
	if err != nil {
		return nil, fmt.Errorf("invalid stance from DB: %w", err)
	}

	comment, err := domain.ExistingComment(
		id,
		createdAt,
		updatedAt,
		content,
		stance,
		postID,
		authorID,
		replyToID,
	)
	if err != nil {
		return nil, err
	}

	return comment, nil
}
