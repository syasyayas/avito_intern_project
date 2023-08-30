package model

import "time"

// features
type NewFeatureRequest struct {
	Slug    string `json:"slug"`
	Percent int    `json:"percent"`
}

type DeleteFeatureRequest struct {
	Slug string `json:"slug"`
}

type FeatureSlugOnly struct {
	Slug string `json:"slug"`
}
type FeatureSlugAndExpire struct {
	Slug      string    `json:"slug"`
	ExpiresAt time.Time `json:"expires_at"`
}

type AddFeaturesToUserRequest struct {
	UserId   string                 `json:"user_id"`
	Features []FeatureSlugAndExpire `json:"features"`
}

type DeleteFeaturesFromUser struct {
	UserId   string            `json:"user_id"`
	Features []FeatureSlugOnly `json:"features"`
}

// user

type UserRequest struct {
	ID string `json:"id"`
}

type UserWithFeaturesResponse struct {
	ID       string                 `json:"id"`
	Features []FeatureSlugAndExpire `json:"features"`
}

// history

type HistoryRequest struct {
	After  time.Time `json:"after"`
	Before time.Time `json:"before"`
}

type HistoryResponse struct {
	URL string `json:"url"`
}
