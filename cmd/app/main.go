package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/vbetsun/todo-app/internal/repository"
	"github.com/vbetsun/todo-app/internal/service"
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
	if err := initConfig(); err != nil {
		logger.Fatal(fmt.Sprintf("can't read config: %v", err))
	}
	db, err := repository.NewPostgres(repository.Config{
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
	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	h := handler.NewHandler(service)
	srv := new(rest.Server)
	port := viper.GetString("port")

	logger.Info("Server is starting on port: " + port)
	if err := srv.Run(port, h.InitRoutes()); err != nil {
		logger.Fatal(fmt.Sprintf("can't start server on port %s, err: %v", port, err))
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	viper.SetConfigFile(".env")
	return viper.MergeInConfig()
}
