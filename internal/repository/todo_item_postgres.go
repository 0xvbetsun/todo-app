package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/vbetsun/todo-app/internal/domain"
)

type TodoItemPostgres struct {
	db *sql.DB
}

func NewTodoItemPostgres(db *sql.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) CreateTodo(listID int, todo domain.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var todoID int
	if err := tx.QueryRow(createTodoQuery(), todo.Title, todo.Description).Scan(&todoID); err != nil {
		tx.Rollback()
		return 0, err
	}
	if _, err := tx.Exec(createListItemsQuery(), listID, todoID); err != nil {
		tx.Rollback()
		return 0, err
	}
	return todoID, tx.Commit()
}

func (r *TodoItemPostgres) GetAllTodos(listID int) ([]domain.TodoItem, error) {
	var todos []domain.TodoItem
	rows, err := r.db.Query(allTodosQuery(), listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var todo domain.TodoItem
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *TodoItemPostgres) GetTodoByID(listID, todoID int) (domain.TodoItem, error) {
	var todo domain.TodoItem
	err := r.db.QueryRow(todoByIDQuery(), listID, todoID).
		Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Done)
	return todo, err
}

func (r *TodoItemPostgres) UpdateTodo(todoID int, data domain.UpdateItemData) error {
	query, args := updateTodo(todoID, data)
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *TodoItemPostgres) DeleteTodo(todoID int) error {
	_, err := r.db.Exec(deleteTodoById(), todoID)
	return err
}

func createTodoQuery() string {
	return fmt.Sprintf(`--sql
		INSERT INTO %s (title, description)
		VALUES ($1, $2)
		RETURNING id
	`, todoItemsTable)
}

func createListItemsQuery() string {
	return fmt.Sprintf(`--sql
		INSERT INTO %s (list_id, item_id)
		VALUES ($1, $2)
	`, listsItemsTable)
}

func allTodosQuery() string {
	return fmt.Sprintf(`--sql
		SELECT ti.id, ti.title, ti.description
		FROM %s AS ti
		INNER JOIN %s AS li ON li.item_id = ti.id
		WHERE li.list_id = $1
	`, todoItemsTable, listsItemsTable)
}

func todoByIDQuery() string {
	return fmt.Sprintf(`--sql
		SELECT ti.id, ti.title, ti.description, ti.done
		FROM %s AS ti
		INNER JOIN %s AS li ON li.item_id = ti.id
		WHERE li.list_id = $1
		AND ti.id = $2
	`, todoItemsTable, listsItemsTable)
}

func updateTodo(todoID int, data domain.UpdateItemData) (string, []interface{}) {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argID := 1
	if data.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argID))
		args = append(args, *data.Title)
		argID++
	}
	if data.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argID))
		args = append(args, *data.Description)
		argID++
	}
	setQuery := strings.Join(setValues, ",")
	args = append(args, todoID)
	return fmt.Sprintf(`--sql
		UPDATE %s
		SET %s
		WHERE id = $%d
	`, todoItemsTable, setQuery, argID), args
}

func deleteTodoById() string {
	return fmt.Sprintf(`--sql
		DELETE FROM %s
		WHERE id = $1
	`, todoItemsTable)
}
