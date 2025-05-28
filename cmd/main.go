package main

import (
	"github.com/gin-gonic/gin"
	"twitter-api/cmd/http"
	"twitter-api/internal/clients/sql"
	"twitter-api/internal/config"
	"twitter-api/internal/feeds"
	"twitter-api/internal/tweets"
	"twitter-api/internal/users"
)

func main() {
	cfg := config.LoadConfig()
	engine := gin.Default()
	redis := cfg.Redis()
	postgres := cfg.Postgres()
	postgresClient := sql.NewSqlClient(postgres)

	usersRepository := users.NewRepository(postgresClient, redis)
	usersService := users.NewService(usersRepository)
	usersHandler := users.NewHandler(usersService)

	tweetsRepository := tweets.NewRepository(postgresClient, redis)
	tweetsService := tweets.NewService(tweetsRepository)
	tweetsHandler := tweets.NewHandler(tweetsService)

	feedsRepository := feeds.NewRepository(redis)
	feedsService := feeds.NewService(feedsRepository)
	go feedsService.CreateFeeds()

	httpServer := http.NewServer(engine, usersHandler, tweetsHandler)
	httpServer.Run(cfg.Port)
}
