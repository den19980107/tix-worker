package main

import (
	"fmt"
	"tix-worker/internal/application"
	"tix-worker/internal/postgres"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load env, err: %s", err))
	}
	db := postgres.Init()
	app := application.New(db)
	app.Run()
}
