package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	apperr "github.com/user/kanban-saas/pkg/errors"
	"github.com/user/kanban-saas/pkg/model"
	"github.com/user/kanban-saas/services/card/internal/service"
)

type LabelHandler struct {
	cardService *service.CardService
}

func NewLabelHandler(cardService *service.CardService) *LabelHandler {
	return &LabelHandler{cardService: cardService}
}

func (h *LabelHandler) Create(w http.ResponseWriter, r *http.Request) {
	workspaceID, err := uuid.Parse(chi.URLParam(r, "workspaceId"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid workspace id"))
		return
	}

	var req model.CreateLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	if req.Name == "" || req.Color == "" {
		apperr.RespondError(w, apperr.BadRequest("name and color are required"))
		return
	}

	label, err := h.cardService.CreateLabel(r.Context(), workspaceID, req)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusCreated, label)
}

func (h *LabelHandler) List(w http.ResponseWriter, r *http.Request) {
	workspaceID, err := uuid.Parse(chi.URLParam(r, "workspaceId"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid workspace id"))
		return
	}

	labels, err := h.cardService.ListLabels(r.Context(), workspaceID)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, labels)
}

func (h *LabelHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid label id"))
		return
	}

	var req model.UpdateLabelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid request body"))
		return
	}

	label, err := h.cardService.UpdateLabel(r.Context(), id, req)
	if err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, label)
}

func (h *LabelHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		apperr.RespondError(w, apperr.BadRequest("invalid label id"))
		return
	}

	if err := h.cardService.DeleteLabel(r.Context(), id); err != nil {
		apperr.RespondError(w, err)
		return
	}

	apperr.RespondJSON(w, http.StatusOK, map[string]string{"message": "label deleted"})
}
