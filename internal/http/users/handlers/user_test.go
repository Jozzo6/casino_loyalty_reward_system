package handlers_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/fakes"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/http/users/handlers"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/test"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	type fields struct {
		userProvider *fakes.FakeUserProvider
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	tests := []struct {
		name           string
		fields         fields
		req            test.TestRequest
		expectedCode   int
		expectedOutput string
	}{
		{
			name: "it should register the user",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
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
			},
			req: test.TestRequest{
				Body: `{"name":"John","password":"password","email":"john@example.com"}`,
			},
			expectedCode:   http.StatusOK,
			expectedOutput: `{"id":"460aec7e-7d58-42fd-93b8-bca05a77bbf5","name":"John","email":"john@example.com","token":"token"}`,
		},
		{
			name: "it should fail register because of email validation",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					RegisterStub: func(ctx context.Context, u types.User) (types.User, string, error) {
						return types.User{}, "", nil
					},
				},
			},
			req: test.TestRequest{
				Body: `{"name":"John","password":"password","email":"john@example"}`,
			},
			expectedCode:   http.StatusBadRequest,
			expectedOutput: `Key: 'RegisterRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag`,
		},
		{
			name: "it should fail register because of password validation",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					RegisterStub: func(ctx context.Context, u types.User) (types.User, string, error) {
						return types.User{}, "", nil
					},
				},
			},
			req: test.TestRequest{
				Body: `{"name":"John","password":"pa","email":"john@example.com"}`,
			},
			expectedCode:   http.StatusBadRequest,
			expectedOutput: `Key: 'RegisterRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag`,
		},
		{
			name: "it should fail register because of name validation",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					RegisterStub: func(ctx context.Context, u types.User) (types.User, string, error) {
						return types.User{}, "", nil
					},
				},
			},
			req: test.TestRequest{
				Body: `{"name":"","password":"password","email":"john@example.com"}`,
			},
			expectedCode:   http.StatusBadRequest,
			expectedOutput: `Key: 'RegisterRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := handlers.NewAccountsRouter(tt.fields.userProvider)
			w := httptest.NewRecorder()
			r, err := tt.req.GetRequest(http.MethodPost)
			require.NoError(t, err)
			router.Register().ServeHTTP(w, r)

			resp := w.Result()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedCode, resp.StatusCode)
			require.Regexp(t, regexp.MustCompile(tt.expectedOutput), string(respBody))
		})
	}
}

func TestLogin(t *testing.T) {
	type fields struct {
		userProvider *fakes.FakeUserProvider
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	tests := []struct {
		name           string
		fields         fields
		req            test.TestRequest
		expectedCode   int
		expectedOutput string
	}{
		{
			name: "it should login the user",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					LoginStub: func(ctx context.Context, u types.User) (types.User, string, error) {
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
			},
			req: test.TestRequest{
				Body: `{"password":"password","email":"john@example.com"}`,
			},
			expectedCode:   http.StatusOK,
			expectedOutput: `{"id":"460aec7e-7d58-42fd-93b8-bca05a77bbf5","name":"John","email":"john@example.com","token":"token"}`,
		},
		{
			name: "it should fail login because of email validation",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					LoginStub: func(ctx context.Context, u types.User) (types.User, string, error) {
						return types.User{}, "", nil
					},
				},
			},
			req: test.TestRequest{
				Body: `{"password":"password","email":"john@example"}`,
			},
			expectedCode:   http.StatusBadRequest,
			expectedOutput: `Key: 'LoginRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag`,
		},
		{
			name: "it should fail login because of password validation",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					LoginStub: func(ctx context.Context, u types.User) (types.User, string, error) {
						return types.User{}, "", nil
					},
				},
			},
			req: test.TestRequest{
				Body: `{"password":"pa","email":"john@example.com"}`,
			},
			expectedCode:   http.StatusBadRequest,
			expectedOutput: `Key: 'LoginRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := handlers.NewAccountsRouter(tt.fields.userProvider)
			w := httptest.NewRecorder()
			r, err := tt.req.GetRequest(http.MethodPost)
			require.NoError(t, err)
			router.Login().ServeHTTP(w, r)

			resp := w.Result()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedCode, resp.StatusCode)
			require.Regexp(t, regexp.MustCompile(tt.expectedOutput), string(respBody))
		})
	}
}

func TestGetUser(t *testing.T) {
	type fields struct {
		userProvider *fakes.FakeUserProvider
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	tests := []struct {
		name           string
		fields         fields
		req            test.TestRequest
		expectedCode   int
		expectedOutput string
	}{
		{
			name: "it should return the user",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					GetUserStub: func(ctx context.Context, u uuid.UUID) (types.User, error) {
						return types.User{
							ID:      ID,
							Name:    "John",
							Email:   "john@example.com",
							Role:    1,
							Balance: 0,
						}, nil
					},
				},
			},
			req: test.TestRequest{
				Vars: map[string]string{
					"id": "460aec7e-7d58-42fd-93b8-bca05a77bbf5",
				},
			},
			expectedCode:   http.StatusOK,
			expectedOutput: `{"id":"460aec7e-7d58-42fd-93b8-bca05a77bbf5","name":"John","email":"john@example.com","role":1,"balance":0,"created":"0001-01-01T00:00:00Z","updated":"0001-01-01T00:00:00Z","Password":""}`,
		},
		{
			name: "it should invalid uuid format",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					GetUserStub: func(ctx context.Context, u uuid.UUID) (types.User, error) {
						return types.User{}, nil
					},
				},
			},
			req: test.TestRequest{
				Vars: map[string]string{
					"id": "460aec7e-7d58-42fd-93b8abca05a77bbf5",
				},
			},
			expectedCode:   http.StatusBadRequest,
			expectedOutput: `invalid UUID format`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := handlers.NewAccountsRouter(tt.fields.userProvider)
			w := httptest.NewRecorder()
			r, err := tt.req.GetRequest(http.MethodPost)
			require.NoError(t, err)
			router.GetUser().ServeHTTP(w, r)

			resp := w.Result()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedCode, resp.StatusCode)
			require.Regexp(t, regexp.MustCompile(tt.expectedOutput), string(respBody))
		})
	}
}

func TestGetUsers(t *testing.T) {
	type fields struct {
		userProvider *fakes.FakeUserProvider
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	tests := []struct {
		name           string
		fields         fields
		req            test.TestRequest
		expectedCode   int
		expectedOutput string
	}{
		{
			name: "it should return the 2 users",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					GetUsersStub: func(ctx context.Context) ([]types.User, error) {
						return []types.User{
							{
								ID:      ID,
								Name:    "John",
								Email:   "john@example.com",
								Role:    1,
								Balance: 0,
							},
						}, nil
					},
				},
			},
			req:            test.TestRequest{},
			expectedCode:   http.StatusOK,
			expectedOutput: ``,
		},
		{
			name: "it should return the empty field",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					GetUsersStub: func(ctx context.Context) ([]types.User, error) {
						return []types.User{}, nil
					},
				},
			},
			req:            test.TestRequest{},
			expectedCode:   http.StatusOK,
			expectedOutput: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := handlers.NewAccountsRouter(tt.fields.userProvider)
			w := httptest.NewRecorder()
			r, err := tt.req.GetRequest(http.MethodPost)
			require.NoError(t, err)
			router.GetUsers().ServeHTTP(w, r)

			resp := w.Result()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedCode, resp.StatusCode)
			require.Regexp(t, regexp.MustCompile(tt.expectedOutput), string(respBody))
		})
	}
}

func TestUpdateUser(t *testing.T) {
	type fields struct {
		userProvider *fakes.FakeUserProvider
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	tests := []struct {
		name           string
		fields         fields
		req            test.TestRequest
		expectedCode   int
		expectedOutput string
	}{
		{
			name: "it should update the user",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					UpdateUserStub: func(ctx context.Context, user types.User) (types.User, error) {
						return types.User{
							ID:      ID,
							Name:    "John",
							Email:   "john@example.com",
							Role:    2,
							Balance: 0,
						}, nil
					},
				},
			},
			req: test.TestRequest{
				Body: `{"id":"460aec7e-7d58-42fd-93b8-bca05a77bbf5","name":"John","email":"john@example.com","role":2,"balance": 0}`,
			},
			expectedCode:   http.StatusOK,
			expectedOutput: ``,
		},
		{
			name: "it should fail to update the user",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					UpdateUserStub: func(ctx context.Context, user types.User) (types.User, error) {
						return types.User{}, nil
					},
				},
			},
			expectedCode:   http.StatusBadRequest,
			expectedOutput: `EOF`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := handlers.NewAccountsRouter(tt.fields.userProvider)
			w := httptest.NewRecorder()
			r, err := tt.req.GetRequest(http.MethodPost)
			require.NoError(t, err)
			router.UpdateUser().ServeHTTP(w, r)

			resp := w.Result()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedCode, resp.StatusCode)
			require.Regexp(t, regexp.MustCompile(tt.expectedOutput), string(respBody))
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type fields struct {
		userProvider *fakes.FakeUserProvider
	}

	tests := []struct {
		name           string
		fields         fields
		req            test.TestRequest
		expectedCode   int
		expectedOutput string
	}{
		{
			name: "it should delete the user",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					DeleteUserStub: func(ctx context.Context, id uuid.UUID) error {
						return nil
					},
				},
			},
			req: test.TestRequest{
				Vars: map[string]string{
					"id": "460aec7e-7d58-42fd-93b8-bca05a77bbf5",
				},
			},
			expectedCode:   http.StatusOK,
			expectedOutput: `"OK"`,
		},
		{
			name: "it should failed to delete ",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					DeleteUserStub: func(ctx context.Context, id uuid.UUID) error {
						return pgx.ErrNoRows
					},
				},
			},
			req: test.TestRequest{
				Vars: map[string]string{
					"id": "460aec7e-7d58-42fd-93b8-bca05a77bbf5",
				},
			},
			expectedCode:   http.StatusNotFound,
			expectedOutput: `"user with 460aec7e-7d58-42fd-93b8-bca05a77bbf5 id was not found"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := handlers.NewAccountsRouter(tt.fields.userProvider)
			w := httptest.NewRecorder()
			r, err := tt.req.GetRequest(http.MethodPost)
			require.NoError(t, err)
			router.DeleteUser().ServeHTTP(w, r)

			resp := w.Result()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedCode, resp.StatusCode)
			require.Regexp(t, regexp.MustCompile(tt.expectedOutput), string(respBody))
		})
	}
}

