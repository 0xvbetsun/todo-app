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

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
