package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

// Key to use when setting the user context.
type ctxKeyUser string

const (
	authHeader            = "Authorization"
	userCtx    ctxKeyUser = "userID"
)

func (h *AuthHandler) UserIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader, ok := r.Header[authHeader]
		if !ok {
			if err := render.Render(w, r, ErrUnauthorized(errors.New("empty auth header"))); err != nil {
				h.log.Error(ErrRenderResp.Error())
			}
			return
		}
		headerParts := strings.Fields(authHeader[0])
		if len(headerParts) != 2 {
			if err := render.Render(w, r, ErrUnauthorized(errors.New("invalid auth header"))); err != nil {
				h.log.Error(ErrRenderResp.Error())
			}
			return
		}
		userID, err := h.service.ParseToken(headerParts[1])
		if err != nil {
			if rErr := render.Render(w, r, ErrUnauthorized(err)); rErr != nil {
				h.log.Error(ErrRenderResp.Error())
			}
			return
		}
		ctx := context.WithValue(r.Context(), userCtx, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
func (h *Handler) Logger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				h.log.Info("Served",
					zap.String("proto", r.Proto),
					zap.String("path", r.URL.Path),
					zap.Duration("lat", time.Since(t1)),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
					zap.String("reqId", middleware.GetReqID(r.Context())))
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

// Recoverer is a middleware for recovering after panics during handling http requests
func (h *Handler) Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}
				if err := render.Render(w, r, ErrInternalServer(errors.New("something went wrong"))); err != nil {
					h.log.Error(ErrRenderResp.Error())
				}
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// SendRequestID is a middleware for sending X-Request-Id to the client
func SendRequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(
			middleware.RequestIDHeader,
			middleware.GetReqID(r.Context()),
		)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func getUserID(w http.ResponseWriter, r *http.Request) (int, error) {
	userID, ok := r.Context().Value(userCtx).(int)
	if !ok {
		return 0, errors.New("userID not found")
	}
	return userID, nil
}
