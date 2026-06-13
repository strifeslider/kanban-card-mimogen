package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/user/kanban-saas/pkg/model"
)

type CardRepository struct {
	db *pgxpool.Pool
}

func NewCardRepository(db *pgxpool.Pool) *CardRepository {
	return &CardRepository{db: db}
}

func (r *CardRepository) Create(ctx context.Context, card *model.Card) error {
	query := `
		INSERT INTO cards (id, column_id, board_id, title, description, position, priority, due_date, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		card.ID, card.ColumnID, card.BoardID, card.Title, card.Description,
		card.Position, card.Priority, card.DueDate, card.CreatedBy,
	).Scan(&card.CreatedAt, &card.UpdatedAt)
}

func (r *CardRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Card, error) {
	query := `
		SELECT id, column_id, board_id, title, description, position, priority, due_date, created_by, created_at, updated_at
		FROM cards WHERE id = $1 AND deleted_at IS NULL`

	card := &model.Card{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&card.ID, &card.ColumnID, &card.BoardID, &card.Title, &card.Description,
		&card.Position, &card.Priority, &card.DueDate, &card.CreatedBy,
		&card.CreatedAt, &card.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get card: %w", err)
	}
	return card, nil
}

func (r *CardRepository) ListByBoard(ctx context.Context, boardID uuid.UUID) ([]model.Card, error) {
	query := `
		SELECT id, column_id, board_id, title, description, position, priority, due_date, created_by, created_at, updated_at
		FROM cards WHERE board_id = $1 AND deleted_at IS NULL
		ORDER BY position ASC`

	rows, err := r.db.Query(ctx, query, boardID)
	if err != nil {
		return nil, fmt.Errorf("list cards: %w", err)
	}
	defer rows.Close()

	var cards []model.Card
	for rows.Next() {
		var card model.Card
		if err := rows.Scan(
			&card.ID, &card.ColumnID, &card.BoardID, &card.Title, &card.Description,
			&card.Position, &card.Priority, &card.DueDate, &card.CreatedBy,
			&card.CreatedAt, &card.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan card: %w", err)
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func (r *CardRepository) ListByColumn(ctx context.Context, columnID uuid.UUID) ([]model.Card, error) {
	query := `
		SELECT id, column_id, board_id, title, description, position, priority, due_date, created_by, created_at, updated_at
		FROM cards WHERE column_id = $1 AND deleted_at IS NULL
		ORDER BY position ASC`

	rows, err := r.db.Query(ctx, query, columnID)
	if err != nil {
		return nil, fmt.Errorf("list cards by column: %w", err)
	}
	defer rows.Close()

	var cards []model.Card
	for rows.Next() {
		var card model.Card
		if err := rows.Scan(
			&card.ID, &card.ColumnID, &card.BoardID, &card.Title, &card.Description,
			&card.Position, &card.Priority, &card.DueDate, &card.CreatedBy,
			&card.CreatedAt, &card.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan card: %w", err)
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func (r *CardRepository) Update(ctx context.Context, card *model.Card) error {
	query := `
		UPDATE cards SET title = $2, description = $3, priority = $4, due_date = $5, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`

	return r.db.QueryRow(ctx, query, card.ID, card.Title, card.Description, card.Priority, card.DueDate).Scan(&card.UpdatedAt)
}

func (r *CardRepository) Move(ctx context.Context, id uuid.UUID, columnID uuid.UUID, position int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`UPDATE cards SET column_id = $2, position = $3, updated_at = NOW() WHERE id = $1`,
		id, columnID, position,
	)
	if err != nil {
		return fmt.Errorf("move card: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *CardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE cards SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *CardRepository) GetMaxPosition(ctx context.Context, columnID uuid.UUID) (int, error) {
	var maxPos int
	query := `SELECT COALESCE(MAX(position), -1) FROM cards WHERE column_id = $1 AND deleted_at IS NULL`
	err := r.db.QueryRow(ctx, query, columnID).Scan(&maxPos)
	return maxPos, err
}

func (r *CardRepository) AddAssignee(ctx context.Context, cardID, userID uuid.UUID) error {
	query := `INSERT INTO card_assignees (id, card_id, user_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, uuid.New(), cardID, userID)
	return err
}

func (r *CardRepository) RemoveAssignee(ctx context.Context, cardID, userID uuid.UUID) error {
	query := `DELETE FROM card_assignees WHERE card_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, cardID, userID)
	return err
}

func (r *CardRepository) GetAssignees(ctx context.Context, cardID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT user_id FROM card_assignees WHERE card_id = $1`
	rows, err := r.db.Query(ctx, query, cardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, nil
}
