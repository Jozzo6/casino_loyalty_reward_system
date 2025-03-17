package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"casino_loyalty_reward_system/internal/component/users"
	"casino_loyalty_reward_system/internal/store"
	"casino_loyalty_reward_system/internal/types"
	utils "casino_loyalty_reward_system/internal/util"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type usersRouter struct {
	component users.Provider
}

func NewAccountsRouter(component users.Provider) *usersRouter {
	return &usersRouter{component: component}
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Token string    `json:"token"`
}

func (ur *usersRouter) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest

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

		createdUser, token, err := ur.component.Register(r.Context(), types.User{
			Email:    req.Email,
			Name:     req.Name,
			Password: req.Password,
		})

		if store.IsErrConflict(err) {
			WriteError(log, w, http.StatusConflict, err)
			return
		} else if err != nil {
			WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(log, w, http.StatusOK, RegisterResponse{
			ID:    createdUser.ID,
			Name:  createdUser.Name,
			Email: createdUser.Email,
			Token: token,
		})
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Token string    `json:"token"`
}

func (ur *usersRouter) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest

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

		user, token, err := ur.component.Login(r.Context(), types.User{
			Email:    req.Email,
			Password: req.Password,
		})
		if errors.Is(err, types.ErrUnauthorized) {
			WriteError(log, w, http.StatusUnauthorized, err)
			return
		}
		if err != nil {
			WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(log, w, http.StatusOK, RegisterResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: token,
		})
	}
}

func (ur *usersRouter) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			log.Errorf("failed to get user id: %s", err)
			WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		user, err := ur.component.GetUser(r.Context(), id)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Errorf("user with id: %s was not found: %s", id.String(), err)
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

func (ur *usersRouter) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		users, err := ur.component.GetUsers(r.Context())
		if err != nil {
			WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(log, w, http.StatusOK, users)
	}
}

func (ur *usersRouter) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.User

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

		user, err := ur.component.UpdateUser(r.Context(), req)
		if err != nil {
			WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(log, w, http.StatusOK, user)

	}
}

func (ur *usersRouter) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			log.Errorf("failed to get user id: %s", err)
			WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		err = ur.component.DeleteUser(r.Context(), id)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Errorf("user with id: %s was not found to be updated: %s", id.String(), err)
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

type UpdateBalanceRequest struct {
	Value           float64               `json:"value" validate:"required"`
	TransactionType types.TransactionType `json:"transaction_type" validate:"required"`
}

func (ur *usersRouter) UpdateBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateBalanceRequest

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

		us, err := types.GetAccountFromContext(r.Context())
		if err != nil {
			WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		user, err := ur.component.UpdateUserBalance(r.Context(), us, req.Value, req.TransactionType)
		if err != nil {
			WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(log, w, http.StatusOK, user)

	}
}

func (ur *usersRouter) AddPromotion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserPromotion

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

		userPromotion, err := ur.component.AddPromotion(r.Context(), req)
		if err != nil {
			log.Errorf("failed to add promotion to user: %s", err)
			if errors.Is(err, types.ErrStartAfterEndDate) ||
				errors.Is(err, types.ErrPromotionNoLongerActive) {
				WriteError(log, w, http.StatusBadRequest, err)
			}
			WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(log, w, http.StatusOK, userPromotion)
	}
}

func (ur *usersRouter) GetUserPromotionByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		userPromotionID, err := uuid.Parse(chi.URLParam(r, "user_prom_id"))
		if err != nil {
			log.Errorf("failed to get user promotion id: %s", err)
			WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		userPromotion, err := ur.component.GetUserPromotionByID(r.Context(), userPromotionID)
		if err != nil {
			WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(log, w, http.StatusOK, userPromotion)
	}
}

func (ur *usersRouter) GetUserPromotions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		userID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			log.Errorf("failed to get user id: %s", err)
			WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		userPromotion, err := ur.component.GetUserPromotions(r.Context(), userID)
		if err != nil {
			WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(log, w, http.StatusOK, userPromotion)
	}
}

func (ur *usersRouter) DeleteUserPromotion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		userPromotionID, err := uuid.Parse(chi.URLParam(r, "user_prom_id"))
		if err != nil {
			log.Errorf("failed to get user promotion id: %s", err)
			WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		err = ur.component.DeleteUserPromotion(r.Context(), userPromotionID)
		if err != nil {
			WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(log, w, http.StatusOK, "ok")
	}
}
func (ur *usersRouter) ClaimPromotion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		userID, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			log.Errorf("failed to get user id: %s", err)
			WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		user, err := types.GetAccountFromContext(r.Context())
		if err != nil {
			log.Errorf("failed to get user from context: %s", err)
			WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		if userID != user.ID {
			log.Error("requestor account ID and path user id not matchting")
			WriteError(log, w, http.StatusBadRequest, types.ErrRequestorIDNotMatching)
			return
		}

		userPromotionID, err := uuid.Parse(chi.URLParam(r, "user_prom_id"))
		if err != nil {
			log.Errorf("failed to get user promotion id: %s", err)
			WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		err = ur.component.ClaimPromotion(r.Context(), userPromotionID)
		if err != nil {
			log.Errorf("failed to get claim promotion: %s", err)
			if errors.Is(err, types.ErrPromotionNoLongerActive) ||
				errors.Is(err, types.ErrPromotionExpired) ||
				errors.Is(err, types.ErrPromotionNotStarted) ||
				errors.Is(err, types.ErrPromotionClaimed) {
				WriteError(log, w, http.StatusBadRequest, err)
				return
			}
			if errors.Is(err, types.ErrRequestorIDNotMatching) {
				WriteError(log, w, http.StatusForbidden, err)
				return
			}
			WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		WriteJSON(log, w, http.StatusOK, "ok")
	}
}
