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
	component promotions.Provider
}

func NewPromotionsRouter(component promotions.Provider) *promotionsRouter {
	return &promotionsRouter{component: component}
}

func (pr *promotionsRouter) CreatePromotion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Promotion

		log := types.GetLoggerFromContext(r.Context())

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		if errs := utils.Validator.Struct(req); errs != nil {
			WriteError(log, w, http.StatusBadRequest, errs)
			return
		}

		promotion, err := pr.component.CreatePromotions(r.Context(), req)
		if err != nil {
			WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(log, w, http.StatusOK, promotion)

	}
}

func (pr *promotionsRouter) GetPromotionByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			log.Errorf("failed to get promotion id: %s", err)
			WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		user, err := pr.component.GetPromotionByID(r.Context(), id)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Errorf("promotion with id: %s was not found: %s", id.String(), err)
			WriteError(log, w, http.StatusNotFound, err)
			return
		}
		if err != nil {
			WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(log, w, http.StatusOK, user)
	}
}

func (pr *promotionsRouter) GetPromotions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		users, err := pr.component.GetPromotions(r.Context())
		if err != nil {
			WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(log, w, http.StatusOK, users)
	}
}

func (pr *promotionsRouter) UpdatePromotion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Promotion

		log := types.GetLoggerFromContext(r.Context())

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		if errs := utils.Validator.Struct(req); errs != nil {
			WriteError(log, w, http.StatusBadRequest, errs)
			return
		}

		promotion, err := pr.component.UpdatePromotion(r.Context(), req)
		if err != nil {
			WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(log, w, http.StatusOK, promotion)

	}
}

func (pr *promotionsRouter) DeletePromotion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			log.Errorf("failed to get promotion id: %s", err)
			WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		err = pr.component.DeletePromotion(r.Context(), id)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Errorf("promotion with id: %s was not found to be deleted: %s", id.String(), err)
			WriteError(log, w, http.StatusNotFound, err)
			return
		}
		if err != nil {
			WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(log, w, http.StatusOK, "OK")
	}
}
