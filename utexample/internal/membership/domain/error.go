package domain

import "errors"

var (
	ErrMembershipNotFound = errors.New("membership not found")
	ErrAlreadyMember      = errors.New("already a member")
	ErrInvalidRole        = errors.New("invalid role")
)
