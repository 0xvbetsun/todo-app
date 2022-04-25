package repository

import (
	"database/sql"
	"fmt"

	"github.com/vbetsun/todo-app/internal/domain"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(u domain.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable)
	if err := r.db.QueryRow(query, u.Name, u.Username, u.Password).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, pwd string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.QueryRow(query, username, pwd).Scan(&user.ID)

	return user, err
}
