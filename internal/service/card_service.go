package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/user/kanban-saas/pkg/model"
)

type CardService struct {
	cardRepo    CardRepository
	labelRepo   LabelRepository
	commentRepo CommentRepository
}

func NewCardService(
	cardRepo CardRepository,
	labelRepo LabelRepository,
	commentRepo CommentRepository,
) *CardService {
	return &CardService{
		cardRepo:    cardRepo,
		labelRepo:   labelRepo,
		commentRepo: commentRepo,
	}
}

func (s *CardService) Create(ctx context.Context, userID uuid.UUID, req model.CreateCardRequest) (*model.Card, error) {
	maxPos, err := s.cardRepo.GetMaxPosition(ctx, req.ColumnID)
	if err != nil {
		return nil, fmt.Errorf("get max position: %w", err)
	}

	pos := maxPos + 1

	priority := 0
	if req.Priority != nil {
		priority = *req.Priority
	}

	card := &model.Card{
		ID:          uuid.New(),
		ColumnID:    req.ColumnID,
		BoardID:     uuid.Nil,
		Title:       req.Title,
		Description: req.Description,
		Position:    pos,
		Priority:    priority,
		DueDate:     req.DueDate,
		CreatedBy:   userID,
	}

	if err := s.cardRepo.Create(ctx, card); err != nil {
		return nil, fmt.Errorf("create card: %w", err)
	}

	return card, nil
}

func (s *CardService) GetByID(ctx context.Context, id uuid.UUID) (*model.Card, error) {
	return s.cardRepo.GetByID(ctx, id)
}

func (s *CardService) ListByBoard(ctx context.Context, boardID uuid.UUID) ([]model.Card, error) {
	return s.cardRepo.ListByBoard(ctx, boardID)
}

func (s *CardService) ListByColumn(ctx context.Context, columnID uuid.UUID) ([]model.Card, error) {
	return s.cardRepo.ListByColumn(ctx, columnID)
}

func (s *CardService) Update(ctx context.Context, id uuid.UUID, req model.UpdateCardRequest) (*model.Card, error) {
	card, err := s.cardRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		card.Title = *req.Title
	}
	if req.Description != nil {
		card.Description = req.Description
	}
	if req.Priority != nil {
		card.Priority = *req.Priority
	}
	if req.DueDate != nil {
		card.DueDate = req.DueDate
	}

	if err := s.cardRepo.Update(ctx, card); err != nil {
		return nil, fmt.Errorf("update card: %w", err)
	}

	return card, nil
}

func (s *CardService) Move(ctx context.Context, id uuid.UUID, req model.MoveCardRequest) error {
	return s.cardRepo.Move(ctx, id, req.ColumnID, req.Position)
}

func (s *CardService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.cardRepo.Delete(ctx, id)
}

func (s *CardService) AddAssignee(ctx context.Context, cardID, userID uuid.UUID) error {
	return s.cardRepo.AddAssignee(ctx, cardID, userID)
}

func (s *CardService) RemoveAssignee(ctx context.Context, cardID, userID uuid.UUID) error {
	return s.cardRepo.RemoveAssignee(ctx, cardID, userID)
}

func (s *CardService) GetAssignees(ctx context.Context, cardID uuid.UUID) ([]uuid.UUID, error) {
	return s.cardRepo.GetAssignees(ctx, cardID)
}

func (s *CardService) AddLabel(ctx context.Context, cardID, labelID uuid.UUID) error {
	return s.labelRepo.AddToCard(ctx, cardID, labelID)
}

func (s *CardService) RemoveLabel(ctx context.Context, cardID, labelID uuid.UUID) error {
	return s.labelRepo.RemoveFromCard(ctx, cardID, labelID)
}

func (s *CardService) GetLabels(ctx context.Context, cardID uuid.UUID) ([]model.Label, error) {
	return s.labelRepo.GetByCard(ctx, cardID)
}

func (s *CardService) CreateComment(ctx context.Context, userID, cardID uuid.UUID, content string) (*model.Comment, error) {
	comment := &model.Comment{
		ID:      uuid.New(),
		CardID:  cardID,
		UserID:  userID,
		Content: content,
	}

	if err := s.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}

	return comment, nil
}

func (s *CardService) ListComments(ctx context.Context, cardID uuid.UUID) ([]model.Comment, error) {
	return s.commentRepo.ListByCard(ctx, cardID)
}

func (s *CardService) UpdateComment(ctx context.Context, id uuid.UUID, content string) (*model.Comment, error) {
	comment, err := s.commentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	comment.Content = content
	if err := s.commentRepo.Update(ctx, comment); err != nil {
		return nil, fmt.Errorf("update comment: %w", err)
	}

	return comment, nil
}

func (s *CardService) DeleteComment(ctx context.Context, id uuid.UUID) error {
	return s.commentRepo.Delete(ctx, id)
}

func (s *CardService) CreateLabel(ctx context.Context, workspaceID uuid.UUID, req model.CreateLabelRequest) (*model.Label, error) {
	label := &model.Label{
		ID:          uuid.New(),
		WorkspaceID: workspaceID,
		Name:        req.Name,
		Color:       req.Color,
	}

	if err := s.labelRepo.Create(ctx, label); err != nil {
		return nil, fmt.Errorf("create label: %w", err)
	}

	return label, nil
}

func (s *CardService) ListLabels(ctx context.Context, workspaceID uuid.UUID) ([]model.Label, error) {
	return s.labelRepo.ListByWorkspace(ctx, workspaceID)
}

func (s *CardService) UpdateLabel(ctx context.Context, id uuid.UUID, req model.UpdateLabelRequest) (*model.Label, error) {
	label, err := s.labelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		label.Name = *req.Name
	}
	if req.Color != nil {
		label.Color = *req.Color
	}

	if err := s.labelRepo.Update(ctx, label); err != nil {
		return nil, fmt.Errorf("update label: %w", err)
	}

	return label, nil
}

func (s *CardService) DeleteLabel(ctx context.Context, id uuid.UUID) error {
	return s.labelRepo.Delete(ctx, id)
}
