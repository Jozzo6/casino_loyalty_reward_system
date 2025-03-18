package userpromotion

import (
	"context"
	"fmt"
	"time"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/store"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/store/redis_pub_sub"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"
	"github.com/google/uuid"
)

type Provider interface {
	AddPromotion(ctx context.Context, userPromotion types.UserPromotion) (types.UserPromotion, error)
	AddWelcomePromotion(ctx context.Context, userID uuid.UUID) (types.UserPromotion, error)
	GetUserPromotions(ctx context.Context, userID uuid.UUID) ([]types.UserPromotion, error)
	GetUserPromotionByID(ctx context.Context, userPromotionID uuid.UUID) (types.UserPromotion, error)
	ClaimPromotion(ctx context.Context, userPromotionID uuid.UUID) error
	DeleteUserPromotion(ctx context.Context, userPromotionID uuid.UUID) error
}

type component struct {
	persistent store.Persistent
	pubsub     store.PubSub
}

var _ Provider = (*component)(nil)

func New(persistent store.Persistent, pubsub store.PubSub) *component {
	return &component{
		persistent: persistent,
		pubsub:     pubsub,
	}
}

func (c *component) AddPromotion(ctx context.Context, userPromotion types.UserPromotion) (types.UserPromotion, error) {
	userPromotion.ID = uuid.New()

	if userPromotion.StartDate.After(userPromotion.EndDate) {
		return types.UserPromotion{}, types.ErrStartAfterEndDate
	}

	promotion, err := c.persistent.PromotionGetByID(ctx, userPromotion.PromotionID)
	if err != nil {
		return types.UserPromotion{}, err
	}

	if !promotion.IsActive {
		return types.UserPromotion{}, types.ErrPromotionNoLongerActive
	}

	up, err := c.persistent.AddPromotion(ctx, userPromotion)
	if err != nil {
		return types.UserPromotion{}, err
	}

	c.pubsub.Publish(ctx, fmt.Sprintf("%s:%s", redis_pub_sub.NotificationsChannel, up.UserID.String()), up)

	return up, err
}

func (c *component) AddWelcomePromotion(ctx context.Context, userID uuid.UUID) (types.UserPromotion, error) {
	promotion, err := c.persistent.PromotionGetByType(ctx, types.WelcomeBonus)
	if err != nil {
		return types.UserPromotion{}, err
	}

	if !promotion.IsActive {
		return types.UserPromotion{}, types.ErrPromotionNoLongerActive
	}

	uP := types.UserPromotion{
		ID:          uuid.New(),
		UserID:      userID,
		PromotionID: promotion.ID,
		StartDate:   time.Now(),
		EndDate:     time.Now().Add(24 * time.Hour),
		Created:     time.Now(),
		Updated:     time.Now(),
	}

	userPromotion, err := c.persistent.AddPromotion(ctx, uP)
	if err != nil {
		return types.UserPromotion{}, err
	}

	c.pubsub.Publish(ctx, fmt.Sprintf("%s:%s", redis_pub_sub.NotificationsChannel, userID.String()), userPromotion)

	return userPromotion, err
}

func (c *component) ClaimPromotion(ctx context.Context, userPromotionID uuid.UUID) error {
	db, err := c.persistent.WithTx(ctx)

	userPromotion, err := db.GetUserPromotionByID(ctx, userPromotionID)
	if err != nil {
		return err
	}

	if userPromotion.Claimed != nil {
		return types.ErrPromotionClaimed
	}

	if !userPromotion.Promotion.IsActive {
		return types.ErrPromotionNoLongerActive
	}

	if time.Now().Before(userPromotion.StartDate) {
		return types.ErrPromotionNotStarted
	}

	if time.Now().After(userPromotion.EndDate) {
		return types.ErrPromotionExpired
	}

	err = db.ClaimPromotion(ctx, userPromotion.ID)
	if err != nil {
		return err
	}

	user, err := db.UserGetBy(ctx, types.UserFilter{ByID: uuid.NullUUID{UUID: userPromotion.UserID, Valid: true}})
	if err != nil {
		return err
	}

	db.UserBalanceUpdate(ctx, user.ID, user.Balance+userPromotion.Promotion.Amount)

	_, err = db.UserBalanceUpdate(ctx, user.ID, userPromotion.Promotion.Amount)
	return err
}

func (c *component) DeleteUserPromotion(ctx context.Context, userPromotionID uuid.UUID) error {
	return c.persistent.DeleteUserPromotion(ctx, userPromotionID)
}

func (c *component) GetUserPromotionByID(ctx context.Context, userPromotionID uuid.UUID) (types.UserPromotion, error) {
	return c.persistent.GetUserPromotionByID(ctx, userPromotionID)
}

func (c *component) GetUserPromotions(ctx context.Context, userPromotionID uuid.UUID) ([]types.UserPromotion, error) {
	return c.persistent.GetUserPromotions(ctx, userPromotionID)
}
