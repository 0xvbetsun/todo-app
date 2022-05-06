package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/vbetsun/todo-app/internal/core"
)

type AuthService interface {
	CreateUser(user core.User) (core.User, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type AuthHandler struct {
	service AuthService
}

type SignUpRequest struct {
	*core.User
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service}
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

func (si *SignInRequest) Bind(r *http.Request) error {
	if si.Username == "" {
		return errors.New("missing required Username field")
	}
	if si.Password == "" {
		return errors.New("missing required Password field")
	}
	return nil
}

func (rd *SignUpResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (rd *SignInResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	data := &SignUpRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	u, err := h.service.CreateUser(*data.User)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.Render(w, r, &SignUpResponse{ID: u.ID, Name: u.Name, Username: u.Username})
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
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
