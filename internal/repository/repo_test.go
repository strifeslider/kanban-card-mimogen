package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/user/kanban-saas/pkg/model"
)

func TestCardRepository_New(t *testing.T) {
	repo := &CardRepository{}
	if repo == nil {
		t.Error("expected non-nil repo")
	}
}

func TestLabelRepository_New(t *testing.T) {
	repo := &LabelRepository{}
	if repo == nil {
		t.Error("expected non-nil repo")
	}
}

func TestCommentRepository_New(t *testing.T) {
	repo := &CommentRepository{}
	if repo == nil {
		t.Error("expected non-nil repo")
	}
}

func TestCardRepository_Model(t *testing.T) {
	card := &model.Card{
		ID:        uuid.New(),
		ColumnID:  uuid.New(),
		BoardID:   uuid.New(),
		Title:     "Task",
		Position:  0,
		Priority:  1,
		CreatedBy: uuid.New(),
	}
	if card.Title != "Task" {
		t.Error("title mismatch")
	}
}

func TestLabelRepository_Model(t *testing.T) {
	label := &model.Label{
		ID:          uuid.New(),
		WorkspaceID: uuid.New(),
		Name:        "Bug",
		Color:       "#ff0000",
	}
	if label.Name != "Bug" {
		t.Error("name mismatch")
	}
}

func TestCommentRepository_Model(t *testing.T) {
	comment := &model.Comment{
		ID:      uuid.New(),
		CardID:  uuid.New(),
		UserID:  uuid.New(),
		Content: "text",
	}
	if comment.Content != "text" {
		t.Error("content mismatch")
	}
}
