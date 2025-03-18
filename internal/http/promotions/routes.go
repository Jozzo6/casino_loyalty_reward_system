package promotions

import (
	"net/http"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/component/promotions"
	userpromotion "github.com/Jozzo6/casino_loyalty_reward_system/internal/component/user_promotion"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/component/users"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/http/middlewares"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/http/promotions/handlers"
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
	promotionsComponent := promotions.New(s.Resource.DB)
	userPromotionComponent := userpromotion.New(s.Resource.DB, s.Resource.PubSub)

	authMiddleware := middlewares.AuthMiddleware(usersComponent)

	promotionsRouter := handlers.NewPromotionsRouter(promotionsComponent)
	userPromotionsRouter := handlers.NewUserPromotionsRouter(userPromotionComponent)

	r.Route("/api/v1", func(r chi.Router) {
		r.With(authMiddleware).Group(func(r chi.Router) {
			r.Route("/user_promotions", func(r chi.Router) {
				r.Get("/{user_id}", userPromotionsRouter.GetUserPromotions())
				r.Get("/{user_id}/promotion/{user_prom_id}", userPromotionsRouter.GetUserPromotionByID())
				r.Put("/{user_id}/promotions/{user_prom_id}/claim", userPromotionsRouter.ClaimPromotion())
				
				r.With(middlewares.RequiredRole(types.Staff)).Group(func(r chi.Router) {
					r.Post("/{user_id}", userPromotionsRouter.AddPromotion())
					r.Delete("/{user_id}/promotions/{user_prom_id}", userPromotionsRouter.GetUserPromotionByID())
				})
			})

			r.Route("/promotions", func(r chi.Router) {
				r.Get("/", promotionsRouter.GetPromotions())
				r.Get("/{id}", promotionsRouter.GetPromotionByID())
				r.With(middlewares.RequiredRole(types.Staff)).Group(func(r chi.Router) {
					r.Post("/", promotionsRouter.CreatePromotion())
					r.Put("/{id}", promotionsRouter.UpdatePromotion())
					r.Delete("/{id}", promotionsRouter.DeletePromotion())
				})
			})
		})
	})

	return r
}
