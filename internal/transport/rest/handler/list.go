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

const listCtx = "list"

type CreateListRequest struct {
	*domain.Todolist
}

type UpdateListRequest struct {
	*domain.UpdateListData
}

type AllListsResponse struct {
	Data []domain.Todolist `json:"data"`
}

type CreateListResponse struct {
	ID int `json:"id"`
}

type GetListResponse struct {
	*domain.Todolist
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
	return nil
}

func (gl *GetListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *Handler) listCtx(next http.Handler) http.Handler {
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
		list, err := h.service.TodoList.GetByID(userID, listID)
		if err != nil {
			render.Render(w, r, ErrNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), listCtx, list)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) getAllLists(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(w, r)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	lists, err := h.service.TodoList.GetAll(userID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.Render(w, r, &AllListsResponse{Data: lists})
}

func (h *Handler) createList(w http.ResponseWriter, r *http.Request) {
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
	id, err := h.service.TodoList.Create(userID, *data.Todolist)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, &CreateListResponse{ID: id})
}

func (h *Handler) getList(w http.ResponseWriter, r *http.Request) {
	list, ok := r.Context().Value(listCtx).(domain.Todolist)
	if !ok {
		render.Render(w, r, ErrInternalServer(errors.New("listID not found")))
		return
	}

	render.Render(w, r, &GetListResponse{Todolist: &list})
}

func (h *Handler) updateList(w http.ResponseWriter, r *http.Request) {
	list, ok := r.Context().Value(listCtx).(domain.Todolist)
	if !ok {
		render.Render(w, r, ErrInternalServer(errors.New("listID not found")))
		return
	}
	data := &UpdateListRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	err := h.service.TodoList.Update(list.ID, *data.UpdateListData)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	render.NoContent(w, r)
}

func (h *Handler) deleteList(w http.ResponseWriter, r *http.Request) {
	list, ok := r.Context().Value(listCtx).(domain.Todolist)
	if !ok {
		render.Render(w, r, ErrInternalServer(errors.New("listID not found")))
		return
	}
	err := h.service.TodoList.Delete(list.ID)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}
	render.NoContent(w, r)
}
