package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	userpromotion "github.com/Jozzo6/casino_loyalty_reward_system/internal/component/user_promotion"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"
	utils "github.com/Jozzo6/casino_loyalty_reward_system/internal/util"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type userPromotionsRouter struct {
	component userpromotion.UserPromotionProvider
}

func NewUserPromotionsRouter(component userpromotion.UserPromotionProvider) *userPromotionsRouter {
	return &userPromotionsRouter{component: component}
}

func (upr *userPromotionsRouter) AddPromotion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserPromotion

		log := types.GetLoggerFromContext(r.Context())

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			utils.WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		if errs := utils.Validator.Struct(req); errs != nil {
			utils.WriteError(log, w, http.StatusBadRequest, errs)
			return
		}

		userPromotion, err := upr.component.AddPromotion(r.Context(), req)
		if err != nil {
			log.Errorf("failed to add promotion to user: %s", err)
			if errors.Is(err, types.ErrStartAfterEndDate) ||
				errors.Is(err, types.ErrPromotionNoLongerActive) {
				utils.WriteError(log, w, http.StatusBadRequest, err)
			}
			utils.WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(log, w, http.StatusOK, userPromotion)
	}
}

func (upr *userPromotionsRouter) GetUserPromotionByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		userPromotionID, err := uuid.Parse(chi.URLParam(r, "user_prom_id"))
		if err != nil {
			log.Errorf("failed to get user promotion id: %s", err)
			utils.WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		userPromotion, err := upr.component.GetUserPromotionByID(r.Context(), userPromotionID)
		if err != nil {
			utils.WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(log, w, http.StatusOK, userPromotion)
	}
}

func (upr *userPromotionsRouter) GetUserPromotions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
		if err != nil {
			log.Errorf("failed to get user id: %s", err)
			utils.WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		userPromotion, err := upr.component.GetUserPromotions(r.Context(), userID)
		if err != nil {
			utils.WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(log, w, http.StatusOK, userPromotion)
	}
}

func (upr *userPromotionsRouter) DeleteUserPromotion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		userPromotionID, err := uuid.Parse(chi.URLParam(r, "user_prom_id"))
		if err != nil {
			log.Errorf("failed to get user promotion id: %s", err)
			utils.WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		err = upr.component.DeleteUserPromotion(r.Context(), userPromotionID)
		if err != nil {
			utils.WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(log, w, http.StatusOK, "ok")
	}
}
func (upr *userPromotionsRouter) ClaimPromotion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		userID, err := uuid.Parse(chi.URLParam(r, "user_id"))
		if err != nil {
			log.Errorf("failed to get user id: %s", err)
			utils.WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		user, err := types.GetAccountFromContext(r.Context())
		if err != nil {
			log.Errorf("failed to get user from context: %s", err)
			utils.WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		if userID != user.ID {
			log.Error("requestor account ID and path user id not matchting")
			utils.WriteError(log, w, http.StatusBadRequest, types.ErrRequestorIDNotMatching)
			return
		}

		userPromotionID, err := uuid.Parse(chi.URLParam(r, "user_prom_id"))
		if err != nil {
			log.Errorf("failed to get user promotion id: %s", err)
			utils.WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		err = upr.component.ClaimPromotion(r.Context(), userPromotionID)
		if err != nil {
			log.Errorf("failed to get claim promotion: %s", err)
			if errors.Is(err, types.ErrPromotionNoLongerActive) ||
				errors.Is(err, types.ErrPromotionExpired) ||
				errors.Is(err, types.ErrPromotionNotStarted) ||
				errors.Is(err, types.ErrPromotionClaimed) {
				utils.WriteError(log, w, http.StatusBadRequest, err)
				return
			}
			if errors.Is(err, types.ErrRequestorIDNotMatching) {
				utils.WriteError(log, w, http.StatusForbidden, err)
				return
			}
			utils.WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(log, w, http.StatusOK, "ok")
	}
}
