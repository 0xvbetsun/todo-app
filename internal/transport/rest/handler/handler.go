package handler

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type Deps struct {
	AuthService     AuthService
	TodoListService TodoListService
	TodoItemService TodoItemService
	Log             *zap.Logger
}

type Handler struct {
	Auth     *AuthHandler
	TodoList *TodoListHandler
	TodoItem *TodoItemHandler
	log      *zap.Logger
}

func New(deps Deps) *Handler {
	return &Handler{
		Auth:     NewAuthHandler(deps.AuthService, deps.Log),
		TodoList: NewTodoListHandler(deps.TodoListService, deps.Log),
		TodoItem: NewTodoItemHandler(deps.TodoItemService, deps.Log),
		log:      deps.Log,
	}
}

func (h *Handler) Routes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(h.Logger())
	r.Use(h.Recoverer)
	r.Use(SendRequestID)
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
			r.Patch("/", h.TodoList.updateList)
			r.Delete("/", h.TodoList.deleteList)
			r.Route("/todos", func(r chi.Router) {
				r.Get("/", h.TodoItem.getAllTodos)
				r.Post("/", h.TodoItem.createTodo)
				r.Route("/{todoID}", func(r chi.Router) {
					r.Use(h.TodoItem.todoCtx)
					r.Get("/", h.TodoItem.getTodo)
					r.Patch("/", h.TodoItem.updateTodo)
					r.Delete("/", h.TodoItem.deleteTodo)
				})
			})
		})
	})
	return r
}
