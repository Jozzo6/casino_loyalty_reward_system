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
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Email      string          `json:"email"`
	Role       UserType        `json:"role,omitempty" `
	Balance    float64         `json:"balance"`
	Created    time.Time       `json:"created"`
	Updated    time.Time       `json:"updated"`
	Promotions []UserPromotion `json:"promotions,omitempty"`
	Password   string
}

type UserType int

const (
	Player UserType = iota
	Staff
)

type TransactionType string

const (
	TransactionTypeRemove TransactionType = "remove"
	TransactionTypeAdd    TransactionType = "add"
)

type UserPromotion struct {
	ID          uuid.UUID  `json:"id"`
	Created     time.Time  `json:"created"`
	Updated     time.Time  `json:"updated"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	Claimed     *time.Time `json:"claimed"`
	UserID      uuid.UUID  `json:"user_id"`
	PromotionID uuid.UUID  `json:"promotion_id"`
	User        *User      `json:"user"`
	Promotion   *Promotion `json:"promotion"`
}
