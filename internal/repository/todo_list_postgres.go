package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/vbetsun/todo-app/internal/domain"
)

type TodoListPostgres struct {
	db *sql.DB
}

func NewTodoListPostgres(db *sql.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) CreateList(userID int, l domain.Todolist) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var id int
	if err := tx.QueryRow(createListQuery(), l.Title, l.Description).Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	if _, err := tx.Exec(createUsersListQuery(), userID, id); err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAllLists(userID int) ([]domain.Todolist, error) {
	var lists []domain.Todolist
	rows, err := r.db.Query(allListsQuery(), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var list domain.Todolist
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

func (r *TodoListPostgres) GetListByID(userID, listID int) (domain.Todolist, error) {
	var list domain.Todolist
	err := r.db.QueryRow(listByIDQuery(), userID, listID).Scan(&list.ID, &list.Title, &list.Description)
	return list, err
}

func (r *TodoListPostgres) UpdateList(listID int, data domain.UpdateListData) error {
	query, args := updateList(listID, data)
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *TodoListPostgres) DeleteList(listID int) error {
	_, err := r.db.Exec(deleteListById(), listID)
	return err
}

func createUsersListQuery() string {
	return fmt.Sprintf(`--sql
		INSERT INTO %s (user_id, list_id) 
		VALUES ($1, $2) RETURNING id
	`, usersListsTable)
}

func createListQuery() string {
	return fmt.Sprintf(`--sql
		INSERT INTO %s (title, description) 
		VALUES ($1, $2) RETURNING id
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

func updateList(listID int, data domain.UpdateListData) (string, []interface{}) {
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
	`, todoListsTable, setQuery, argID), args
}

func deleteListById() string {
	return fmt.Sprintf(`--sql
		DELETE FROM %s
		WHERE id = $1
	`, todoListsTable)
}
