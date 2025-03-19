//go:build integration

package postgresdb_test

import (
	"context"
	"testing"
	"time"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/store/postgresdb"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetUserBy(t *testing.T) {
	defer truncate()

	log, err := zap.NewDevelopment()
	require.NoError(t, err)

	var (
		ctx = context.Background()

		databaseManager = postgresdb.New(log.Sugar(), testDB)
	)

	ID, err := uuid.Parse("c94e17df-ce34-4196-a8ca-5497e478d95d")
	require.NoError(t, err)

	ID2, err := uuid.Parse("7d469fa9-187e-4b43-9e16-b20e1939e84a")
	require.NoError(t, err)

	tm := time.Now()

	usr1 := types.User{
		ID:       ID,
		Name:     "John1",
		Email:    "john1@example.com",
		Password: "password",
		Role:     1,
		Balance:  0,
		Created:  tm,
		Updated:  tm,
	}
	usr2 := types.User{
		ID:       ID2,
		Name:     "John2",
		Email:    "john2@example.com",
		Password: "password",
		Role:     1,
		Balance:  0,
		Created:  tm,
		Updated:  tm,
	}

	_, err = testDB.Exec(ctx, `
		INSERT INTO users (
		id,
		name,
		email,
		password,
		role,
		created,
		updated
	) VALUES ($1, $2, $3, $4, $5, $6, $7), ($8, $9, $10, $11, $12, $13, $14)`,
		usr1.ID, usr1.Name, usr1.Email, usr1.Password, usr1.Role, usr1.Created, usr1.Updated,
		usr2.ID, usr2.Name, usr2.Email, usr2.Password, usr2.Role, usr2.Created, usr2.Updated,
	)
	require.NoError(t, err)

	tests := []struct {
		name           string
		filter         types.UserFilter
		expectedResult types.User
		expectedError  error
	}{
		{
			name:           "get user by id",
			filter:         types.UserFilter{ByID: uuid.NullUUID{UUID: usr1.ID, Valid: true}},
			expectedResult: usr1,
			expectedError:  nil,
		},
		{
			name:           "get other user by id",
			filter:         types.UserFilter{ByID: uuid.NullUUID{UUID: usr2.ID, Valid: true}},
			expectedResult: usr2,
			expectedError:  nil,
		},
		{
			name:           "get account by email",
			filter:         types.UserFilter{ByEmail: &usr2.Email},
			expectedResult: usr2,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := databaseManager.UserGetBy(context.Background(), tt.filter)

			require.ErrorIs(t, err, tt.expectedError)
			if tt.expectedError == nil {
				require.EqualValues(t, user.ID, tt.expectedResult.ID)
				require.EqualValues(t, user.Name, tt.expectedResult.Name)
				require.EqualValues(t, user.Email, tt.expectedResult.Email)
				require.EqualValues(t, user.Role, tt.expectedResult.Role)
			}
		})
	}
}
