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

const listCtx = "list"

type TodoListService interface {
	CreateList(userID int, list core.Todolist) (int, error)
	GetAllLists(userID int) ([]core.Todolist, error)
	GetListByID(userID, listID int) (core.Todolist, error)
	UpdateList(listID int, data core.UpdateListData) error
	DeleteList(listID int) error
}

type TodoListHandler struct {
	service TodoListService
}

type CreateListRequest struct {
	*core.Todolist
}

type UpdateListRequest struct {
	*core.UpdateListData
}

type AllListsResponse struct {
	Data []core.Todolist `json:"data"`
}

type CreateListResponse struct {
	ID int `json:"id"`
}

type GetListResponse struct {
	*core.Todolist
}

func NewTodoListHandler(service TodoListService) *TodoListHandler {
	return &TodoListHandler{service}
}

func (cl *CreateListRequest) Bind(r *http.Request) error {
	if cl.Title == "" {
		return errors.New("missing required Title field")
	}
	return nil
}

func (ul *UpdateListRequest) Bind(r *http.Request) error {
	if ul.Title == nil && ul.Description == nil {
		return errors.New("you should provide one of Title or Description")
	}
	return nil
}

func (cl *CreateListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (al *AllListsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if len(al.Data) == 0 {
		al.Data = make([]core.Todolist, 0)
	}
	return nil
}

func (gl *GetListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *TodoListHandler) listCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserID(w, r)
		if err != nil {
			render.Render(w, r, ErrInternalServer(err))
			return
		}
		listID, err := strconv.Atoi(chi.URLParam(r, "listID"))
		if err != nil {
			render.Render(w, r, ErrInternalServer(err))
			return
		}
		list, err := h.service.GetListByID(userID, listID)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), listCtx, list)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *TodoListHandler) getAllLists(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(w, r)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	lists, err := h.service.GetAllLists(userID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.Render(w, r, &AllListsResponse{Data: lists})
}

func (h *TodoListHandler) createList(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(w, r)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	data := &CreateListRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	id, err := h.service.CreateList(userID, *data.Todolist)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, &CreateListResponse{ID: id})
}

func (h *TodoListHandler) getList(w http.ResponseWriter, r *http.Request) {
	list, ok := r.Context().Value(listCtx).(core.Todolist)
	if !ok {
		render.Render(w, r, ErrInternalServer(errors.New("listID not found")))
		return
	}

	render.Render(w, r, &GetListResponse{Todolist: &list})
}

func (h *TodoListHandler) updateList(w http.ResponseWriter, r *http.Request) {
	list, ok := r.Context().Value(listCtx).(core.Todolist)
	if !ok {
		render.Render(w, r, ErrInternalServer(errors.New("listID not found")))
		return
	}
	data := &UpdateListRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	err := h.service.UpdateList(list.ID, *data.UpdateListData)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.NoContent(w, r)
}

func (h *TodoListHandler) deleteList(w http.ResponseWriter, r *http.Request) {
	list, ok := r.Context().Value(listCtx).(core.Todolist)
	if !ok {
		render.Render(w, r, ErrInternalServer(errors.New("listID not found")))
		return
	}
	err := h.service.DeleteList(list.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.NoContent(w, r)
}
