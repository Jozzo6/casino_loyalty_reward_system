package types

import (
	"time"

	"github.com/google/uuid"
)

type UserFilter struct {
	ByID    uuid.NullUUID
	ByEmail *string
}

type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password,omitempty"`
	Role     UserType  `json:"role,omitempty" `
	Balance  float64   `json:"balance"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

type UserType int

const (
	Admin UserType = iota
	Staff
	Player
)

type TransactionType string

const (
	Remove TransactionType = "remove"
	Add    TransactionType = "add"
)
