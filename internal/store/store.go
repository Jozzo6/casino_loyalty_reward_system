package store

import (
	"context"
	"database/sql"
	"errors"

	"casino_loyalty_reward_system/internal/types"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
)

type Tx interface {
	WithTx(ctx context.Context) (Persistent, error)
	CommitTx(ctx context.Context) error
	RollbackTx(ctx context.Context) error
}

type UserManager interface {
	UserCreate(ctx context.Context, user types.User) (types.User, error)
	UserGetBy(ctx context.Context, filter types.UserFilter) (types.User, error)
	GetUsers(ctx context.Context) ([]types.User, error)
}

type Persistent interface {
	Tx
	UserManager
}

type Cache interface {
}

type PubSub interface {
	Publish(ctx context.Context, channel string, data any) *redis.IntCmd
	Subscribe(ctx context.Context, channel string) *redis.PubSub
}

func IsErrNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) || errors.Is(err, redis.Nil)
}

func IsErrConflict(err error) bool {
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return pgErr.Code == "23505"
		}
	}
	return false
}