func TestUpdateUserBalance(t *testing.T) {
	type fields struct {
		userProvider *fakes.FakeUserProvider
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	tests := []struct {
		name           string
		fields         fields
		req            test.TestRequest
		expectedCode   int
		expectedOutput string
	}{
		{
			name: "it should update the user balance add",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					UpdateUserBalanceStub: func(ctx context.Context, user types.User, value float64, transactionType types.TransactionType) (types.User, error) {
						return types.User{
							ID:      ID,
							Name:    "John",
							Email:   "john@example.com",
							Role:    2,
							Balance: 10,
						}, nil
					},
				},
			},
			req: test.TestRequest{
				Context: context.WithValue(context.Background(), types.CtxKeyAccount, types.User{
					ID:      ID,
					Name:    "John",
					Email:   "john@example.com",
					Role:    2,
					Balance: 0,
				}),
				Body: `{"value": 10,"transaction_type":"add"}`,
			},
			expectedCode:   http.StatusOK,
			expectedOutput: ``,
		},
		{
			name: "it should update the user balance remove",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					UpdateUserBalanceStub: func(ctx context.Context, user types.User, value float64, transactionType types.TransactionType) (types.User, error) {
						return types.User{
							ID:      ID,
							Name:    "John",
							Email:   "john@example.com",
							Role:    2,
							Balance: 90,
						}, nil
					},
				},
			},
			req: test.TestRequest{
				Context: context.WithValue(context.Background(), types.CtxKeyAccount, types.User{
					ID:      ID,
					Name:    "John",
					Email:   "john@example.com",
					Role:    2,
					Balance: 100,
				}),
				Body: `{"value": 10,"transaction_type":"remove"}`,
			},
			expectedCode:   http.StatusOK,
			expectedOutput: `{"id":"460aec7e-7d58-42fd-93b8-bca05a77bbf5","name":"John","email":"john@example.com","role":2,"balance":90,"created":"0001-01-01T00:00:00Z","updated":"0001-01-01T00:00:00Z","Password":""}`,
		},
		{
			name: "it should fail update the user balance",
			fields: fields{
				userProvider: &fakes.FakeUserProvider{
					UpdateUserBalanceStub: func(ctx context.Context, user types.User, value float64, transactionType types.TransactionType) (types.User, error) {
						return types.User{}, types.ErrInsufficientBalance
					},
				},
			},
			req: test.TestRequest{
				Context: context.WithValue(context.Background(), types.CtxKeyAccount, types.User{
					ID:      ID,
					Name:    "John",
					Email:   "john@example.com",
					Role:    2,
					Balance: 0,
				}),
				Body: `{"value": 10,"transaction_type":"remove"}`,
			},
			expectedCode:   http.StatusBadRequest,
			expectedOutput: `"Insufficient balance"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := handlers.NewAccountsRouter(tt.fields.userProvider)
			w := httptest.NewRecorder()
			r, err := tt.req.GetRequest(http.MethodPost)
			require.NoError(t, err)
			router.UpdateBalance().ServeHTTP(w, r)

			resp := w.Result()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedCode, resp.StatusCode)
			require.Regexp(t, regexp.MustCompile(tt.expectedOutput), string(respBody))
		})
	}
}
