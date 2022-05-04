package psql

import (
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v4/log/zapadapter"
	_ "github.com/jackc/pgx/v4/stdlib"
	"go.uber.org/zap"
)

const (
	usersTable      = "users"
	todoListsTable  = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
	Logger   *zap.Logger
}

type Storage struct {
	Auth     *Auth
	TodoList *TodoList
	TodoItem *TodoItem
}

func (cfg Config) String() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.DBName,
		cfg.Password,
		cfg.SSLMode,
	)
}

func NewDB(cfg Config) (*sql.DB, error) {
	zapadapter.NewLogger(cfg.Logger)
	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Auth:     NewAuth(db),
		TodoList: NewTodoList(db),
		TodoItem: NewTodoItem(db),
	}
}
