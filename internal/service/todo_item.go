package service

import (
	"github.com/vbetsun/todo-app/internal/core"
)

type TodoItemStorage interface {
	CreateTodo(listID int, todo core.TodoItem) (int, error)
	GetAllTodos(listID int) ([]core.TodoItem, error)
	GetTodoByID(listID, todoID int) (core.TodoItem, error)
	UpdateTodo(todoID int, data core.UpdateItemData) error
	DeleteTodo(todoID int) error
}

type TodoItemService struct {
	storage TodoItemStorage
}

func NewTodoItemService(storage TodoItemStorage) *TodoItemService {
	return &TodoItemService{storage}
}

func (s *TodoItemService) CreateTodo(listID int, todo core.TodoItem) (int, error) {
	return s.storage.CreateTodo(listID, todo)
}

func (s *TodoItemService) GetAllTodos(listID int) ([]core.TodoItem, error) {
	return s.storage.GetAllTodos(listID)
}

func (s *TodoItemService) GetTodoByID(listID, todoID int) (core.TodoItem, error) {
	return s.storage.GetTodoByID(listID, todoID)
}

func (s *TodoItemService) UpdateTodo(todoID int, data core.UpdateItemData) error {
	return s.storage.UpdateTodo(todoID, data)
}

func (s *TodoItemService) DeleteTodo(todoID int) error {
	return s.storage.DeleteTodo(todoID)
}
