package handler

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
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
	r.Use(middleware.RequestID)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Mount("/auth", h.authRouter())
	r.With(h.UserIdentity).Mount("/api", h.apiRouter())
	return r
}

func (h *Handler) authRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/sign-in", h.SignIn)
	r.Post("/sign-up", h.SignUp)
	return r
}

func (h *Handler) apiRouter() chi.Router {
	r := chi.NewRouter()
	r.Route("/lists", func(r chi.Router) {
		r.Get("/", h.getAllLists)
		r.Post("/", h.createList)
		r.Route("/{listID}", func(r chi.Router) {
			r.Use(h.listCtx)
			r.Get("/", h.getList)
			r.Put("/", h.updateList)
			r.Delete("/", h.deleteList)
			r.Route("/todos", func(r chi.Router) {
				r.Get("/", h.getAllTodos)
				r.Post("/", h.createTodo)
				r.Route("/{todoID}", func(r chi.Router) {
					r.Use(h.todoCtx)
					r.Get("/", h.getTodo)
					r.Put("/", h.updateTodo)
					r.Delete("/", h.deleteTodo)
				})
			})
		})
	})
	return r
}
