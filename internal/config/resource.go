package config

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/store"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/store/postgresdb"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/store/redis_pub_sub"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Resource struct {
	Config     *Config
	Log        *zap.SugaredLogger
	HTTPClient *http.Client
	DB         store.Persistent
	PubSub     store.PubSub
	Close      func() error
}

func InitResource(ctx context.Context) (*Resource, error) {
	var r Resource

	c, err := newConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}
	r.Config = c

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development: r.Config.Environment == "staging",
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	r.Log = logger.Sugar()

	r.HTTPClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	pool, closer, err := postgresdb.PgxPool(ctx, r.Config.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Postgres connection pool: %w", err)
	}

	r.DB = postgresdb.New(r.Log, pool)

	redisOpt, err := redis.ParseURL(r.Config.RedisURI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URI: %w", err)
	}

	redisClient := redis.NewClient(redisOpt)
	err = redisClient.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	r.PubSub = redis_pub_sub.New(redisClient, r.Log)

	r.Close = func() error {
		return errors.Join(
			closer(),
			redisClient.Close(),
		)
	}

	return &r, nil
}
