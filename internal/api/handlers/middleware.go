package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/component/users"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"

	"go.uber.org/zap"
)

func LoggerMiddleware(log *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := log.With()

			ctx := context.WithValue(r.Context(), types.CtxKeyLogger, log)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthMiddleware(component users.Provider) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			log := types.GetLoggerFromContext(ctx)
			token := r.Header.Get("Authorization")

			if token == "" {
				WriteErrorMessage(log, w, http.StatusUnauthorized, "missing token in header")
				return
			}

			parts := strings.Split(token, " ")
			if len(parts) != 2 {
				WriteErrorMessage(log, w, http.StatusUnauthorized, "invalid token in header")
				return
			}

			if parts[0] != "Bearer" {
				WriteErrorMessage(log, w, http.StatusUnauthorized, "not token bearer")
				return
			}

			account, err := component.Auth(ctx, parts[1], r.URL.Path, r.Method)
			if err != nil {
				log.Errorf("failed to auth user: %s", err)
				WriteErrorMessage(log, w, http.StatusUnauthorized, "failed to auth user")
				return
			}

			ctx = context.WithValue(ctx, types.CtxKeyAccount, account)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
