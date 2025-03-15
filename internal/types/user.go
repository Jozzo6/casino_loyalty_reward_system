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
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

type UserType int

const (
	Admin UserType = iota
	Staff
	Player
)
