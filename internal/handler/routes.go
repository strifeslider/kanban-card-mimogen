package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/user/kanban-saas/pkg/auth"
)

func SetupRoutes(
	r chi.Router,
	cardHandler *CardHandler,
	labelHandler *LabelHandler,
	commentHandler *CommentHandler,
	jwtCfg auth.JWTConfig,
) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(auth.RequireAuth(jwtCfg))

		r.Route("/cards", func(r chi.Router) {
			r.Post("/", cardHandler.Create)
			r.Get("/{id}", cardHandler.Get)
			r.Put("/{id}", cardHandler.Update)
			r.Delete("/{id}", cardHandler.Delete)
			r.Post("/{id}/move", cardHandler.Move)

			r.Post("/{id}/assignees", cardHandler.AddAssignee)
			r.Delete("/{id}/assignees/{userId}", cardHandler.RemoveAssignee)

			r.Post("/{id}/labels/{labelId}", cardHandler.AddLabel)
			r.Delete("/{id}/labels/{labelId}", cardHandler.RemoveLabel)
		})

		r.Get("/boards/{boardId}/cards", cardHandler.ListByBoard)

		r.Route("/workspaces/{workspaceId}/labels", func(r chi.Router) {
			r.Post("/", labelHandler.Create)
			r.Get("/", labelHandler.List)
		})

		r.Route("/labels", func(r chi.Router) {
			r.Put("/{id}", labelHandler.Update)
			r.Delete("/{id}", labelHandler.Delete)
		})

		r.Route("/cards/{cardId}/comments", func(r chi.Router) {
			r.Post("/", commentHandler.Create)
			r.Get("/", commentHandler.List)
		})

		r.Route("/comments", func(r chi.Router) {
			r.Put("/{id}", commentHandler.Update)
			r.Delete("/{id}", commentHandler.Delete)
		})
	})
}
