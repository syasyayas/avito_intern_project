package service

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

var (
	ErrFeatureAlreadyExists = errors.New("feature already exists")
	ErrFeatureNotFound      = errors.New("feature not found")
	ErrFeatureEmpty         = errors.New("no slug provided")
	ErrFeatureInvalid       = errors.New("invalid feature")
)

var (
	ErrRelationAlreayExists = errors.New("user already has feature")
)
