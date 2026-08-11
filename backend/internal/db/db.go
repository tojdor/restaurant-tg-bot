package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func Connect() (*pgx.Conn, error) {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dsn := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal("Error while trying to connect to DB")
		return nil, err
	}

	log.Printf("Connected to db")

	return conn, nil
}
