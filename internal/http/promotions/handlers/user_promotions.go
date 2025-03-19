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

// AddPromotion adds a promotion to a user.
// @Summary Add a promotion to a user
// @Description Assign a promotion to a user with the provided details
// @Tags User Promotions
// @Accept json
// @Produce json
// @Param userPromotion body types.UserPromotion true "User Promotion details"
// @Success 200 {object} types.UserPromotion "Added user promotion"
// @Failure 400 {object} types.ErrorResponse "Invalid input or business rule violation"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/user-promotions/{user_id} [post]
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

// GetUserPromotionByID retrieves a user promotion by its ID.
// @Summary Get a user promotion by ID
// @Description Retrieve a user promotion using its unique ID
// @Tags User Promotions
// @Accept json
// @Produce json
// @Param user_prom_id path string true "User Promotion ID"
// @Success 200 {object} types.UserPromotion "Retrieved user promotion"
// @Failure 400 {object} types.ErrorResponse "Invalid ID format"
// @Failure 404 {object} types.ErrorResponse "User promotion not found"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/user-promotions/{user_id}/promotion/{user_prom_id} [get]
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

// GetUserPromotions retrieves all promotions for a specific user.
// @Summary Get all promotions for a user
// @Description Retrieve a list of all promotions assigned to a specific user
// @Tags User Promotions
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {array} types.UserPromotion "List of user promotions"
// @Failure 400 {object} types.ErrorResponse "Invalid ID format"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/user-promotions/{user_id} [get]
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

// DeleteUserPromotion deletes a user promotion by its ID.
// @Summary Delete a user promotion
// @Description Delete a user promotion using its unique ID
// @Tags User Promotions
// @Accept json
// @Produce json
// @Param user_prom_id path string true "User Promotion ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} types.ErrorResponse "Invalid ID format"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/user-promotions/{user_id}/promotions/{user_prom_id} [delete]
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

// ClaimPromotion allows a user to claim a promotion.
// @Summary Claim a promotion
// @Description Allows a user to claim a promotion if eligible
// @Tags User Promotions
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param user_prom_id path string true "User Promotion ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} types.ErrorResponse "Invalid input or business rule violation"
// @Failure 403 {object} types.ErrorResponse "Forbidden - Requestor ID does not match"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/user-promotions/{user_id}/promotions/{user_prom_id}/claim [post]
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
