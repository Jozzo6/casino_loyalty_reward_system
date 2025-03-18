package types

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

type ctxKey int

const (
	CtxKeyAccount ctxKey = iota
	CtxKeyLogger
	CtxUserRole
)

func GetLoggerFromContext(ctx context.Context) *zap.SugaredLogger {
	logger, ok := ctx.Value(CtxKeyLogger).(*zap.SugaredLogger)
	if ok {
		return logger
	}

	newLogger, _ := zap.NewProduction()
	return newLogger.Sugar()
}

func GetAccountFromContext(ctx context.Context) (User, error) {
	user, ok := ctx.Value(CtxKeyAccount).(User)
	if ok {
		return user, nil
	}

	return User{}, errors.New("user not in ctx")
}

func GetUserRoleFromContext(ctx context.Context) (UserType, error) {
	userRole, ok := ctx.Value(CtxUserRole).(UserType)
	if ok {
		return userRole, nil
	}

	return 0, errors.New("failed to get user role from context")
}
