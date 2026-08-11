package main

import (
	"backend/internal/db"
	"fmt"
)

func main() {

	fmt.Println(db.Connect())
}
