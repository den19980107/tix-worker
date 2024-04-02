package main

import (
	"fmt"
	"log"
	"os"
	"tix-worker/internal/application"
	"tix-worker/internal/postgres"

	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("ENV")

	log.Printf("env is %s", env)
	if env != "production" {
		err := godotenv.Load()
		if err != nil {
			panic(fmt.Sprintf("failed to load env, err: %s", err))
		}
	}

	db := postgres.Init()
	app := application.New(db)
	app.Run()
}
