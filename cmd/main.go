package main

import (
	"stt/api/route"
	"stt/bootstrap"
	db "stt/database/postgres/sqlc"
	"time"
)

func main() {
	// set up bootstrap for db connection and environment variable, configuration
	app := bootstrap.NewServerApp(".")

	// Set up database store for quering and handle db transaction (core)
	dbStore := db.NewStore(app.PostgresConnectionPool)
	timeout := time.Duration(app.Env.ContextTimeout) * time.Second

	route.Setup(app.Env, timeout, dbStore, app.Engine)

	app.Engine.Run(app.Env.ServerAddress)

	defer app.CloseDbConnection()

}
