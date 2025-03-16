package types

import "errors"

type ErrorResponse struct {
	Message string `json:"message"`
}

var (
	ErrUnauthorized        = errors.New("Unauthorized")
	ErrInsufficientBalance = errors.New("Insufficient balance")
)
