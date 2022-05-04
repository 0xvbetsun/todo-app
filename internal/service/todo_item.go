package service

import (
	"github.com/vbetsun/todo-app/internal/domain"
	"github.com/vbetsun/todo-app/internal/repository"
)

type TodoItemService struct {
	repo repository.TodoItem
}

func NewTodoItemService(repo repository.TodoItem) *TodoItemService {
	return &TodoItemService{repo: repo}
}

func (s *TodoItemService) CreateTodo(listID int, todo domain.TodoItem) (int, error) {
	return s.repo.CreateTodo(listID, todo)
}

func (s *TodoItemService) GetAllTodos(listID int) ([]domain.TodoItem, error) {
	return s.repo.GetAllTodos(listID)
}

func (s *TodoItemService) GetTodoByID(listID, todoID int) (domain.TodoItem, error) {
	return s.repo.GetTodoByID(listID, todoID)
}

func (s *TodoItemService) UpdateTodo(todoID int, data domain.UpdateItemData) error {
	return s.repo.UpdateTodo(todoID, data)
}

func (s *TodoItemService) DeleteTodo(todoID int) error {
	return s.repo.DeleteTodo(todoID)
}
