package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/component/users"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/store"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"
	utils "github.com/Jozzo6/casino_loyalty_reward_system/internal/util"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type usersRouter struct {
	component users.UserProvider
}

func NewAccountsRouter(component users.UserProvider) *usersRouter {
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

// Register handles user registration.
// @Summary Register a new user
// @Description Creates a new user account and returns the user details along with a token.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration details"
// @Success 200 {object} RegisterResponse "User registered successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid request payload"
// @Failure 409 {object} types.ErrorResponse "User already exists"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/register [post]
func (ur *usersRouter) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest

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

		createdUser, token, err := ur.component.Register(r.Context(), types.User{
			Email:    req.Email,
			Name:     req.Name,
			Password: req.Password,
		})

		if store.IsErrConflict(err) {
			utils.WriteError(log, w, http.StatusConflict, err)
			return
		} else if err != nil {
			utils.WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(log, w, http.StatusOK, RegisterResponse{
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

// Login handles user login.
// @Summary Login a user
// @Description Authenticates a user and returns their details along with a token.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login details"
// @Success 200 {object} LoginResponse "User logged in successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid request payload"
// @Failure 401 {object} types.ErrorResponse "Unauthorized"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/login [post]
func (ur *usersRouter) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest

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

		user, token, err := ur.component.Login(r.Context(), types.User{
			Email:    req.Email,
			Password: req.Password,
		})
		if errors.Is(err, types.ErrUnauthorized) {
			utils.WriteError(log, w, http.StatusUnauthorized, err)
			return
		}
		if err != nil {
			utils.WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(log, w, http.StatusOK, RegisterResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: token,
		})
	}
}

// GetUser retrieves a user by ID.
// @Summary Get a user by ID
// @Description Retrieves the details of a user by their unique ID.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} types.User "User details retrieved successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid user ID"
// @Failure 404 {object} types.ErrorResponse "User not found"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id} [get]
func (ur *usersRouter) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			log.Errorf("failed to get user id: %s", err)
			utils.WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		user, err := ur.component.GetUser(r.Context(), id)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Errorf("user with id: %s was not found: %s", id.String(), err)
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

// GetUsers retrieves all users.
// @Summary Get all users
// @Description Retrieves a list of all users.
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} types.User "List of users retrieved successfully"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/users [get]
func (ur *usersRouter) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		users, err := ur.component.GetUsers(r.Context())
		if err != nil {
			utils.WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(log, w, http.StatusOK, users)
	}
}

// UpdateUser updates a user's details.
// @Summary Update a user
// @Description Updates the details of an existing user.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body types.User true "User details to update"
// @Success 200 {object} types.User "User updated successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid request payload"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id} [put]
func (ur *usersRouter) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.User

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

		user, err := ur.component.UpdateUser(r.Context(), req)
		if err != nil {
			utils.WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(log, w, http.StatusOK, user)

	}
}

// DeleteUser deletes a user by ID.
// @Summary Delete a user
// @Description Deletes a user by their unique ID.
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {string} string "User deleted successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid user ID"
// @Failure 404 {object} types.ErrorResponse "User not found"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id} [delete]
func (ur *usersRouter) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := types.GetLoggerFromContext(r.Context())

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			log.Errorf("failed to get user id: %s", err)
			utils.WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		err = ur.component.DeleteUser(r.Context(), id)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Errorf("user with id: %s was not found to be updated: %s", id.String(), err)
			utils.WriteError(log, w, http.StatusNotFound, fmt.Errorf("user with %s id was not found", id.String()))
			return
		}
		if err != nil {
			utils.WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(log, w, http.StatusOK, "OK")
	}
}

type UpdateBalanceRequest struct {
	Value           float64               `json:"value" validate:"required"`
	TransactionType types.TransactionType `json:"transaction_type" validate:"required"`
}

// UpdateBalance updates a user's balance.
// @Summary Update user balance
// @Description Updates the balance of a user based on the transaction type and value.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body UpdateBalanceRequest true "Balance update details"
// @Success 200 {object} types.User "User balance updated successfully"
// @Failure 400 {object} types.ErrorResponse "Invalid request payload or insufficient balance"
// @Failure 500 {object} types.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id}/balance [put]
func (ur *usersRouter) UpdateBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateBalanceRequest

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

		us, err := types.GetAccountFromContext(r.Context())
		if err != nil {
			utils.WriteError(log, w, http.StatusBadRequest, err)
			return
		}

		user, err := ur.component.UpdateUserBalance(r.Context(), us, req.Value, req.TransactionType)
		if err != nil {
			if errors.Is(err, types.ErrInsufficientBalance) {
				utils.WriteError(log, w, http.StatusBadRequest, err)
				return
			}
			utils.WriteError(log, w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(log, w, http.StatusOK, user)

	}
}
