package main

import (
	"stt/api/route"
	"stt/bootstrap"
	db "stt/database/postgres/sqlc"
	"stt/services"
	"time"
)

func main() {
	// set up bootstrap for db connection and environment variable, configuration
	app := bootstrap.NewServerApp(".")

	// Set up database store for quering and handle db transaction (core)
	dbStore := db.NewStore(app.PostgresConnectionPool)
	timeout := time.Duration(app.Env.ContextTimeout) * time.Second

	accountService := services.InitAccountService(dbStore, timeout)

	protectedRouterGroup := app.Engine.Group("")
	route.InitAccountRouter(protectedRouterGroup, accountService)

	app.Engine.Run(app.Env.ServerAddress)

	defer app.CloseDbConnection()

}
