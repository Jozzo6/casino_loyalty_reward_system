package redis_pub_sub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func New(client *redis.Client, log *zap.SugaredLogger) *PubSub {
	return &PubSub{client: client, log: log}
}

type PubSub struct {
	log    *zap.SugaredLogger
	client *redis.Client
}

func (ps *PubSub) Publish(ctx context.Context, channel string, data any) *redis.IntCmd {
	bytes, err := json.Marshal(data)
	if err != nil {
		ps.log.Warn(fmt.Sprintf("err marshalling pubsub publish data %s", err.Error()))
	}

	return ps.client.Publish(ctx, channel, string(bytes))
}

func (ps *PubSub) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return ps.client.Subscribe(ctx, channel)
}
//cd2ef8d2-da53-4a46-8981-4ac4902ea283