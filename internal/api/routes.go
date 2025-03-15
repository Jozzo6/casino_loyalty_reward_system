package api

import (
	"casino_loyalty_reward_system/internal/api/handlers"
	"casino_loyalty_reward_system/internal/component/users"
	"net/http"

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

	r.Use(handlers.LoggerMiddleware(s.Resource.Log))

	usersComponent := users.New(s.Resource.DB, []byte(s.Resource.Config.JWTKey), s.Resource.Config.JWTDuration)

	authMiddleware := handlers.AuthMiddleware(usersComponent)

	usersRouter := handlers.NewAccountsRouter(usersComponent)

	r.Route("/api/v1", func(r chi.Router) {

		r.Post("/register", usersRouter.Register())
		r.Post("/login", usersRouter.Login())

		r.With(authMiddleware).Route("/users", func(r chi.Router) {
		})

	})

	return r
}
