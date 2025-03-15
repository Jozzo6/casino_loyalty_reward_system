package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"casino_loyalty_reward_system/internal/component/users"
	"casino_loyalty_reward_system/internal/store"
	"casino_loyalty_reward_system/internal/types"
	utils "casino_loyalty_reward_system/internal/util"

	"github.com/google/uuid"
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
