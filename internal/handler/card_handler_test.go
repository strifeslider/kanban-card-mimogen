package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/user/kanban-saas/pkg/auth"
)

func TestNewCardHandler(t *testing.T) {
	h := &CardHandler{}
	if h == nil {
		t.Error("expected non-nil handler")
	}
}

func TestNewLabelHandler(t *testing.T) {
	h := &LabelHandler{}
	if h == nil {
		t.Error("expected non-nil handler")
	}
}

func TestNewCommentHandler(t *testing.T) {
	h := &CommentHandler{}
	if h == nil {
		t.Error("expected non-nil handler")
	}
}

func TestSetupRoutes(t *testing.T) {
	r := chi.NewRouter()
	ch := &CardHandler{}
	lh := &LabelHandler{}
	cmh := &CommentHandler{}
	jwtCfg := auth.JWTConfig{Secret: "test"}

	SetupRoutes(r, ch, lh, cmh, jwtCfg)

	// Verify routes exist
	routes := []string{
		"/api/v1/cards",
		"/api/v1/boards/test-id/cards",
		"/api/v1/workspaces/test-id/labels",
		"/api/v1/cards/test-id/comments",
	}

	for _, route := range routes {
		req := httptest.NewRequest("GET", route, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code == http.StatusNotFound {
			t.Errorf("route %s not found", route)
		}
	}
}

func TestCardHandler_Create_EmptyBody(t *testing.T) {
	h := &CardHandler{}
	req := httptest.NewRequest("POST", "/api/v1/cards", nil)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCardHandler_Get_InvalidID(t *testing.T) {
	h := &CardHandler{}
	req := httptest.NewRequest("GET", "/api/v1/cards/invalid", nil)
	w := httptest.NewRecorder()

	h.Get(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestLabelHandler_Create_EmptyBody(t *testing.T) {
	h := &LabelHandler{}
	req := httptest.NewRequest("POST", "/api/v1/workspaces/test/labels", nil)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCommentHandler_Create_EmptyBody(t *testing.T) {
	h := &CommentHandler{}
	req := httptest.NewRequest("POST", "/api/v1/cards/test/comments", nil)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
