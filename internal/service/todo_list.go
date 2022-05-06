package service

import (
	"github.com/vbetsun/todo-app/internal/core"
)

type TodoListStorage interface {
	CreateList(userID int, list core.Todolist) (core.Todolist, error)
	GetAllLists(userID int) ([]core.Todolist, error)
	GetListByID(userID, listID int) (core.Todolist, error)
	UpdateList(listID int, data core.UpdateListData) (core.Todolist, error)
	DeleteList(listID int) error
}
type TodoListService struct {
	storage TodoListStorage
}

func NewTodoListService(storage TodoListStorage) *TodoListService {
	return &TodoListService{storage}
}

func (s *TodoListService) CreateList(userID int, list core.Todolist) (core.Todolist, error) {
	return s.storage.CreateList(userID, list)
}

func (s *TodoListService) GetAllLists(userID int) ([]core.Todolist, error) {
	return s.storage.GetAllLists(userID)
}

func (s *TodoListService) GetListByID(userID, listID int) (core.Todolist, error) {
	return s.storage.GetListByID(userID, listID)
}

func (s *TodoListService) UpdateList(listID int, data core.UpdateListData) (core.Todolist, error) {
	return s.storage.UpdateList(listID, data)
}

func (s *TodoListService) DeleteList(listID int) error {
	return s.storage.DeleteList(listID)
}
