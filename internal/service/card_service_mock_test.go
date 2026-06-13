package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/user/kanban-saas/pkg/mock"
	"github.com/user/kanban-saas/pkg/model"
)

func newTestCardService() (*CardService, *mock.MockCardRepo, *mock.MockLabelRepo, *mock.MockCommentRepo) {
	cardRepo := mock.NewMockCardRepo()
	labelRepo := mock.NewMockLabelRepo()
	commentRepo := mock.NewMockCommentRepo()
	svc := NewCardService(cardRepo, labelRepo, commentRepo)
	return svc, cardRepo, labelRepo, commentRepo
}

func TestCardService_Create(t *testing.T) {
	svc, _, _, _ := newTestCardService()
	ctx := context.Background()

	card, err := svc.Create(ctx, uuid.New(), model.CreateCardRequest{
		ColumnID: uuid.New(),
		Title:    "New Card",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if card.Title != "New Card" {
		t.Errorf("expected title 'New Card', got '%s'", card.Title)
	}
}

func TestCardService_Create_WithPriority(t *testing.T) {
	svc, _, _, _ := newTestCardService()
	ctx := context.Background()

	priority := 3
	card, err := svc.Create(ctx, uuid.New(), model.CreateCardRequest{
		ColumnID: uuid.New(),
		Title:    "High Priority",
		Priority: &priority,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if card.Priority != 3 {
		t.Errorf("expected priority 3, got %d", card.Priority)
	}
}

func TestCardService_GetByID(t *testing.T) {
	svc, cardRepo, _, _ := newTestCardService()
	ctx := context.Background()

	cardID := uuid.New()
	cardRepo.Cards[cardID] = &model.Card{
		ID:    cardID,
		Title: "Test",
	}

	card, err := svc.GetByID(ctx, cardID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if card.Title != "Test" {
		t.Errorf("expected title 'Test', got '%s'", card.Title)
	}
}

func TestCardService_ListByBoard(t *testing.T) {
	svc, cardRepo, _, _ := newTestCardService()
	ctx := context.Background()

	boardID := uuid.New()
	cardRepo.Cards[uuid.New()] = &model.Card{BoardID: boardID, Title: "Card1"}
	cardRepo.Cards[uuid.New()] = &model.Card{BoardID: boardID, Title: "Card2"}

	cards, err := svc.ListByBoard(ctx, boardID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cards) != 2 {
		t.Errorf("expected 2 cards, got %d", len(cards))
	}
}

func TestCardService_ListByColumn(t *testing.T) {
	svc, cardRepo, _, _ := newTestCardService()
	ctx := context.Background()

	columnID := uuid.New()
	cardRepo.Cards[uuid.New()] = &model.Card{ColumnID: columnID, Title: "Card1"}

	cards, err := svc.ListByColumn(ctx, columnID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cards) != 1 {
		t.Errorf("expected 1 card, got %d", len(cards))
	}
}

func TestCardService_Update(t *testing.T) {
	svc, cardRepo, _, _ := newTestCardService()
	ctx := context.Background()

	cardID := uuid.New()
	cardRepo.Cards[cardID] = &model.Card{
		ID:    cardID,
		Title: "Old Title",
	}

	newTitle := "New Title"
	card, err := svc.Update(ctx, cardID, model.UpdateCardRequest{Title: &newTitle})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if card.Title != "New Title" {
		t.Errorf("expected title 'New Title', got '%s'", card.Title)
	}
}

func TestCardService_Move(t *testing.T) {
	svc, cardRepo, _, _ := newTestCardService()
	ctx := context.Background()

	cardID := uuid.New()
	cardRepo.Cards[cardID] = &model.Card{
		ID:       cardID,
		ColumnID: uuid.New(),
		Position: 0,
	}

	newColumnID := uuid.New()
	err := svc.Move(ctx, cardID, model.MoveCardRequest{
		ColumnID: newColumnID,
		Position: 5,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCardService_Delete(t *testing.T) {
	svc, cardRepo, _, _ := newTestCardService()
	ctx := context.Background()

	cardID := uuid.New()
	cardRepo.Cards[cardID] = &model.Card{ID: cardID}

	err := svc.Delete(ctx, cardID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := cardRepo.Cards[cardID]; ok {
		t.Error("card should be deleted")
	}
}

func TestCardService_AddAssignee(t *testing.T) {
	svc, _, _, _ := newTestCardService()
	ctx := context.Background()

	err := svc.AddAssignee(ctx, uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCardService_RemoveAssignee(t *testing.T) {
	svc, _, _, _ := newTestCardService()
	ctx := context.Background()

	err := svc.RemoveAssignee(ctx, uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCardService_GetAssignees(t *testing.T) {
	svc, _, _, _ := newTestCardService()
	ctx := context.Background()

	assignees, err := svc.GetAssignees(ctx, uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if assignees != nil {
		t.Error("expected nil assignees")
	}
}

func TestCardService_AddLabel(t *testing.T) {
	svc, _, _, _ := newTestCardService()
	ctx := context.Background()

	err := svc.AddLabel(ctx, uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCardService_RemoveLabel(t *testing.T) {
	svc, _, _, _ := newTestCardService()
	ctx := context.Background()

	err := svc.RemoveLabel(ctx, uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCardService_GetLabels(t *testing.T) {
	svc, _, _, _ := newTestCardService()
	ctx := context.Background()

	labels, err := svc.GetLabels(ctx, uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if labels != nil {
		t.Error("expected nil labels")
	}
}

func TestCardService_CreateComment(t *testing.T) {
	svc, _, _, _ := newTestCardService()
	ctx := context.Background()

	comment, err := svc.CreateComment(ctx, uuid.New(), uuid.New(), "Hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.Content != "Hello" {
		t.Errorf("expected content 'Hello', got '%s'", comment.Content)
	}
}

func TestCardService_ListComments(t *testing.T) {
	svc, _, _, commentRepo := newTestCardService()
	ctx := context.Background()

	cardID := uuid.New()
	commentRepo.Comments[uuid.New()] = &model.Comment{CardID: cardID, Content: "msg1"}

	comments, err := svc.ListComments(ctx, cardID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(comments))
	}
}

func TestCardService_UpdateComment(t *testing.T) {
	svc, _, _, commentRepo := newTestCardService()
	ctx := context.Background()

	commentID := uuid.New()
	commentRepo.Comments[commentID] = &model.Comment{
		ID:      commentID,
		Content: "old",
	}

	comment, err := svc.UpdateComment(ctx, commentID, "new")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.Content != "new" {
		t.Errorf("expected content 'new', got '%s'", comment.Content)
	}
}

func TestCardService_DeleteComment(t *testing.T) {
	svc, _, _, commentRepo := newTestCardService()
	ctx := context.Background()

	commentID := uuid.New()
	commentRepo.Comments[commentID] = &model.Comment{ID: commentID}

	err := svc.DeleteComment(ctx, commentID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCardService_CreateLabel(t *testing.T) {
	svc, _, _, _ := newTestCardService()
	ctx := context.Background()

	label, err := svc.CreateLabel(ctx, uuid.New(), model.CreateLabelRequest{
		Name:  "Bug",
		Color: "#ff0000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if label.Name != "Bug" {
		t.Errorf("expected name 'Bug', got '%s'", label.Name)
	}
}

func TestCardService_ListLabels(t *testing.T) {
	svc, _, labelRepo, _ := newTestCardService()
	ctx := context.Background()

	workspaceID := uuid.New()
	labelRepo.Labels[uuid.New()] = &model.Label{WorkspaceID: workspaceID, Name: "Label1"}

	labels, err := svc.ListLabels(ctx, workspaceID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(labels) != 1 {
		t.Errorf("expected 1 label, got %d", len(labels))
	}
}

func TestCardService_UpdateLabel(t *testing.T) {
	svc, _, labelRepo, _ := newTestCardService()
	ctx := context.Background()

	labelID := uuid.New()
	labelRepo.Labels[labelID] = &model.Label{
		ID:    labelID,
		Name:  "Old",
		Color: "#000000",
	}

	newName := "Updated"
	newColor := "#ffffff"
	label, err := svc.UpdateLabel(ctx, labelID, model.UpdateLabelRequest{
		Name:  &newName,
		Color: &newColor,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if label.Name != "Updated" {
		t.Errorf("expected name 'Updated', got '%s'", label.Name)
	}
}

func TestCardService_DeleteLabel(t *testing.T) {
	svc, _, _, _ := newTestCardService()
	ctx := context.Background()

	err := svc.DeleteLabel(ctx, uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewCardService(t *testing.T) {
	svc, _, _, _ := newTestCardService()
	if svc == nil {
		t.Error("expected non-nil service")
	}
}
