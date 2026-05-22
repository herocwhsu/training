package domain

import "errors"

var (
	ErrCompanyNotFound = errors.New("company not found")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidName     = errors.New("invalid name")
)
