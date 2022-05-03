package repository

import (
	"database/sql"

	"github.com/vbetsun/todo-app/internal/domain"
)

type Authorization interface {
	CreateUser(domain.User) (int, error)
	GetUser(username, password string) (domain.User, error)
}

type TodoList interface {
	Create(userID int, list domain.Todolist) (int, error)
	GetAll(userID int) ([]domain.Todolist, error)
	GetByID(userID, listID int) (domain.Todolist, error)
	Update(listID int, data domain.UpdateListData) error
	Delete(listID int) error
}

type TodoItem interface{}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
	}
}
