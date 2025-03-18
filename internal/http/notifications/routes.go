package notifications

import (
	"net/http"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/component/users"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/http/middlewares"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/http/notifications/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *server) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middlewares.LoggerMiddleware(s.Resource.Log))

	usersComponent := users.New(s.Resource.DB, s.Resource.PubSub, []byte(s.Resource.Config.JWTKey), s.Resource.Config.JWTDuration)

	authMiddleware := middlewares.AuthMiddleware(usersComponent)

	notificationRouter := handlers.NewNotificationsRouter(usersComponent)

	r.With(authMiddleware).Route("/api/v1", func(r chi.Router) {
		r.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})
		r.HandleFunc("/notifications", notificationRouter.ListenToNotifications)
	})

	return r
}
