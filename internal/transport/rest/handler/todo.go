package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) todoCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		todoID := chi.URLParam(r, "todoID")

		ctx := context.WithValue(r.Context(), "todo", todoID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) getAllTodos(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getAllTodos"))
}

func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("createTodo"))
}

func (h *Handler) getTodo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getTodo"))
}

func (h *Handler) updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("updateTodo"))
}

func (h *Handler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("deleteTodo"))
}
