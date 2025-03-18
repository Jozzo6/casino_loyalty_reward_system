package notifications

import (
	"context"
	"fmt"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/store"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/store/redis_pub_sub"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type Provider interface {
	ListenToNotifications(ctx context.Context, conn *websocket.Conn, userID uuid.UUID) error
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

func (c *component) ListenToNotifications(ctx context.Context, conn *websocket.Conn, userID uuid.UUID) error {
	sub := c.pubsub.Subscribe(ctx, fmt.Sprintf("%s:%s", redis_pub_sub.NotificationsChannel, userID))
	defer sub.Close()

	ch := sub.Channel()

	for msg := range ch {
		err := conn.Write(ctx, websocket.MessageText, []byte(msg.Payload))
		if err != nil {
			return err
		}
	}

	return nil
}
