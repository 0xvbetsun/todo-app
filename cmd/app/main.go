package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/vbetsun/todo-app/internal/repository"
	"github.com/vbetsun/todo-app/internal/service"
	"github.com/vbetsun/todo-app/internal/transport/rest"
	"github.com/vbetsun/todo-app/internal/transport/rest/handler"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("Can not read config %s", err.Error())
	}
	db, err := repository.NewPostgres(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		Password: viper.GetString("POSTGRES_PASSWORD"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("Can not connect to the DB %s", err.Error())
	}
	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	h := handler.NewHandler(service)
	srv := new(rest.Server)
	port := viper.GetString("port")

	log.Printf("Server is starting on port %s", port)
	if err := srv.Run(port, h.InitRoutes()); err != nil {
		log.Fatalf("Can not start server on port %s, err: %s", port, err.Error())
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
