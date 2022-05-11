package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/vbetsun/todo-app/internal/core"
	"go.uber.org/zap"
)

// Key to use when setting the todo context.
type ctxKeyTodo string

const todoCtx ctxKeyTodo = "todo"

type TodoItemService interface {
	CreateTodo(listID int, todo core.TodoItem) (core.TodoItem, error)
	GetAllTodos(listID int) ([]core.TodoItem, error)
	GetTodoByID(listID, todoID int) (core.TodoItem, error)
	UpdateTodo(todoID int, data core.UpdateItemData) (core.TodoItem, error)
	DeleteTodo(todoID int) error
}

type TodoItemHandler struct {
	service TodoItemService
	log     *zap.Logger
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

func NewTodoItemHandler(service TodoItemService, log *zap.Logger) *TodoItemHandler {
	return &TodoItemHandler{service, log}
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
			if err := render.Render(w, r, ErrInternalServer(ErrListNotFound)); err != nil {
				h.log.Error(ErrRenderResp.Error())
			}
			return
		}
		todoID, err := strconv.Atoi(chi.URLParam(r, "todoID"))
		if err != nil {
			if rErr := render.Render(w, r, ErrInternalServer(err)); rErr != nil {
				h.log.Error(ErrRenderResp.Error())
			}
			return
		}
		todo, err := h.service.GetTodoByID(list.ID, todoID)
		if err != nil {
			if rErr := render.Render(w, r, ErrNotFound); rErr != nil {
				h.log.Error(ErrRenderResp.Error())
			}
			return
		}
		ctx := context.WithValue(r.Context(), todoCtx, todo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *TodoItemHandler) getAllTodos(w http.ResponseWriter, r *http.Request) {
	list, ok := r.Context().Value(listCtx).(core.Todolist)
	if !ok {
		if err := render.Render(w, r, ErrInternalServer(ErrListNotFound)); err != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	todos, err := h.service.GetAllTodos(list.ID)
	if err != nil {
		if rErr := render.Render(w, r, ErrInternalServer(err)); rErr != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	if err := render.Render(w, r, &AllTodosResponse{Data: todos}); err != nil {
		h.log.Error(ErrRenderResp.Error())
	}
}

func (h *TodoItemHandler) createTodo(w http.ResponseWriter, r *http.Request) {
	list, ok := r.Context().Value(listCtx).(core.Todolist)
	if !ok {
		if err := render.Render(w, r, ErrInternalServer(ErrListNotFound)); err != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	data := &CreateTodoRequest{}
	if err := render.Bind(r, data); err != nil {
		if rErr := render.Render(w, r, ErrInvalidRequest(err)); rErr != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	todo, err := h.service.CreateTodo(list.ID, *data.TodoItem)
	if err != nil {
		if rErr := render.Render(w, r, ErrInternalServer(err)); rErr != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	render.Status(r, http.StatusCreated)
	if err := render.Render(w, r, &TodoResponse{TodoItem: &todo}); err != nil {
		h.log.Error(ErrRenderResp.Error())
	}
}

func (h *TodoItemHandler) getTodo(w http.ResponseWriter, r *http.Request) {
	todo, ok := r.Context().Value(todoCtx).(core.TodoItem)
	if !ok {
		if err := render.Render(w, r, ErrInternalServer(ErrTodoNotFound)); err != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	if err := render.Render(w, r, &TodoResponse{TodoItem: &todo}); err != nil {
		h.log.Error(ErrRenderResp.Error())
	}
}

func (h *TodoItemHandler) updateTodo(w http.ResponseWriter, r *http.Request) {
	todo, ok := r.Context().Value(todoCtx).(core.TodoItem)
	if !ok {
		if err := render.Render(w, r, ErrInternalServer(ErrTodoNotFound)); err != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	data := &UpdateTodoRequest{}
	if err := render.Bind(r, data); err != nil {
		if rErr := render.Render(w, r, ErrInvalidRequest(err)); rErr != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	todo, err := h.service.UpdateTodo(todo.ID, *data.UpdateItemData)
	if err != nil {
		if rErr := render.Render(w, r, ErrInvalidRequest(err)); rErr != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	if err := render.Render(w, r, &TodoResponse{TodoItem: &todo}); err != nil {
		h.log.Error(ErrRenderResp.Error())
	}
}

func (h *TodoItemHandler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	todo, ok := r.Context().Value(todoCtx).(core.TodoItem)
	if !ok {
		if err := render.Render(w, r, ErrInternalServer(ErrTodoNotFound)); err != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	err := h.service.DeleteTodo(todo.ID)
	if err != nil {
		if rErr := render.Render(w, r, ErrInternalServer(err)); rErr != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	render.NoContent(w, r)
}
