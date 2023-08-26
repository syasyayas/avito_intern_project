package model

import "time"

type Feature struct {
	ID        int       `json:"id" db:"id"`
	Slug      string    `json:"slug" db:"slug"`
	ExpiresAt time.Time `json:"expires_at"`
}
