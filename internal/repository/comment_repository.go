package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/user/kanban-saas/pkg/model"
)

type CommentRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(ctx context.Context, comment *model.Comment) error {
	query := `
		INSERT INTO comments (id, card_id, user_id, content)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		comment.ID, comment.CardID, comment.UserID, comment.Content,
	).Scan(&comment.CreatedAt, &comment.UpdatedAt)
}

func (r *CommentRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error) {
	query := `
		SELECT id, card_id, user_id, content, created_at, updated_at
		FROM comments WHERE id = $1 AND deleted_at IS NULL`

	comment := &model.Comment{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&comment.ID, &comment.CardID, &comment.UserID, &comment.Content,
		&comment.CreatedAt, &comment.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get comment: %w", err)
	}
	return comment, nil
}

func (r *CommentRepository) ListByCard(ctx context.Context, cardID uuid.UUID) ([]model.Comment, error) {
	query := `
		SELECT id, card_id, user_id, content, created_at, updated_at
		FROM comments WHERE card_id = $1 AND deleted_at IS NULL
		ORDER BY created_at ASC`

	rows, err := r.db.Query(ctx, query, cardID)
	if err != nil {
		return nil, fmt.Errorf("list comments: %w", err)
	}
	defer rows.Close()

	var comments []model.Comment
	for rows.Next() {
		var comment model.Comment
		if err := rows.Scan(&comment.ID, &comment.CardID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *CommentRepository) Update(ctx context.Context, comment *model.Comment) error {
	query := `
		UPDATE comments SET content = $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query, comment.ID, comment.Content).Scan(&comment.UpdatedAt)
}

func (r *CommentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE comments SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
