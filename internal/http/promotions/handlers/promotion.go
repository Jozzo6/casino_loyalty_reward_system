package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/component/promotions"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"
	utils "github.com/Jozzo6/casino_loyalty_reward_system/internal/util"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type promotionsRouter struct {
	component promotions.PromotionProvider
}

func NewPromotionsRouter(component promotions.PromotionProvider) *promotionsRouter {
	return &promotionsRouter{component: component}
}

// CreatePromotion handles the creation of a new promotion.
// @Summary Create a new promotion
// @Description Create a new promotion with the provided details
// @Tags Promotions
// @Accept json
// @Produce json
// @Param promotion body types.Promotion true "Promotion details"
// @Success 200 {object} types.Promotion "Created promotion"
// @Failure 400 {object} types.ErrorResponse "Invalid input"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/promotions [post]
func (pr *promotionsRouter) CreatePromotion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Promotion

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

		promotion, err := pr.component.CreatePromotions(r.Context(), req)
		if err != nil {
			utils.WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(log, w, http.StatusOK, promotion)

	}
}

// GetPromotionByID retrieves a promotion by its ID.
// @Summary Get a promotion by ID
// @Description Retrieve a promotion using its unique ID
// @Tags Promotions
// @Accept json
// @Produce json
// @Param id path string true "Promotion ID"
// @Success 200 {object} types.Promotion "Retrieved promotion"
// @Failure 400 {object} types.ErrorResponse "Invalid ID format"
// @Failure 404 {object} types.ErrorResponse "Promotion not found"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/promotions/{id} [get]
func (pr *promotionsRouter) GetPromotionByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			log.Errorf("failed to get promotion id: %s", err)
			utils.WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		user, err := pr.component.GetPromotionByID(r.Context(), id)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Errorf("promotion with id: %s was not found: %s", id.String(), err)
			utils.WriteError(log, w, http.StatusNotFound, err)
			return
		}
		if err != nil {
			utils.WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(log, w, http.StatusOK, user)
	}
}

// GetPromotions retrieves all promotions.
// @Summary Get all promotions
// @Description Retrieve a list of all promotions
// @Tags Promotions
// @Accept json
// @Produce json
// @Success 200 {array} types.Promotion "List of promotions"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/promotions [get]
func (pr *promotionsRouter) GetPromotions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		users, err := pr.component.GetPromotions(r.Context())
		if err != nil {
			utils.WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(log, w, http.StatusOK, users)
	}
}

// UpdatePromotion updates an existing promotion.
// @Summary Update a promotion
// @Description Update an existing promotion with the provided details
// @Tags Promotions
// @Accept json
// @Produce json
// @Param promotion body types.Promotion true "Updated promotion details"
// @Success 200 {object} types.Promotion "Updated promotion"
// @Failure 400 {object} types.ErrorResponse "Invalid input"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/promotions/{id} [put]
func (pr *promotionsRouter) UpdatePromotion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Promotion

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

		promotion, err := pr.component.UpdatePromotion(r.Context(), req)
		if err != nil {
			utils.WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(log, w, http.StatusOK, promotion)

	}
}

// DeletePromotion deletes a promotion by its ID.
// @Summary Delete a promotion
// @Description Delete a promotion using its unique ID
// @Tags Promotions
// @Accept json
// @Produce json
// @Param id path string true "Promotion ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} types.ErrorResponse "Invalid ID format"
// @Failure 404 {object} types.ErrorResponse "Promotion not found"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/promotions/{id} [delete]
func (pr *promotionsRouter) DeletePromotion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			log.Errorf("failed to get promotion id: %s", err)
			utils.WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		err = pr.component.DeletePromotion(r.Context(), id)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Errorf("promotion with id: %s was not found to be deleted: %s", id.String(), err)
			utils.WriteError(log, w, http.StatusNotFound, err)
			return
		}
		if err != nil {
			utils.WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(log, w, http.StatusOK, "OK")
	}
}
