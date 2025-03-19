package handlers_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/fakes"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/http/promotions/handlers"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/test"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestCreatePromotion(t *testing.T) {
	type fields struct {
		promotionsProvider *fakes.FakePromotionProvider
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
			name: "it should create promotion",
			fields: fields{
				promotionsProvider: &fakes.FakePromotionProvider{
					CreatePromotionsStub: func(ctx context.Context, p types.Promotion) (types.Promotion, error) {
						return types.Promotion{
							ID:          ID,
							Title:       "Title",
							Description: "Description",
							Type:        types.Regular,
							IsActive:    true,
							Amount:      10,
						}, nil
					},
				},
			},
			req: test.TestRequest{
				Body: `{"title":"Title","description":"Description","type":"regular","is_active":true,"amount":10}`,
			},
			expectedCode:   http.StatusOK,
			expectedOutput: `{"id":"460aec7e-7d58-42fd-93b8-bca05a77bbf5","title":"Title","description":"Description","amount":10,"is_active":true,"type":"regular","created":"0001-01-01T00:00:00Z","updated":"0001-01-01T00:00:00Z"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := handlers.NewPromotionsRouter(tt.fields.promotionsProvider)
			w := httptest.NewRecorder()
			r, err := tt.req.GetRequest(http.MethodPost)
			require.NoError(t, err)
			router.CreatePromotion().ServeHTTP(w, r)

			resp := w.Result()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedCode, resp.StatusCode)
			require.Regexp(t, regexp.MustCompile(tt.expectedOutput), string(respBody))
		})
	}
}

func TestGetPromotionByID(t *testing.T) {
	type fields struct {
		promotionsProvider *fakes.FakePromotionProvider
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
			name: "it should get promotion by id",
			fields: fields{
				promotionsProvider: &fakes.FakePromotionProvider{
					GetPromotionByIDStub: func(ctx context.Context, u uuid.UUID) (types.Promotion, error) {
						return types.Promotion{
							ID:          ID,
							Title:       "Title",
							Description: "Description",
							Type:        types.Regular,
							IsActive:    true,
							Amount:      10,
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
			expectedOutput: `{"id":"460aec7e-7d58-42fd-93b8-bca05a77bbf5","title":"Title","description":"Description","amount":10,"is_active":true,"type":"regular","created":"0001-01-01T00:00:00Z","updated":"0001-01-01T00:00:00Z"}`,
		},
		{
			name: "it should fail to get promotion by id not found",
			fields: fields{
				promotionsProvider: &fakes.FakePromotionProvider{
					GetPromotionByIDStub: func(ctx context.Context, u uuid.UUID) (types.Promotion, error) {
						return types.Promotion{}, pgx.ErrNoRows
					},
				},
			},
			req: test.TestRequest{
				Vars: map[string]string{
					"id": "460aec7e-7d58-42fd-93b8-bca05a77bbf5",
				},
			},
			expectedCode:   http.StatusNotFound,
			expectedOutput: `{"message":"no rows in result set"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := handlers.NewPromotionsRouter(tt.fields.promotionsProvider)
			w := httptest.NewRecorder()
			r, err := tt.req.GetRequest(http.MethodPost)
			require.NoError(t, err)
			router.GetPromotionByID().ServeHTTP(w, r)

			resp := w.Result()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedCode, resp.StatusCode)
			require.Regexp(t, regexp.MustCompile(tt.expectedOutput), string(respBody))
		})
	}
}

