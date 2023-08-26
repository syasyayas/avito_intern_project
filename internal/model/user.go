package model

type User struct {
	ID       string    `json:"id" db:"id"`
	Features []Feature `json:"features"`
}
