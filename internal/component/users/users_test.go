package users_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/component/users"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/fakes"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/store"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

const (
	jwtKey      = "test_key"
	jwtDuration = time.Duration(time.Hour)
)

type fields struct {
	persistentStore store.Persistent
	pubsub          store.PubSub
	tester          users.UserProvider
}

func TestGetUser(t *testing.T) {
	type args struct {
		userID uuid.UUID
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	tests := []struct {
		name           string
		fields         fields
		vars           map[string]string
		args           args
		expectedError  error
		expectedOutput types.User
	}{
		{
			name: "it should get user",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					UserGetByStub: func(ctx context.Context, uf types.UserFilter) (types.User, error) {
						return types.User{
							ID:       ID,
							Name:     "John",
							Email:    "john@example.com",
							Password: "password",
							Role:     1,
							Balance:  0,
						}, nil
					},
				},
				tester: &fakes.FakeUserProvider{},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{userID: ID},
			expectedOutput: types.User{
				ID:       ID,
				Name:     "John",
				Email:    "john@example.com",
				Password: "password",
				Role:     1,
				Balance:  0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := users.New(tt.fields.persistentStore, tt.fields.pubsub, []byte(jwtKey), jwtDuration)
			res, err := c.GetUser(context.Background(), tt.args.userID)

			require.ErrorIs(t, err, tt.expectedError)
			if tt.expectedError == nil {
				require.Equal(t, tt.expectedOutput, res)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	type args struct {
		user types.User
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	conflictErr := errors.New(`ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)`)

	tests := []struct {
		name           string
		fields         fields
		vars           map[string]string
		args           args
		expectedError  error
		expectedOutput types.User
	}{
		{
			name: "it should get register",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					UserCreateStub: func(ctx context.Context, u types.User) (types.User, error) {
						return types.User{
							ID:       ID,
							Name:     "John",
							Email:    "john@example.com",
							Password: "password",
							Role:     1,
							Balance:  0,
						}, nil
					},
				},
				tester: &fakes.FakeUserProvider{
					RegisterStub: func(ctx context.Context, u types.User) (types.User, string, error) {
						return types.User{
							ID:       ID,
							Name:     "John",
							Email:    "john@example.com",
							Password: "password",
							Role:     1,
							Balance:  0,
						}, "token", nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{user: types.User{
				Name:     "John",
				Email:    "john@example.com",
				Password: "password",
			}},
			expectedOutput: types.User{
				ID:       ID,
				Name:     "John",
				Email:    "john@example.com",
				Password: "password",
				Role:     1,
				Balance:  0,
			},
		},
		{
			name: "it should fail email in use",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					UserCreateStub: func(ctx context.Context, u types.User) (types.User, error) {
						return types.User{}, conflictErr
					},
				},
				tester: &fakes.FakeUserProvider{},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{user: types.User{
				Name:     "John",
				Email:    "john@example.com",
				Password: "password",
			}},
			expectedOutput: types.User{},
			expectedError:  conflictErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := users.New(tt.fields.persistentStore, tt.fields.pubsub, []byte(jwtKey), jwtDuration)
			res, token, err := c.Register(context.Background(), tt.args.user)

			require.ErrorIs(t, err, tt.expectedError)
			if tt.expectedError == nil {
				require.NotEmpty(t, token)
				require.Equal(t, tt.expectedOutput, res)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	type args struct {
		user types.User
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	user := types.User{
		ID:      ID,
		Name:    "John",
		Email:   "john@example.com",
		Role:    1,
		Balance: 0,
	}

	tests := []struct {
		name           string
		fields         fields
		vars           map[string]string
		args           args
		expectedError  error
		expectedOutput types.User
	}{
		{
			name: "it should login",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					UserGetByStub: func(ctx context.Context, uf types.UserFilter) (types.User, error) {
						return types.User{
							ID:       ID,
							Name:     "John",
							Email:    "john@example.com",
							Password: "$2a$10$slqGr93DMCar8kc6BkCY0.EeZ3/a70D7bq1/gD25pcSw2k0c9d2gW",
							Role:     1,
							Balance:  0,
						}, nil
					},
				},
				tester: &fakes.FakeUserProvider{
					LoginStub: func(ctx context.Context, u types.User) (types.User, string, error) {
						return user, "token", nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				user: types.User{
					Email:    "john@example.com",
					Password: "password",
				}},
			expectedOutput: user,
		},
		{
			name: "it should fail user not found",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					UserGetByStub: func(ctx context.Context, uf types.UserFilter) (types.User, error) {
						return types.User{}, pgx.ErrNoRows
					},
				},
				tester: &fakes.FakeUserProvider{},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{user: types.User{
				Name:     "John",
				Email:    "john@example.com",
				Password: "password",
			}},
			expectedOutput: types.User{},
			expectedError:  types.ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := users.New(tt.fields.persistentStore, tt.fields.pubsub, []byte(jwtKey), jwtDuration)
			res, token, err := c.Login(context.Background(), tt.args.user)

			require.ErrorIs(t, err, tt.expectedError)
			if tt.expectedError == nil {
				require.NotEmpty(t, token)
				require.Equal(t, tt.expectedOutput, res)
			}
		})
	}
}

func TestAuth(t *testing.T) {

}

func TestGetUsers(t *testing.T) {
	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	tests := []struct {
		name           string
		fields         fields
		vars           map[string]string
		expectedError  error
		expectedOutput []types.User
	}{
		{
			name: "it should be empty",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					GetUsersStub: func(ctx context.Context) ([]types.User, error) {
						return []types.User{}, nil
					},
				},
				tester: &fakes.FakeUserProvider{
					GetUsersStub: func(ctx context.Context) ([]types.User, error) {
						return []types.User{}, nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			expectedOutput: []types.User{},
		},
		{
			name: "it should get users",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					GetUsersStub: func(ctx context.Context) ([]types.User, error) {
						return []types.User{
							{
								ID:       ID,
								Name:     "John",
								Email:    "john@example.com",
								Password: "password",
								Role:     1,
								Balance:  0,
							},
						}, nil
					},
				},
				tester: &fakes.FakeUserProvider{
					GetUsersStub: func(ctx context.Context) ([]types.User, error) {
						return []types.User{
							{
								ID:       ID,
								Name:     "John",
								Email:    "john@example.com",
								Password: "password",
								Role:     1,
								Balance:  0,
							},
						}, nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			expectedOutput: []types.User{
				{
					ID:       ID,
					Name:     "John",
					Email:    "john@example.com",
					Password: "password",
					Role:     1,
					Balance:  0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := users.New(tt.fields.persistentStore, tt.fields.pubsub, []byte(jwtKey), jwtDuration)
			res, err := c.GetUsers(context.Background())

			require.ErrorIs(t, err, tt.expectedError)
			if tt.expectedError == nil {
				require.Equal(t, tt.expectedOutput, res)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	type args struct {
		user types.User
	}

	tests := []struct {
		name           string
		fields         fields
		vars           map[string]string
		args           args
		expectedError  error
		expectedOutput types.User
	}{
		{
			name: "it should update",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					UserUpdateStub: func(ctx context.Context, u types.User) (types.User, error) {
						return types.User{
							ID:       ID,
							Name:     "John",
							Email:    "john@example.com",
							Password: "password",
							Role:     0,
							Balance:  0,
						}, nil
					},
				},
				tester: &fakes.FakeUserProvider{
					UpdateUserStub: func(ctx context.Context, u types.User) (types.User, error) {
						return types.User{
							ID:       ID,
							Name:     "John",
							Email:    "john@example.com",
							Password: "password",
							Role:     0,
							Balance:  0,
						}, nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				user: types.User{
					ID:       ID,
					Name:     "John",
					Email:    "john@example.com",
					Password: "password",
					Role:     0,
					Balance:  0,
				},
			},
			expectedOutput: types.User{
				ID:       ID,
				Name:     "John",
				Email:    "john@example.com",
				Password: "password",
				Role:     0,
				Balance:  0,
			},
		},
		{
			name: "it should fail not found",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					UserUpdateStub: func(ctx context.Context, u types.User) (types.User, error) {
						return types.User{}, pgx.ErrNoRows
					},
				},
				tester: &fakes.FakeUserProvider{
					UpdateUserStub: func(ctx context.Context, u types.User) (types.User, error) {
						return types.User{}, pgx.ErrNoRows
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				user: types.User{
					ID:       ID,
					Name:     "John",
					Email:    "john@example.com",
					Password: "password",
					Role:     1,
					Balance:  0,
				},
			},
			expectedOutput: types.User{},
			expectedError:  pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := users.New(tt.fields.persistentStore, tt.fields.pubsub, []byte(jwtKey), jwtDuration)
			res, err := c.UpdateUser(context.Background(), tt.args.user)

			require.ErrorIs(t, err, tt.expectedError)
			if tt.expectedError == nil {
				require.Equal(t, tt.expectedOutput, res)
			}
		})
	}
}

func TestUpdateUserBalance(t *testing.T) {
	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	type args struct {
		user        types.User
		transaction types.TransactionType
		value       float64
	}

	tests := []struct {
		name           string
		fields         fields
		vars           map[string]string
		args           args
		expectedError  error
		expectedOutput types.User
	}{
		{
			name: "it should update balance add",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					UserBalanceUpdateStub: func(ctx context.Context, u uuid.UUID, f float64) (types.User, error) {
						return types.User{
							ID:       ID,
							Name:     "John",
							Email:    "john@example.com",
							Password: "password",
							Role:     0,
							Balance:  10,
						}, nil
					},
				},
				tester: &fakes.FakeUserProvider{
					UpdateUserBalanceStub: func(ctx context.Context, u types.User, f float64, tt types.TransactionType) (types.User, error) {
						return types.User{
							ID:       ID,
							Name:     "John",
							Email:    "john@example.com",
							Password: "password",
							Role:     0,
							Balance:  10,
						}, nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				user: types.User{
					ID:       ID,
					Name:     "John",
					Email:    "john@example.com",
					Password: "password",
					Role:     0,
					Balance:  0,
				},
				transaction: types.TransactionTypeAdd,
				value:       10,
			},
			expectedOutput: types.User{
				ID:       ID,
				Name:     "John",
				Email:    "john@example.com",
				Password: "password",
				Role:     0,
				Balance:  10,
			},
		},
		{
			name: "it should update balance remove",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					UserBalanceUpdateStub: func(ctx context.Context, u uuid.UUID, f float64) (types.User, error) {
						return types.User{
							ID:       ID,
							Name:     "John",
							Email:    "john@example.com",
							Password: "password",
							Role:     0,
							Balance:  10,
						}, nil
					},
				},
				tester: &fakes.FakeUserProvider{
					UpdateUserBalanceStub: func(ctx context.Context, u types.User, f float64, tt types.TransactionType) (types.User, error) {
						return types.User{
							ID:       ID,
							Name:     "John",
							Email:    "john@example.com",
							Password: "password",
							Role:     0,
							Balance:  10,
						}, nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				user: types.User{
					ID:       ID,
					Name:     "John",
					Email:    "john@example.com",
					Password: "password",
					Role:     0,
					Balance:  20,
				},
				transaction: types.TransactionTypeRemove,
				value:       10,
			},
			expectedOutput: types.User{
				ID:       ID,
				Name:     "John",
				Email:    "john@example.com",
				Password: "password",
				Role:     0,
				Balance:  10,
			},
		},
		{
			name: "it should fail update balance remove insufficient funds",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					UserBalanceUpdateStub: func(ctx context.Context, u uuid.UUID, f float64) (types.User, error) {
						return types.User{}, nil
					},
				},
				tester: &fakes.FakeUserProvider{
					UpdateUserBalanceStub: func(ctx context.Context, u types.User, f float64, tt types.TransactionType) (types.User, error) {
						return types.User{}, types.ErrInsufficientBalance
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				user: types.User{
					ID:       ID,
					Name:     "John",
					Email:    "john@example.com",
					Password: "password",
					Role:     0,
					Balance:  5,
				},
				transaction: types.TransactionTypeRemove,
				value:       10,
			},
			expectedOutput: types.User{},
			expectedError:  types.ErrInsufficientBalance,
		},
		{
			name: "it should fail not found",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					UserBalanceUpdateStub: func(ctx context.Context, u uuid.UUID, f float64) (types.User, error) {
						return types.User{}, pgx.ErrNoRows
					},
				},
				tester: &fakes.FakeUserProvider{
					UpdateUserBalanceStub: func(ctx context.Context, u types.User, f float64, tt types.TransactionType) (types.User, error) {
						return types.User{}, pgx.ErrNoRows
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				user: types.User{
					ID:       ID,
					Name:     "John",
					Email:    "john@example.com",
					Password: "password",
					Role:     1,
					Balance:  0,
				},
				transaction: types.TransactionTypeAdd,
				value:       10,
			},
			expectedOutput: types.User{},
			expectedError:  pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := users.New(tt.fields.persistentStore, tt.fields.pubsub, []byte(jwtKey), jwtDuration)
			res, err := c.UpdateUserBalance(context.Background(), tt.args.user, tt.args.value, tt.args.transaction)

			require.ErrorIs(t, err, tt.expectedError)
			if tt.expectedError == nil {
				require.Equal(t, tt.expectedOutput, res)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	type args struct {
		userID uuid.UUID
	}

	tests := []struct {
		name          string
		fields        fields
		vars          map[string]string
		args          args
		expectedError error
	}{
		{
			name: "it should delete user",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					UserBalanceUpdateStub: func(ctx context.Context, u uuid.UUID, f float64) (types.User, error) {
						return types.User{
							ID:       ID,
							Name:     "John",
							Email:    "john@example.com",
							Password: "password",
							Role:     0,
							Balance:  10,
						}, nil
					},
				},
				tester: &fakes.FakeUserProvider{
					UpdateUserBalanceStub: func(ctx context.Context, u types.User, f float64, tt types.TransactionType) (types.User, error) {
						return types.User{
							ID:       ID,
							Name:     "John",
							Email:    "john@example.com",
							Password: "password",
							Role:     0,
							Balance:  10,
						}, nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				userID: ID,
			},
		},
		{
			name: "it should fail not found",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					UserDeleteStub: func(ctx context.Context, u uuid.UUID) error {
						return pgx.ErrNoRows
					},
				},
				tester: &fakes.FakeUserProvider{
					DeleteUserStub: func(ctx context.Context, u uuid.UUID) error {
						return pgx.ErrNoRows
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				userID: ID,
			},
			expectedError: pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := users.New(tt.fields.persistentStore, tt.fields.pubsub, []byte(jwtKey), jwtDuration)
			err := c.DeleteUser(context.Background(), tt.args.userID)

			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}
