package service

import (
	"github.com/vbetsun/todo-app/internal/domain"
	"github.com/vbetsun/todo-app/internal/repository"
)

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}
type TodoList interface {
	CreateList(userID int, list domain.Todolist) (int, error)
	GetAllLists(userID int) ([]domain.Todolist, error)
	GetListByID(userID, listID int) (domain.Todolist, error)
	UpdateList(listID int, data domain.UpdateListData) error
	DeleteList(listID int) error
}
type TodoItem interface {
	CreateTodo(listID int, todo domain.TodoItem) (int, error)
	GetAllTodos(listID int) ([]domain.TodoItem, error)
	GetTodoByID(listID, todoID int) (domain.TodoItem, error)
	UpdateTodo(todoID int, data domain.UpdateItemData) error
	DeleteTodo(todoID int) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo),
		TodoList:      NewTodoListService(repo),
		TodoItem:      NewTodoItemService(repo),
	}
}
