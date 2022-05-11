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

// Key to use when setting the list context.
type ctxKeyList string

const listCtx ctxKeyList = "list"

type TodoListService interface {
	CreateList(userID int, list core.Todolist) (core.Todolist, error)
	GetAllLists(userID int) ([]core.Todolist, error)
	GetListByID(userID, listID int) (core.Todolist, error)
	UpdateList(listID int, data core.UpdateListData) (core.Todolist, error)
	DeleteList(listID int) error
}

type TodoListHandler struct {
	service TodoListService
	log     *zap.Logger
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

type ListResponse struct {
	*core.Todolist
}

func NewTodoListHandler(service TodoListService, log *zap.Logger) *TodoListHandler {
	return &TodoListHandler{service, log}
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

func (cl *ListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (al *AllListsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if len(al.Data) == 0 {
		al.Data = make([]core.Todolist, 0)
	}
	return nil
}

func (h *TodoListHandler) listCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserID(w, r)
		if err != nil {
			if rErr := render.Render(w, r, ErrInternalServer(err)); rErr != nil {
				h.log.Error(ErrRenderResp.Error())
			}
			return
		}
		listID, err := strconv.Atoi(chi.URLParam(r, "listID"))
		if err != nil {
			if rErr := render.Render(w, r, ErrInternalServer(err)); rErr != nil {
				h.log.Error(ErrRenderResp.Error())
			}
			return
		}
		list, err := h.service.GetListByID(userID, listID)
		if err != nil {
			if rErr := render.Render(w, r, ErrNotFound); rErr != nil {
				h.log.Error(ErrRenderResp.Error())
			}
			return
		}
		ctx := context.WithValue(r.Context(), listCtx, list)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *TodoListHandler) getAllLists(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(w, r)
	if err != nil {
		if rErr := render.Render(w, r, ErrInternalServer(err)); rErr != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	lists, err := h.service.GetAllLists(userID)
	if err != nil {
		if rErr := render.Render(w, r, ErrInternalServer(err)); rErr != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	if err := render.Render(w, r, &AllListsResponse{Data: lists}); err != nil {
		h.log.Error(ErrRenderResp.Error())
	}
}

func (h *TodoListHandler) createList(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(w, r)
	if err != nil {
		if rErr := render.Render(w, r, ErrInternalServer(err)); rErr != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	data := &CreateListRequest{}
	if err := render.Bind(r, data); err != nil {
		if rErr := render.Render(w, r, ErrInvalidRequest(err)); rErr != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	list, err := h.service.CreateList(userID, *data.Todolist)
	if err != nil {
		if rErr := render.Render(w, r, ErrInternalServer(err)); rErr != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	render.Status(r, http.StatusCreated)
	if err := render.Render(w, r, &ListResponse{Todolist: &list}); err != nil {
		h.log.Error(ErrRenderResp.Error())
	}
}

func (h *TodoListHandler) getList(w http.ResponseWriter, r *http.Request) {
	list, ok := r.Context().Value(listCtx).(core.Todolist)
	if !ok {
		if err := render.Render(w, r, ErrInternalServer(ErrListNotFound)); err != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	if err := render.Render(w, r, &ListResponse{Todolist: &list}); err != nil {
		h.log.Error(ErrRenderResp.Error())
	}
}

func (h *TodoListHandler) updateList(w http.ResponseWriter, r *http.Request) {
	list, ok := r.Context().Value(listCtx).(core.Todolist)
	if !ok {
		if err := render.Render(w, r, ErrInternalServer(ErrListNotFound)); err != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	data := &UpdateListRequest{}
	if err := render.Bind(r, data); err != nil {
		if rErr := render.Render(w, r, ErrInvalidRequest(err)); rErr != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	list, err := h.service.UpdateList(list.ID, *data.UpdateListData)
	if err != nil {
		if rErr := render.Render(w, r, ErrInvalidRequest(err)); rErr != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	if err := render.Render(w, r, &ListResponse{Todolist: &list}); err != nil {
		h.log.Error(ErrRenderResp.Error())
	}
}

func (h *TodoListHandler) deleteList(w http.ResponseWriter, r *http.Request) {
	list, ok := r.Context().Value(listCtx).(core.Todolist)
	if !ok {
		if err := render.Render(w, r, ErrInternalServer(ErrListNotFound)); err != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	err := h.service.DeleteList(list.ID)
	if err != nil {
		if rErr := render.Render(w, r, ErrInternalServer(err)); rErr != nil {
			h.log.Error(ErrRenderResp.Error())
		}
		return
	}
	render.NoContent(w, r)
}
