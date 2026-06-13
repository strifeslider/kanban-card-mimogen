package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/user/kanban-saas/pkg/model"
)

type LabelRepository struct {
	db *pgxpool.Pool
}

func NewLabelRepository(db *pgxpool.Pool) *LabelRepository {
	return &LabelRepository{db: db}
}

func (r *LabelRepository) Create(ctx context.Context, label *model.Label) error {
	query := `
		INSERT INTO labels (id, workspace_id, name, color)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at`

	return r.db.QueryRow(ctx, query,
		label.ID, label.WorkspaceID, label.Name, label.Color,
	).Scan(&label.CreatedAt)
}

func (r *LabelRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Label, error) {
	query := `
		SELECT id, workspace_id, name, color, created_at
		FROM labels WHERE id = $1`

	label := &model.Label{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&label.ID, &label.WorkspaceID, &label.Name, &label.Color, &label.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get label: %w", err)
	}
	return label, nil
}

func (r *LabelRepository) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.Label, error) {
	query := `
		SELECT id, workspace_id, name, color, created_at
		FROM labels WHERE workspace_id = $1
		ORDER BY name ASC`

	rows, err := r.db.Query(ctx, query, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("list labels: %w", err)
	}
	defer rows.Close()

	var labels []model.Label
	for rows.Next() {
		var label model.Label
		if err := rows.Scan(&label.ID, &label.WorkspaceID, &label.Name, &label.Color, &label.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan label: %w", err)
		}
		labels = append(labels, label)
	}
	return labels, nil
}

func (r *LabelRepository) Update(ctx context.Context, label *model.Label) error {
	query := `
		UPDATE labels SET name = $2, color = $3
		WHERE id = $1`

	_, err := r.db.Exec(ctx, query, label.ID, label.Name, label.Color)
	return err
}

func (r *LabelRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM labels WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *LabelRepository) AddToCard(ctx context.Context, cardID, labelID uuid.UUID) error {
	query := `INSERT INTO card_labels (card_id, label_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, cardID, labelID)
	return err
}

func (r *LabelRepository) RemoveFromCard(ctx context.Context, cardID, labelID uuid.UUID) error {
	query := `DELETE FROM card_labels WHERE card_id = $1 AND label_id = $2`
	_, err := r.db.Exec(ctx, query, cardID, labelID)
	return err
}

func (r *LabelRepository) GetByCard(ctx context.Context, cardID uuid.UUID) ([]model.Label, error) {
	query := `
		SELECT l.id, l.workspace_id, l.name, l.color, l.created_at
		FROM labels l
		INNER JOIN card_labels cl ON l.id = cl.label_id
		WHERE cl.card_id = $1
		ORDER BY l.name ASC`

	rows, err := r.db.Query(ctx, query, cardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var labels []model.Label
	for rows.Next() {
		var label model.Label
		if err := rows.Scan(&label.ID, &label.WorkspaceID, &label.Name, &label.Color, &label.CreatedAt); err != nil {
			return nil, err
		}
		labels = append(labels, label)
	}
	return labels, nil
}
