package main

import (
	"go-api/cmd/api"
	"go-api/internal/env"
	initializers "go-api/internal/intializers"
)

func init() {
	initializers.EnvironmentVariables()
	initializers.InitializeDB()
}

func main() {
	port := env.GetString("PORT", ":8080")
	db := initializers.GetDB()
	server := api.NewApiServer(port, db)
	r := server.Init("v1")

	if err := server.Start(r); err != nil {
		panic(err)
	}
}
