package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/vbetsun/todo-app/internal/domain"
)

const todoCtx = "todo"

type CreateTodoRequest struct {
	*domain.TodoItem
}

type UpdateTodoRequest struct {
	*domain.UpdateItemData
}

type CreateTodoResponse struct {
	ID int `json:"id"`
}

type AllTodosResponse struct {
	Data []domain.TodoItem `json:"data"`
}

type GetTodoResponse struct {
	*domain.TodoItem
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
		at.Data = make([]domain.TodoItem, 0)
	}
	return nil
}

func (ct *CreateTodoResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (ct *GetTodoResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *Handler) todoCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		list, ok := r.Context().Value(listCtx).(domain.Todolist)
		if !ok {
			render.Render(w, r, ErrInternalServer(errors.New("listID not found")))
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

func (h *Handler) getAllTodos(w http.ResponseWriter, r *http.Request) {
	list, ok := r.Context().Value(listCtx).(domain.Todolist)
	if !ok {
		render.Render(w, r, ErrInternalServer(errors.New("listID not found")))
		return
	}
	todos, err := h.service.GetAllTodos(list.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.Render(w, r, &AllTodosResponse{Data: todos})
}

func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) {
	list, ok := r.Context().Value(listCtx).(domain.Todolist)
	if !ok {
		render.Render(w, r, ErrInternalServer(errors.New("listID not found")))
		return
	}
	data := &CreateTodoRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	id, err := h.service.CreateTodo(list.ID, *data.TodoItem)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, &CreateTodoResponse{ID: id})
}

func (h *Handler) getTodo(w http.ResponseWriter, r *http.Request) {
	todo, ok := r.Context().Value(todoCtx).(domain.TodoItem)
	if !ok {
		render.Render(w, r, ErrInternalServer(errors.New("todoID not found")))
		return
	}

	render.Render(w, r, &GetTodoResponse{TodoItem: &todo})
}

func (h *Handler) updateTodo(w http.ResponseWriter, r *http.Request) {
	todo, ok := r.Context().Value(todoCtx).(domain.TodoItem)
	if !ok {
		render.Render(w, r, ErrInternalServer(errors.New("todoID not found")))
		return
	}
	data := &UpdateTodoRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	err := h.service.UpdateTodo(todo.ID, *data.UpdateItemData)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.NoContent(w, r)
}

func (h *Handler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	todo, ok := r.Context().Value(todoCtx).(domain.TodoItem)
	if !ok {
		render.Render(w, r, ErrInternalServer(errors.New("todoID not found")))
		return
	}
	err := h.service.DeleteTodo(todo.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.NoContent(w, r)
}
