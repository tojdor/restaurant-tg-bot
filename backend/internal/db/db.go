package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func Connect() (*pgxpool.Pool, error) {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dsn := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Error while trying to connect to DB")
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Printf("Error pinging db")
		return nil, err
	}

	log.Printf("Connected to db")

	return pool, nil
}
