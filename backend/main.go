package main

import (
	"backend/internal/db"
	"backend/internal/user"
	"log"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db.RunMigrations()

	pool, err := db.Connect()
	if err != nil {
		log.Fatal("Error while trying to connect to db")
	}
	defer pool.Close()

	userStorage := user.NewStorage(pool)
	userService := user.NewService(userStorage)
	userHandler := user.NewHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("/login", userHandler.IsRegisteredHandler)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
