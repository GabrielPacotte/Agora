package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/GabrielPacotte/Agora/internal/domain"
	"github.com/GabrielPacotte/Agora/internal/filters"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostRepositoryPG struct {
	db *pgxpool.Pool
}

func NewPostRepositoryPG(db *pgxpool.Pool) *PostRepositoryPG {
	return &PostRepositoryPG{db: db}
}

type stanceRecord struct {
	ID          uuid.UUID `json:"id"`
	Label       string    `json:"label"`
	Color       string    `json:"color"`
	Description string    `json:"description"`
}

func (r *PostRepositoryPG) Create(ctx context.Context, post *domain.Post) error {
	tagVals := make([]string, 0)
	for _, t := range post.Tags() {
		tagVals = append(tagVals, t.String())
	}

	stances := post.Stances()
	jsonStances := make([]stanceRecord, 0, len(stances))
	for _, s := range stances {
		jsonStances = append(jsonStances, stanceRecord{
			ID:          s.ID(),
			Label:       s.Label(),
			Color:       s.Color().String(),
			Description: s.Description(),
		})
	}

	stBytes, err := json.Marshal(jsonStances)
	if err != nil {
		return fmt.Errorf("marshal stances: %w", err)
	}

	cmd := `
		SELECT agora.create_post(
			$1, $2, $3, $4, $5, $6, $7, $8
		);
	`

	_, err = r.db.Exec(ctx, cmd,
		post.ID(),
		post.AuthorID(),
		post.Title(),
		post.Content(),
		post.CreatedAt(),
		post.UpdatedAt(),
		tagVals,
		stBytes,
	)
	return wrapPGError(err)
}

func (r *PostRepositoryPG) Update(ctx context.Context, post *domain.Post) error {
	tagVals := make([]string, 0)
	for _, t := range post.Tags() {
		tagVals = append(tagVals, t.String())
	}

	stances := post.Stances()
	jsonStances := make([]stanceRecord, 0, len(stances))
	for _, s := range stances {
		jsonStances = append(jsonStances, stanceRecord{
			ID:          s.ID(),
			Label:       s.Label(),
			Color:       s.Color().String(),
			Description: s.Description(),
		})
	}

	stBytes, err := json.Marshal(jsonStances)
	if err != nil {
		return fmt.Errorf("marshal stances: %w", err)
	}

	const isInvisible = false

	cmd := `
		SELECT agora.update_post(
			$1, $2, $3, $4, $5, $6, $7
		);
	`

	_, err = r.db.Exec(ctx, cmd,
		post.ID(),
		post.Title(),
		post.Content(),
		isInvisible,
		post.UpdatedAt(),
		tagVals,
		stBytes,
	)
	return wrapPGError(err)
}

func (r *PostRepositoryPG) Delete(ctx context.Context, postID uuid.UUID) error {
	cmd := `SELECT agora.delete_post($1);`
	_, err := r.db.Exec(ctx, cmd, postID)
	return wrapPGError(err)
}

func (r *PostRepositoryPG) GetByID(ctx context.Context, postID uuid.UUID) (*domain.Post, error) {
	cmd := `
		SELECT
			id,
			author_id,
			title,
			content,
			created_at,
			updated_at,
			is_invisible,
			tags,
			stances
		FROM agora.get_post_by_id($1);
	`

	row := r.db.QueryRow(ctx, cmd, postID)
	return scanPost(row)
}

func (r *PostRepositoryPG) ListByUserID(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.Post, error) {
	cmd := `
		SELECT
			id,
			author_id,
			title,
			content,
			created_at,
			updated_at,
			is_invisible,
			tags,
			stances
		FROM agora.list_posts_by_user($1, $2, $3);
	`

	rows, err := r.db.Query(ctx, cmd, userID, page.Limit(), page.Offset())
	if err != nil {
		return nil, wrapPGError(err)
	}
	defer rows.Close()

	var posts []domain.Post

	for rows.Next() {
		p, err := scanPost(rows)
		if err != nil {
			return nil, err
		}
		posts = append(posts, *p)
	}

	if rows.Err() != nil {
		return nil, wrapPGError(rows.Err())
	}

	return posts, nil
}

func (r *PostRepositoryPG) GetFeed(ctx context.Context, userID uuid.UUID, page filters.Page) ([]domain.Post, error) {
	cmd := `
		SELECT
			id,
			author_id,
			title,
			content,
			created_at,
			updated_at,
			is_invisible,
			tags,
			stances
		FROM agora.get_feed_posts($1, $2, $3);
	`

	rows, err := r.db.Query(ctx, cmd, userID, page.Limit(), page.Offset())
	if err != nil {
		return nil, wrapPGError(err)
	}
	defer rows.Close()

	out := make([]domain.Post, 0, page.Limit())
	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *post)
	}

	if err := rows.Err(); err != nil {
		return nil, wrapPGError(err)
	}

	return out, nil
}

func (r *PostRepositoryPG) Search(ctx context.Context, f filters.SearchPostFilter) ([]domain.Post, error) {
	var tags []string
	if len(f.Tags) > 0 {
		tags = make([]string, 0, len(f.Tags))
		for _, t := range f.Tags {
			tags = append(tags, t.String())
		}
	}

	cmd := `
		SELECT
			id,
			author_id,
			title,
			content,
			created_at,
			updated_at,
			is_invisible,
			tags,
			stances
		FROM agora.search_posts($1, $2, $3, $4, $5);
	`

	rows, err := r.db.Query(
		ctx,
		cmd,
		tags,
		f.FromDate,
		f.ToDate,
		f.Page.Limit(),
		f.Page.Offset(),
	)
	if err != nil {
		return nil, wrapPGError(err)
	}
	defer rows.Close()

	out := make([]domain.Post, 0, f.Page.Limit())
	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *post)
	}

	if err := rows.Err(); err != nil {
		return nil, wrapPGError(err)
	}

	return out, nil
}

func scanPost(row pgx.Row) (*domain.Post, error) {
	var (
		id          uuid.UUID
		authorID    uuid.UUID
		title       string
		content     string
		createdAt   time.Time
		updatedAt   time.Time
		isInvisible bool
		tagArray    []string
		stanceJSON  []byte
	)

	if err := row.Scan(
		&id,
		&authorID,
		&title,
		&content,
		&createdAt,
		&updatedAt,
		&isInvisible,
		&tagArray,
		&stanceJSON,
	); err != nil {
		return nil, wrapPGError(err)
	}

	tags := make([]domain.Tag, 0, len(tagArray))
	for _, t := range tagArray {
		tags = append(tags, domain.Tag(t))
	}

	var rawStances []stanceRecord
	if len(stanceJSON) > 0 {
		if err := json.Unmarshal(stanceJSON, &rawStances); err != nil {
			return nil, fmt.Errorf("decode stances JSON: %w", err)
		}
	}

	stances := make([]domain.Stance, 0, len(rawStances))
	for _, s := range rawStances {
		color, err := domain.NewColor(s.Color)
		if err != nil {
			return nil, fmt.Errorf("invalid color from DB: %w", err)
		}
		stance, err := domain.NewStance(s.ID, s.Label, color, s.Description)
		if err != nil {
			return nil, fmt.Errorf("invalid stance from DB: %w", err)
		}
		stances = append(stances, stance)
	}

	post, err := domain.ExistingPost(
		id,
		createdAt,
		updatedAt,
		title,
		content,
		tags,
		stances,
		authorID,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}
