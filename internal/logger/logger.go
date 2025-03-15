package logger

import (
	"casino_loyalty_reward_system/internal/types"
	"context"
	"net/http"

	"go.uber.org/zap"
)

func Logger(log *zap.SugaredLogger, projectID string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := log.With()

			ctx := context.WithValue(r.Context(), types.CtxKeyLogger, log)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
