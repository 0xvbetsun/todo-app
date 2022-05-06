package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/vbetsun/todo-app/internal/core"
)

const todoCtx = "todo"

type TodoItemService interface {
	CreateTodo(listID int, todo core.TodoItem) (core.TodoItem, error)
	GetAllTodos(listID int) ([]core.TodoItem, error)
	GetTodoByID(listID, todoID int) (core.TodoItem, error)
	UpdateTodo(todoID int, data core.UpdateItemData) (core.TodoItem, error)
	DeleteTodo(todoID int) error
}

type TodoItemHandler struct {
	service TodoItemService
}

type CreateTodoRequest struct {
	*core.TodoItem
}

type UpdateTodoRequest struct {
	*core.UpdateItemData
}

type AllTodosResponse struct {
	Data []core.TodoItem `json:"data"`
}

type TodoResponse struct {
	*core.TodoItem
}

func NewTodoItemHandler(service TodoItemService) *TodoItemHandler {
	return &TodoItemHandler{service}
}

func (ct *CreateTodoRequest) Bind(r *http.Request) error {
	if ct.Title == "" {
		return errors.New("missing required Title field")
	}
	return nil
}

func (ut *UpdateTodoRequest) Bind(r *http.Request) error {
	if ut.Title == nil && ut.Description == nil {
		return errors.New("you should provide one of Title or Description")
	}
	return nil
}

func (at *AllTodosResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if len(at.Data) == 0 {
		at.Data = make([]core.TodoItem, 0)
	}
	return nil
}

func (ct *TodoResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *TodoItemHandler) todoCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		list, ok := r.Context().Value(listCtx).(core.Todolist)
		if !ok {
			render.Render(w, r, ErrInternalServer(ErrListNotFound))
			return
		}
		todoID, err := strconv.Atoi(chi.URLParam(r, "todoID"))
		if err != nil {
			render.Render(w, r, ErrInternalServer(err))
			return
		}
		todo, err := h.service.GetTodoByID(list.ID, todoID)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), todoCtx, todo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *TodoItemHandler) getAllTodos(w http.ResponseWriter, r *http.Request) {
	list, ok := r.Context().Value(listCtx).(core.Todolist)
	if !ok {
		render.Render(w, r, ErrInternalServer(ErrListNotFound))
		return
	}
	todos, err := h.service.GetAllTodos(list.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.Render(w, r, &AllTodosResponse{Data: todos})
}

func (h *TodoItemHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	list, ok := r.Context().Value(listCtx).(core.Todolist)
	if !ok {
		render.Render(w, r, ErrInternalServer(ErrListNotFound))
		return
	}
	data := &CreateTodoRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	todo, err := h.service.CreateTodo(list.ID, *data.TodoItem)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, &TodoResponse{TodoItem: &todo})
}

func (h *TodoItemHandler) getTodo(w http.ResponseWriter, r *http.Request) {
	todo, ok := r.Context().Value(todoCtx).(core.TodoItem)
	if !ok {
		render.Render(w, r, ErrInternalServer(ErrTodoNotFound))
		return
	}

	render.Render(w, r, &TodoResponse{TodoItem: &todo})
}

func (h *TodoItemHandler) updateTodo(w http.ResponseWriter, r *http.Request) {
	todo, ok := r.Context().Value(todoCtx).(core.TodoItem)
	if !ok {
		render.Render(w, r, ErrInternalServer(ErrTodoNotFound))
		return
	}
	data := &UpdateTodoRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	todo, err := h.service.UpdateTodo(todo.ID, *data.UpdateItemData)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, &TodoResponse{TodoItem: &todo})
}

func (h *TodoItemHandler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	todo, ok := r.Context().Value(todoCtx).(core.TodoItem)
	if !ok {
		render.Render(w, r, ErrInternalServer(ErrTodoNotFound))
		return
	}
	err := h.service.DeleteTodo(todo.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.NoContent(w, r)
}
