package psql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/vbetsun/todo-app/internal/core"
)

type TodoList struct {
	db *sql.DB
}

func NewTodoList(db *sql.DB) *TodoList {
	return &TodoList{db}
}

func (r *TodoList) CreateList(userID int, l core.Todolist) (core.Todolist, error) {
	var list core.Todolist
	tx, err := r.db.Begin()
	if err != nil {
		return list, err
	}
	err = tx.QueryRow(createListQuery(), l.Title, l.Description).Scan(&list.ID, &list.Title, &list.Description)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return list, fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return list, err
	}
	if _, err := tx.Exec(createUsersListQuery(), userID, list.ID); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return list, fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return list, err
	}
	return list, tx.Commit()
}

func (r *TodoList) GetAllLists(userID int) ([]core.Todolist, error) {
	var lists []core.Todolist
	rows, err := r.db.Query(allListsQuery(), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var list core.Todolist
		if err := rows.Scan(&list.ID, &list.Title, &list.Description); err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return lists, nil
}

func (r *TodoList) GetListByID(userID, listID int) (core.Todolist, error) {
	var list core.Todolist
	err := r.db.QueryRow(listByIDQuery(), userID, listID).Scan(&list.ID, &list.Title, &list.Description)
	return list, err
}

func (r *TodoList) UpdateList(listID int, data core.UpdateListData) (core.Todolist, error) {
	var list core.Todolist
	query, args := updateList(listID, data)
	err := r.db.QueryRow(query, args...).Scan(&list.ID, &list.Title, &list.Description)
	return list, err
}

func (r *TodoList) DeleteList(listID int) error {
	_, err := r.db.Exec(deleteListById(), listID)
	return err
}

func createUsersListQuery() string {
	return fmt.Sprintf(`--sql
		INSERT INTO %s (user_id, list_id) 
		VALUES ($1, $2) 
		RETURNING id
	`, usersListsTable)
}

func createListQuery() string {
	return fmt.Sprintf(`--sql
		INSERT INTO %s (title, description) 
		VALUES ($1, $2) 
		RETURNING id, title, description
	`, todoListsTable)
}

func allListsQuery() string {
	return fmt.Sprintf(`--sql
		SELECT tl.id, tl.title, tl.description 
		FROM %s AS tl 
		INNER JOIN %s AS ul ON tl.id = ul.list_id 
		WHERE ul.user_id = $1
	`, todoListsTable, usersListsTable)
}

func listByIDQuery() string {
	return fmt.Sprintf(`--sql
		SELECT tl.id, tl.title, tl.description 
		FROM %s AS tl 
		INNER JOIN %s AS ul ON tl.id = ul.list_id 
		WHERE ul.user_id = $1
		AND tl.id = $2
	`, todoListsTable, usersListsTable)
}

func updateList(listID int, data core.UpdateListData) (string, []interface{}) {
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
	args = append(args, listID)
	return fmt.Sprintf(`--sql
		UPDATE %s
		SET %s
		WHERE id = $%d
		RETURNING id, title, description
	`, todoListsTable, setQuery, argID), args
}

func deleteListById() string {
	return fmt.Sprintf(`--sql
		DELETE FROM %s
		WHERE id = $1
	`, todoListsTable)
}
