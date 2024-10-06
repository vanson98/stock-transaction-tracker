package main

import (
	"stt/api/route"
	"stt/bootstrap"
	db "stt/database/postgres/sqlc"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// set up bootstrap for db connection and environment variable, configuration
	app := bootstrap.App(".")
	timeout := time.Duration(app.Env.ContextTimeout) * time.Second
	defer app.CloseDbConnection()

	// Set up database store for quering and handle transaction (core)
	dbStore := db.NewStore(app.PostgresConnectionPool)

	// Handle routing to controllers
	gin := gin.Default()
	route.Setup(app.Env, timeout, dbStore, gin)

	gin.Run(app.Env.ServerAddress)
}
