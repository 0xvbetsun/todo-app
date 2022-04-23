package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func todoCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		todoID := chi.URLParam(r, "todoID")

		ctx := context.WithValue(r.Context(), "todo", todoID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAllTodos(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getAllTodos"))
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("createTodo"))
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getTodo"))
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("updateTodo"))
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("deleteTodo"))
}
