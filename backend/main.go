package main

import (
	"backend/internal/db"
	"log"
)

func main() {

	pool, err := db.Connect()

	if err != nil {
		log.Fatal("Error while trying to connect to db")
	}

	defer pool.Close()

	// userStorage := user.NewStorage(pool)
	// ctx := context.Background()
}