func TestGetPromotions(t *testing.T) {
	type fields struct {
		promotionsProvider *fakes.FakePromotionProvider
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
			name: "it should get promotion by id",
			fields: fields{
				promotionsProvider: &fakes.FakePromotionProvider{
					GetPromotionsStub: func(ctx context.Context) ([]types.Promotion, error) {
						return []types.Promotion{
							{
								ID:          ID,
								Title:       "Title",
								Description: "Description",
								Type:        types.Regular,
								IsActive:    true,
								Amount:      10,
							},
						}, nil
					},
				},
			},
			req:            test.TestRequest{},
			expectedCode:   http.StatusOK,
			expectedOutput: `[{"id":"460aec7e-7d58-42fd-93b8-bca05a77bbf5","title":"Title","description":"Description","amount":10,"is_active":true,"type":"regular","created":"0001-01-01T00:00:00Z","updated":"0001-01-01T00:00:00Z"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := handlers.NewPromotionsRouter(tt.fields.promotionsProvider)
			w := httptest.NewRecorder()
			r, err := tt.req.GetRequest(http.MethodPost)
			require.NoError(t, err)
			router.GetPromotions().ServeHTTP(w, r)

			resp := w.Result()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedCode, resp.StatusCode)
			require.Regexp(t, regexp.MustCompile(regexp.QuoteMeta(tt.expectedOutput)), string(respBody))
		})
	}
}

func TestUpdatePromotion(t *testing.T) {
	type fields struct {
		promotionsProvider *fakes.FakePromotionProvider
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
			name: "it should update promotion",
			fields: fields{
				promotionsProvider: &fakes.FakePromotionProvider{
					UpdatePromotionStub: func(ctx context.Context, p types.Promotion) (types.Promotion, error) {
						return types.Promotion{
							ID:          ID,
							Title:       "Title",
							Description: "Description",
							Type:        types.Regular,
							IsActive:    true,
							Amount:      10,
						}, nil
					},
				},
			},
			req: test.TestRequest{
				Vars: map[string]string{
					"id": "460aec7e-7d58-42fd-93b8-bca05a77bbf5",
				},
				Body: `{"id":"460aec7e-7d58-42fd-93b8-bca05a77bbf5","title":"Title","description":"Description","amount":10,"is_active":true,"type":"regular","created":"0001-01-01T00:00:00Z","updated":"0001-01-01T00:00:00Z"}`,
			},
			expectedCode:   http.StatusOK,
			expectedOutput: `{"id":"460aec7e-7d58-42fd-93b8-bca05a77bbf5","title":"Title","description":"Description","amount":10,"is_active":true,"type":"regular","created":"0001-01-01T00:00:00Z","updated":"0001-01-01T00:00:00Z"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := handlers.NewPromotionsRouter(tt.fields.promotionsProvider)
			w := httptest.NewRecorder()
			r, err := tt.req.GetRequest(http.MethodPost)
			require.NoError(t, err)
			router.UpdatePromotion().ServeHTTP(w, r)

			resp := w.Result()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedCode, resp.StatusCode)
			require.Regexp(t, regexp.MustCompile(regexp.QuoteMeta(tt.expectedOutput)), string(respBody))
		})
	}
}

func TestDeletePromotion(t *testing.T) {
	type fields struct {
		promotionsProvider *fakes.FakePromotionProvider
	}

	tests := []struct {
		name           string
		fields         fields
		req            test.TestRequest
		expectedCode   int
		expectedOutput string
	}{
		{
			name: "it should update promotion",
			fields: fields{
				promotionsProvider: &fakes.FakePromotionProvider{
					DeletePromotionStub: func(ctx context.Context, u uuid.UUID) error {
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
				promotionsProvider: &fakes.FakePromotionProvider{
					DeletePromotionStub: func(ctx context.Context, id uuid.UUID) error {
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
			expectedOutput: `"promotion with 460aec7e-7d58-42fd-93b8-bca05a77bbf5 id was not found"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := handlers.NewPromotionsRouter(tt.fields.promotionsProvider)
			w := httptest.NewRecorder()
			r, err := tt.req.GetRequest(http.MethodPost)
			require.NoError(t, err)
			router.DeletePromotion().ServeHTTP(w, r)

			resp := w.Result()

			respBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tt.expectedCode, resp.StatusCode)
			require.Regexp(t, regexp.MustCompile(regexp.QuoteMeta(tt.expectedOutput)), string(respBody))
		})
	}
}
