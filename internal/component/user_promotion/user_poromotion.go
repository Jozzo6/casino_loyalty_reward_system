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

type UserPromotionProvider interface {
	AddPromotion(ctx context.Context, userPromotion types.UserPromotion) (types.UserPromotion, error)
	AddWelcomePromotion(ctx context.Context, userID uuid.UUID) (types.UserPromotion, error)
	GetUserPromotions(ctx context.Context, userID uuid.UUID) ([]types.UserPromotion, error)
	GetUserPromotionByID(ctx context.Context, userPromotionID uuid.UUID) (types.UserPromotion, error)
	ClaimPromotion(ctx context.Context, userPromotionID uuid.UUID) error
	DeleteUserPromotion(ctx context.Context, userPromotionID uuid.UUID) error
	ListenToRegisterEvent(ctx context.Context) error
}

type component struct {
	persistent store.Persistent
	pubsub     store.PubSub
}

var _ UserPromotionProvider = (*component)(nil)

func New(persistent store.Persistent, pubsub store.PubSub) *component {
	comp := &component{
		persistent: persistent,
		pubsub:     pubsub,
	}

	go func() {
		err := comp.ListenToRegisterEvent(context.Background())
		if err != nil {
			fmt.Printf("error in ListenToRegisterEvent: %v", err)
		}
	}()

	return comp
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

func (c *component) ListenToRegisterEvent(ctx context.Context) error {
	sub := c.pubsub.Subscribe(ctx, fmt.Sprintf(redis_pub_sub.RegistrationChannel))
	defer sub.Close()

	log := types.GetLoggerFromContext(ctx)

	ch := sub.Channel()

	for msg := range ch {
		ID, err := uuid.Parse(msg.Payload)
		if err != nil {
			log.Errorf("failed to parse uuid: %s", err)
			continue
		}
		_, err = c.AddWelcomePromotion(ctx, ID)
		if err != nil {
			log.Errorf("failed to add welcome promotion: %s", err)
			continue
		}
	}

	return nil
}
