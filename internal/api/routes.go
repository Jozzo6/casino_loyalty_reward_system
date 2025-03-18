package api

import (
	"net/http"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/api/handlers"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/component/promotions"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/component/users"

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

	usersComponent := users.New(s.Resource.DB, s.Resource.PubSub, []byte(s.Resource.Config.JWTKey), s.Resource.Config.JWTDuration)
	promotionsComponent := promotions.New(s.Resource.DB)

	authMiddleware := handlers.AuthMiddleware(usersComponent)

	usersRouter := handlers.NewAccountsRouter(usersComponent)
	promotionsRouter := handlers.NewPromotionsRouter(promotionsComponent)

	r.Route("/api/v1", func(r chi.Router) {

		r.Post("/register", usersRouter.Register())
		r.Post("/login", usersRouter.Login())

		r.With(authMiddleware).Group(func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Get("/", usersRouter.GetUsers())
				r.Get("/{id}", usersRouter.GetUser())
				r.Put("/{id}", usersRouter.UpdateUser())
				r.Put("/{id}/balance", usersRouter.UpdateBalance())
				r.Delete("/{id}", usersRouter.DeleteUser())
				r.Post("/{id}/promotions", usersRouter.AddPromotion())
				r.Get("/{id}/promotions", usersRouter.GetUserPromotions())
				r.Get("/{id}/promotions/{user_prom_id}", usersRouter.GetUserPromotionByID())
				r.Put("/{id}/promotions/{user_prom_id}/claim", usersRouter.ClaimPromotion())
				r.Delete("/{id}/promotions/{user_prom_id}", usersRouter.GetUserPromotionByID())
			})

			r.HandleFunc("/notifications", usersRouter.ListenToNotifications)

			r.Route("/promotions", func(r chi.Router) {
				r.Post("/", promotionsRouter.CreatePromotion())
				r.Get("/", promotionsRouter.GetPromotions())
				r.Get("/{id}", promotionsRouter.GetPromotionByID())
				r.Put("/{id}", promotionsRouter.UpdatePromotion())
				r.Delete("/{id}", promotionsRouter.DeletePromotion())
			})
		})

	})

	return r
}
