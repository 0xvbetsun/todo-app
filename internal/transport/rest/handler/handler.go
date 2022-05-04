package handler

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Deps struct {
	AuthService     AuthService
	TodoListService TodoListService
	TodoItemService TodoItemService
}

type Handler struct {
	Auth     *AuthHandler
	TodoList *TodoListHandler
	TodoItem *TodoItemHandler
}

func New(deps Deps) *Handler {
	return &Handler{
		Auth:     NewAuthHandler(deps.AuthService),
		TodoList: NewTodoListHandler(deps.TodoListService),
		TodoItem: NewTodoItemHandler(deps.TodoItemService),
	}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.URLFormat)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Mount("/auth", h.authRouter())
	r.With(h.Auth.UserIdentity).Mount("/api", h.apiRouter())
	return r
}

func (h *Handler) authRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/sign-up", h.Auth.SignUp)
	r.Post("/sign-in", h.Auth.SignIn)
	return r
}

func (h *Handler) apiRouter() chi.Router {
	r := chi.NewRouter()
	r.Route("/lists", func(r chi.Router) {
		r.Get("/", h.TodoList.getAllLists)
		r.Post("/", h.TodoList.createList)
		r.Route("/{listID}", func(r chi.Router) {
			r.Use(h.TodoList.listCtx)
			r.Get("/", h.TodoList.getList)
			r.Put("/", h.TodoList.updateList)
			r.Delete("/", h.TodoList.deleteList)
			r.Route("/todos", func(r chi.Router) {
				r.Get("/", h.TodoItem.getAllTodos)
				r.Post("/", h.TodoItem.createTodo)
				r.Route("/{todoID}", func(r chi.Router) {
					r.Use(h.TodoItem.todoCtx)
					r.Get("/", h.TodoItem.getTodo)
					r.Put("/", h.TodoItem.updateTodo)
					r.Delete("/", h.TodoItem.deleteTodo)
				})
			})
		})
	})
	return r
}
