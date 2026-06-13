package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/user/kanban-saas/pkg/auth"
	apperr "github.com/user/kanban-saas/pkg/errors"
	"github.com/user/kanban-saas/services/card/internal/service"
)

type CommentHandler struct {
	cardService *service.CardService
}

func NewCommentHandler(cardService *service.CardService) *CommentHandler {
	return &CommentHandler{cardService: cardService}
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	cardID, err := uuid.Parse(chi.URLParam(r, "cardId"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid card id"))
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	if req.Content == "" {
		apperr.RespondError(w, apperr.BadRequest("content is required"))
		return
	}

	comment, err := h.cardService.CreateComment(r.Context(), userID, cardID, req.Content)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusCreated, comment)
}

func (h *CommentHandler) List(w http.ResponseWriter, r *http.Request) {
	cardID, err := uuid.Parse(chi.URLParam(r, "cardId"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid card id"))
		return
	}

	comments, err := h.cardService.ListComments(r.Context(), cardID)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, comments)
}

func (h *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid comment id"))
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	if req.Content == "" {
		apperr.RespondError(w, apperr.BadRequest("content is required"))
		return
	}

	comment, err := h.cardService.UpdateComment(r.Context(), id, req.Content)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, comment)
}

func (h *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid comment id"))
		return
	}

	if err := h.cardService.DeleteComment(r.Context(), id); err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, map[string]string{"message": "comment deleted"})
}
