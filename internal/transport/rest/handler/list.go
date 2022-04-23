package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func listCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		listID := chi.URLParam(r, "listID")

		ctx := context.WithValue(r.Context(), "list", listID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getAllLists(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getAllLists"))
}

func createList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("createList"))
}

func getList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getList"))
}

func updateList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("updateList"))
}

func deleteList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("deleteList"))
}
