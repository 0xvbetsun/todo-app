package repository

import (
	"database/sql"

	"github.com/vbetsun/todo-app/internal/domain"
)

type Authorization interface {
	CreateUser(domain.User) (int, error)
	GetUser(username, password string) (domain.User, error)
}

type TodoList interface{}
type TodoItem interface{}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
