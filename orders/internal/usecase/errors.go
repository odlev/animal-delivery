// Package usecase is a nice package
package usecase

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrInvalidUUID = errors.New("invalid uuid")
	ErrInvalidStatus = errors.New("invalid status")
)
