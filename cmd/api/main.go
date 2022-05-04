package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"github.com/vbetsun/todo-app/internal/service"
	"github.com/vbetsun/todo-app/internal/storage/psql"
	"github.com/vbetsun/todo-app/internal/transport/rest"
	"github.com/vbetsun/todo-app/internal/transport/rest/handler"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	if err := LoadConfig("configs"); err != nil {
		logger.Fatal(fmt.Sprintf("can't read config: %v", err))
	}
	db, err := psql.NewDB(psql.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		Password: viper.GetString("POSTGRES_PASSWORD"),
		SSLMode:  viper.GetString("db.sslmode"),
		Logger:   logger,
	})
	if err != nil {
		logger.Fatal(fmt.Sprintf("can't connect to the DB %v", err))
	}
	store := psql.NewStorage(db)
	service := service.NewService(service.Deps{
		AuthStorage:     store.Auth,
		TodoListStorage: store.TodoList,
		TodoItemStorage: store.TodoItem,
	})
	h := handler.New(handler.Deps{
		AuthService:     service.Auth,
		TodoListService: service.TodoList,
		TodoItemService: service.TodoItem,
	})
	srv := new(rest.Server)
	port := viper.GetString("port")

	go func() {
		if err := srv.Run(port, h.InitRoutes()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(fmt.Sprintf("can't start server on port %s, err: %v", port, err))
		} else {
			logger.Info("Server stopped gracefully")
		}
	}()
	logger.Info("Server is starting on port: " + port)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-exit
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Error("Error occurred while server is shutting down " + err.Error())
	}
	if err := db.Close(); err != nil {
		logger.Error("Error occurred while db is closing " + err.Error())
	}
}

func LoadConfig(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	viper.AddConfigPath("deployments")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	return viper.MergeInConfig()
}
