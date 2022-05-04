package service

import (
	"github.com/vbetsun/todo-app/internal/domain"
	"github.com/vbetsun/todo-app/internal/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) CreateList(userID int, list domain.Todolist) (int, error) {
	return s.repo.CreateList(userID, list)
}

func (s *TodoListService) GetAllLists(userID int) ([]domain.Todolist, error) {
	return s.repo.GetAllLists(userID)
}

func (s *TodoListService) GetListByID(userID, listID int) (domain.Todolist, error) {
	return s.repo.GetListByID(userID, listID)
}

func (s *TodoListService) UpdateList(listID int, data domain.UpdateListData) error {
	return s.repo.UpdateList(listID, data)
}

func (s *TodoListService) DeleteList(listID int) error {
	return s.repo.DeleteList(listID)
}
