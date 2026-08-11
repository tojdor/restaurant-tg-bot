package models

type User struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Password_hash string `json:"password_hash"`
	Role          string `json:"role"`
}
