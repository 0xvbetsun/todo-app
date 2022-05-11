package psql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/vbetsun/todo-app/internal/core"
)

// TodoItem represents Todo repository
type TodoItem struct {
	db *sql.DB
}

// NewTodoItem returns instance of Todo repository
func NewTodoItem(db *sql.DB) *TodoItem {
	return &TodoItem{db}
}

// CreateTodo creates new Todo in DB and links it to the List
func (r *TodoItem) CreateTodo(listID int, t core.TodoItem) (core.TodoItem, error) {
	var todo core.TodoItem
	tx, err := r.db.Begin()
	if err != nil {
		return todo, err
	}
	err = tx.QueryRow(createTodoQuery(), t.Title, t.Description).Scan(&todo.ID, &todo.Title, &todo.Description)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return todo, fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return todo, err
	}
	if _, err := tx.Exec(createListItemsQuery(), listID, todo.ID); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return todo, fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return todo, err
	}
	return todo, tx.Commit()
}

// GetAllTodos returns all todos which related to the given List
func (r *TodoItem) GetAllTodos(listID int) ([]core.TodoItem, error) {
	var todos []core.TodoItem
	rows, err := r.db.Query(allTodosQuery(), listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var todo core.TodoItem
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

// GetTodoByID returns todo by ID which related to the given list
func (r *TodoItem) GetTodoByID(listID, todoID int) (core.TodoItem, error) {
	var todo core.TodoItem
	err := r.db.QueryRow(todoByIDQuery(), listID, todoID).
		Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Done)
	return todo, err
}

// UpdateTodo save Todo changes to the db
func (r *TodoItem) UpdateTodo(todoID int, data core.UpdateItemData) (core.TodoItem, error) {
	var t core.TodoItem
	query, args := updateTodo(todoID, data)
	err := r.db.QueryRow(query, args...).Scan(&t.ID, &t.Title, &t.Description, &t.Done)
	return t, err
}

// DeleteTodo removes todo from DB by ID
func (r *TodoItem) DeleteTodo(todoID int) error {
	_, err := r.db.Exec(deleteTodoById(), todoID)
	return err
}

func createTodoQuery() string {
	return fmt.Sprintf(`--sql
		INSERT INTO %s (title, description)
		VALUES ($1, $2)
		RETURNING id, title, description
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

func updateTodo(todoID int, data core.UpdateItemData) (string, []interface{}) {
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
		RETURNING id, title, description, done
	`, todoItemsTable, setQuery, argID), args
}

func deleteTodoById() string {
	return fmt.Sprintf(`--sql
		DELETE FROM %s
		WHERE id = $1
	`, todoItemsTable)
}
