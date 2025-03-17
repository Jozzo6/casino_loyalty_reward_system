package types

import (
	"time"

	"github.com/google/uuid"
)

type Promotion struct {
	ID          uuid.UUID     `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Amount      float64       `json:"amount"`
	IsActive    bool          `json:"is_active"`
	Type        PromotionType `json:"type"`
	Created     time.Time     `json:"created"`
	Updated     time.Time     `json:"updated"`
}

type PromotionType string

const (
	Regular      PromotionType = "regular"
	WelcomeBonus PromotionType = "welcome_bonus"
)
