package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/component/users"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"
	utils "github.com/Jozzo6/casino_loyalty_reward_system/internal/util"

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

func AuthMiddleware(component users.UserProvider) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			log := types.GetLoggerFromContext(ctx)
			token := r.Header.Get("Authorization")

			if token == "" {
				utils.WriteErrorMessage(log, w, http.StatusUnauthorized, "missing token in header")
				return
			}

			parts := strings.Split(token, " ")
			if len(parts) != 2 {
				utils.WriteErrorMessage(log, w, http.StatusUnauthorized, "invalid token in header")
				return
			}

			if parts[0] != "Bearer" {
				utils.WriteErrorMessage(log, w, http.StatusUnauthorized, "not token bearer")
				return
			}

			account, err := component.Auth(ctx, parts[1], r.URL.Path, r.Method)
			if err != nil {
				log.Errorf("failed to auth user: %s", err)
				utils.WriteErrorMessage(log, w, http.StatusUnauthorized, "failed to auth user")
				return
			}

			ctx = context.WithValue(ctx, types.CtxKeyAccount, account)
			ctx = context.WithValue(ctx, types.CtxUserRole, account.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequiredRole(requiredRole types.UserType) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := types.GetLoggerFromContext(ctx).With("handler", "middleware.required_role")

			userRole, err := types.GetUserRoleFromContext(ctx)
			if err != nil {
				log.Errorf("failed to get user role from context: %s", err)
				utils.WriteError(log, w, http.StatusInternalServerError, err)
				return
			}

			if userRole < requiredRole {
				log.Infof("user does not have required role: %s", requiredRole)
				utils.WriteErrorMessage(log, w, http.StatusForbidden, "user does not have required role")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
