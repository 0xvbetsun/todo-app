package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vbetsun/todo-app/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/auth", h.authRouter())
	r.Mount("/api", h.apiRouter())
	return r
}

func (h *Handler) authRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/sign-in", signIn)
	r.Post("/sign-up", signUp)
	return r
}

func (h *Handler) apiRouter() http.Handler {
	r := chi.NewRouter()
	r.Mount("/lists", h.listsRouter())
	return r
}

func (h *Handler) listsRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getAllLists)
	r.Post("/", createList)
	r.Route("/{listID}", func(r chi.Router) {
		r.Use(listCtx)
		r.Get("/", getList)
		r.Put("/", updateList)
		r.Delete("/", deleteList)
		r.Mount("/todos", h.todosRouter())
	})
	return r
}

func (h *Handler) todosRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getAllTodos)
	r.Post("/", createTodo)
	r.Route("/{todoID}", func(r chi.Router) {
		r.Use(todoCtx)
		r.Get("/", getTodo)
		r.Put("/", updateTodo)
		r.Delete("/", deleteTodo)
	})
	return r
}
