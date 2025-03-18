package promotions

import (
	"context"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/store"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"

	"github.com/google/uuid"
)

type Provider interface {
	CreatePromotions(ctx context.Context, promotion types.Promotion) (types.Promotion, error)
	GetPromotions(ctx context.Context) ([]types.Promotion, error)
	GetPromotionByID(ctx context.Context, ID uuid.UUID) (types.Promotion, error)
	UpdatePromotion(ctx context.Context, promotion types.Promotion) (types.Promotion, error)
	DeletePromotion(ctx context.Context, ID uuid.UUID) error
}

type component struct {
	persistent store.Persistent
}

var _ Provider = (*component)(nil)

func New(persistent store.Persistent) *component {
	return &component{persistent: persistent}
}

func (c *component) CreatePromotions(ctx context.Context, promotion types.Promotion) (types.Promotion, error) {
	promotion.ID = uuid.New()

	createdPromotion, err := c.persistent.PromotionCreate(ctx, promotion)
	if err != nil {
		return types.Promotion{}, err
	}

	return createdPromotion, err
}

func (c *component) GetPromotionByID(ctx context.Context, ID uuid.UUID) (types.Promotion, error) {
	return c.persistent.PromotionGetByID(ctx, ID)
}

func (c *component) GetPromotions(ctx context.Context) ([]types.Promotion, error) {
	return c.persistent.GetPromotions(ctx)
}

func (c *component) UpdatePromotion(ctx context.Context, promotion types.Promotion) (types.Promotion, error) {
	return c.persistent.PromotionUpdate(ctx, promotion)
}

func (c *component) DeletePromotion(ctx context.Context, ID uuid.UUID) error {
	return c.persistent.PromotionDelete(ctx, ID)
}
