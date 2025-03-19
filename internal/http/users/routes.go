package users

import (
	"net/http"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/component/users"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/http/middlewares"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/http/users/handlers"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"

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

	usersRouter := handlers.NewAccountsRouter(usersComponent)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", usersRouter.Register())
		r.Post("/login", usersRouter.Login())

		r.With(authMiddleware).Group(func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", usersRouter.GetUsers())
				r.Get("/{id}", usersRouter.GetUser())
				r.Put("/{id}", usersRouter.UpdateUser())
				r.Put("/{id}/balance", usersRouter.UpdateBalance())
				r.With(middlewares.RequiredRole(types.Staff)).Delete("/{id}", usersRouter.DeleteUser())
			})
		})
	})

	return r
}
