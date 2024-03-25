package postgres

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	pgUser := os.Getenv("POSTGRES_USER")
	pgPwd := os.Getenv("POSTGRES_PASSWORD")
	pgHost := os.Getenv("POSTGRES_HOST")
	pgPort := os.Getenv("POSTGRES_PORT")
	pgDatabase := os.Getenv("POSTGRES_DATABASE")

	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", pgUser, pgPwd, pgHost, pgPort, pgDatabase)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic(fmt.Sprintf("failed to init postgres, err: %s", err))
	}

	return db
}
