package service

import (
	"github.com/vbetsun/todo-app/internal/domain"
	"github.com/vbetsun/todo-app/internal/repository"
)

type Authorization interface {
	CreateUser(domain.User) (int, error)
	GenerateToken(username, password string) (string, error)
}
type TodoList interface{}
type TodoItem interface{}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo),
	}
}
