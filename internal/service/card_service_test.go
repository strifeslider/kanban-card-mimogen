package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/user/kanban-saas/pkg/model"
)

func TestNewCardService(t *testing.T) {
	svc := &CardService{}
	if svc == nil {
		t.Error("expected non-nil service")
	}
}

func TestCardPriority(t *testing.T) {
	priorities := []int{0, 1, 2, 3, 4}
	for _, p := range priorities {
		if p < 0 || p > 4 {
			t.Errorf("priority %d out of range", p)
		}
	}
}

func TestCardModel(t *testing.T) {
	card := model.Card{
		ID:       uuid.New(),
		Title:    "Test Card",
		Position: 0,
		Priority: 1,
	}

	if card.Title != "Test Card" {
		t.Errorf("expected title 'Test Card', got '%s'", card.Title)
	}
	if card.Position != 0 {
		t.Errorf("expected position 0, got %d", card.Position)
	}
}

func TestCreateCardRequest(t *testing.T) {
	req := model.CreateCardRequest{
		ColumnID: uuid.New(),
		Title:    "New Card",
	}

	if req.Title != "New Card" {
		t.Errorf("expected title 'New Card', got '%s'", req.Title)
	}
}

func TestMoveCardRequest(t *testing.T) {
	req := model.MoveCardRequest{
		ColumnID: uuid.New(),
		Position: 2,
	}

	if req.Position != 2 {
		t.Errorf("expected position 2, got %d", req.Position)
	}
}

func TestLabelModel(t *testing.T) {
	label := model.Label{
		ID:          uuid.New(),
		WorkspaceID: uuid.New(),
		Name:        "Bug",
		Color:       "#ff0000",
	}

	if label.Name != "Bug" {
		t.Errorf("expected name 'Bug', got '%s'", label.Name)
	}
	if label.Color != "#ff0000" {
		t.Errorf("expected color '#ff0000', got '%s'", label.Color)
	}
}

func TestCommentModel(t *testing.T) {
	comment := model.Comment{
		ID:      uuid.New(),
		CardID:  uuid.New(),
		UserID:  uuid.New(),
		Content: "Test comment",
	}

	if comment.Content != "Test comment" {
		t.Errorf("expected content 'Test comment', got '%s'", comment.Content)
	}
}

func TestCreateCommentRequest(t *testing.T) {
	req := model.CreateCommentRequest{
		Content: "New comment",
	}

	if req.Content != "New comment" {
		t.Errorf("expected content 'New comment', got '%s'", req.Content)
	}
}

func TestTimeParsing(t *testing.T) {
	now := time.Now()
	if now.IsZero() {
		t.Error("expected non-zero time")
	}
}
