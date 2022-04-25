package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/vbetsun/todo-app/internal/domain"
)

type SignUpRequest struct {
	*domain.User
}

type SignUpResponse struct {
	ID int `json:"id"`
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

func (si *SignInRequest) Bind(r *http.Request) error {
	if si.Username == "" {
		return errors.New("missing required Username field")
	}
	if si.Password == "" {
		return errors.New("missing required Password field")
	}
	return nil
}

func (sr *SignUpRequest) Bind(r *http.Request) error {
	if sr.User == nil {
		return errors.New("missing required User fields")
	}
	if sr.Name == "" {
		return errors.New("missing required Name field")
	}
	if sr.Username == "" {
		return errors.New("missing required Username field")
	}
	if sr.Password == "" {
		return errors.New("missing required Password field")
	}

	return nil
}

func (rd *SignUpResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	data := &SignInRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	token, err := h.service.GenerateToken(data.Username, data.Password)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.Render(w, r, &SignInResponse{Token: token})
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	data := &SignUpRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	userID, err := h.service.CreateUser(*data.User)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.Render(w, r, &SignUpResponse{ID: userID})
}
