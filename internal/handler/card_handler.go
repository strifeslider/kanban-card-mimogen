package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/user/kanban-saas/pkg/auth"
	apperr "github.com/user/kanban-saas/pkg/errors"
	"github.com/user/kanban-saas/pkg/model"
	"github.com/user/kanban-saas/services/card/internal/service"
)

type CardHandler struct {
	cardService *service.CardService
}

func NewCardHandler(cardService *service.CardService) *CardHandler {
	return &CardHandler{cardService: cardService}
}

func (h *CardHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)

	var req model.CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	if req.Title == "" {
		apperr.RespondError(w, apperr.BadRequest("title is required"))
		return
	}

	card, err := h.cardService.Create(r.Context(), userID, req)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusCreated, card)
}

func (h *CardHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid card id"))
		return
	}

	card, err := h.cardService.GetByID(r.Context(), id)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, card)
}

func (h *CardHandler) ListByBoard(w http.ResponseWriter, r *http.Request) {
	boardID, err := uuid.Parse(chi.URLParam(r, "boardId"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid board id"))
		return
	}

	cards, err := h.cardService.ListByBoard(r.Context(), boardID)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, cards)
}

func (h *CardHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid card id"))
		return
	}

	var req model.UpdateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	card, err := h.cardService.Update(r.Context(), id, req)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, card)
}

func (h *CardHandler) Move(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid card id"))
		return
	}

	var req model.MoveCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	if err := h.cardService.Move(r.Context(), id, req); err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, map[string]string{"message": "card moved"})
}

func (h *CardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid card id"))
		return
	}

	if err := h.cardService.Delete(r.Context(), id); err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, map[string]string{"message": "card deleted"})
}

func (h *CardHandler) AddAssignee(w http.ResponseWriter, r *http.Request) {
	cardID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid card id"))
		return
	}

	var req struct {
		UserID uuid.UUID `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	if err := h.cardService.AddAssignee(r.Context(), cardID, req.UserID); err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusCreated, map[string]string{"message": "assignee added"})
}

func (h *CardHandler) RemoveAssignee(w http.ResponseWriter, r *http.Request) {
	cardID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid card id"))
		return
	}

	userID, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid user id"))
		return
	}

	if err := h.cardService.RemoveAssignee(r.Context(), cardID, userID); err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, map[string]string{"message": "assignee removed"})
}

func (h *CardHandler) AddLabel(w http.ResponseWriter, r *http.Request) {
	cardID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid card id"))
		return
	}

	labelID, err := uuid.Parse(chi.URLParam(r, "labelId"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid label id"))
		return
	}

	if err := h.cardService.AddLabel(r.Context(), cardID, labelID); err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusCreated, map[string]string{"message": "label added"})
}

func (h *CardHandler) RemoveLabel(w http.ResponseWriter, r *http.Request) {
	cardID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid card id"))
		return
	}

	labelID, err := uuid.Parse(chi.URLParam(r, "labelId"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid label id"))
		return
	}

	if err := h.cardService.RemoveLabel(r.Context(), cardID, labelID); err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, map[string]string{"message": "label removed"})
}
