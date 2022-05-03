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

func (s *TodoListService) Create(userID int, list domain.Todolist) (int, error) {
	return s.repo.Create(userID, list)
}

func (s *TodoListService) GetAll(userID int) ([]domain.Todolist, error) {
	return s.repo.GetAll(userID)
}

func (s *TodoListService) GetByID(userID, listID int) (domain.Todolist, error) {
	return s.repo.GetByID(userID, listID)
}

func (s *TodoListService) Update(listID int, data domain.UpdateListData) error {
	return s.repo.Update(listID, data)
}

func (s *TodoListService) Delete(listID int) error {
	return s.repo.Delete(listID)
}
