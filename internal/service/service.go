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
	Create(userID int, list domain.Todolist) (int, error)
	GetAll(userID int) ([]domain.Todolist, error)
	GetByID(userID, listID int) (domain.Todolist, error)
	Update(listID int, data domain.UpdateListData) error
	Delete(listID int) error
}
type TodoItem interface{}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo),
		TodoList:      NewTodoListService(repo),
	}
}
