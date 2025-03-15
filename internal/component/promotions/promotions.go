package promotions

import (
	"time"

	"github.com/google/uuid"
)

type Promotion struct {
	PromotionID uuid.UUID `json:"promotion_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	IsActive    bool      `json:"is_active"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}
