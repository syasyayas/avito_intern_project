package model

import "time"

type History struct {
	UserID    string    `json:"user_id" csv:"user_id"`
	FeatureID int       `json:"feature_id" csv:"feature_id"`
	Operation string    `json:"operation" csv:"operation"`
	Date      time.Time `json:"date" csv:"date"`
}
