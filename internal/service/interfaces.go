package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/kanban-saas/pkg/model"
)

type CardRepository interface {
	Create(ctx context.Context, card *model.Card) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Card, error)
	ListByBoard(ctx context.Context, boardID uuid.UUID) ([]model.Card, error)
	ListByColumn(ctx context.Context, columnID uuid.UUID) ([]model.Card, error)
	Update(ctx context.Context, card *model.Card) error
	Move(ctx context.Context, id uuid.UUID, columnID uuid.UUID, position int) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetMaxPosition(ctx context.Context, columnID uuid.UUID) (int, error)
	AddAssignee(ctx context.Context, cardID, userID uuid.UUID) error
	RemoveAssignee(ctx context.Context, cardID, userID uuid.UUID) error
	GetAssignees(ctx context.Context, cardID uuid.UUID) ([]uuid.UUID, error)
}

type LabelRepository interface {
	Create(ctx context.Context, label *model.Label) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Label, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]model.Label, error)
	Update(ctx context.Context, label *model.Label) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddToCard(ctx context.Context, cardID, labelID uuid.UUID) error
	RemoveFromCard(ctx context.Context, cardID, labelID uuid.UUID) error
	GetByCard(ctx context.Context, cardID uuid.UUID) ([]model.Label, error)
}

type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Comment, error)
	ListByCard(ctx context.Context, cardID uuid.UUID) ([]model.Comment, error)
	Update(ctx context.Context, comment *model.Comment) error
	Delete(ctx context.Context, id uuid.UUID) error
}
