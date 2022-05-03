package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

const (
	authHeader = "Authorization"
	userCtx    = "userID"
)

func (h *Handler) UserIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader, ok := r.Header[authHeader]
		if !ok {
			render.Render(w, r, ErrUnauthorized(errors.New("empty auth header")))
			return
		}
		headerParts := strings.Fields(authHeader[0])
		if len(headerParts) != 2 {
			render.Render(w, r, ErrUnauthorized(errors.New("invalid auth header")))
			return
		}
		userID, err := h.service.ParseToken(headerParts[1])
		if err != nil {
			render.Render(w, r, ErrUnauthorized(err))
			return
		}
		ctx := context.WithValue(r.Context(), userCtx, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserID(w http.ResponseWriter, r *http.Request) (int, error) {
	userID, ok := r.Context().Value(userCtx).(int)
	if !ok {
		return 0, errors.New("userID not found")
	}
	return userID, nil
}
