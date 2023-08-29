package model

import "time"

type Feature struct {
	Slug      string    `json:"slug" db:"slug"`
	ExpiresAt time.Time `json:"expires_at"`
}
