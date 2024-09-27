package main

import (
	"stt/api/route"
	"stt/bootstrap"
	db "stt/database/postgres/sqlc"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	app := bootstrap.App()
	env := app.Env

	connectionPool := app.PostgresConnectionPool
	queries := db.New(connectionPool)
	defer app.CloseDbConnection()

	timeout := time.Duration(env.ContextTimeout) * time.Second

	gin := gin.Default()

	route.Setup(env, timeout, queries, gin)

	gin.Run(env.ServerAddress)
}
