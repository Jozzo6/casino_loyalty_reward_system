package promotions_test

import (
	"context"
	"testing"
	"time"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/component/promotions"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/fakes"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/store"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

type fields struct {
	persistentStore store.Persistent
	pubsub          store.PubSub
	tester          promotions.PromotionProvider
}

func TestCreatePromotions(t *testing.T) {
	type args struct {
		promotion types.Promotion
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	promotion := types.Promotion{
		ID:          ID,
		Title:       "Welcome",
		Description: "Description",
		IsActive:    true,
		Amount:      10,
		Type:        types.WelcomeBonus,
	}

	tests := []struct {
		name           string
		fields         fields
		vars           map[string]string
		args           args
		expectedError  error
		expectedOutput types.Promotion
	}{
		{
			name: "it should get register",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					PromotionCreateStub: func(ctx context.Context, p types.Promotion) (types.Promotion, error) {
						return promotion, nil
					},
				},
				tester: &fakes.FakePromotionProvider{
					CreatePromotionsStub: func(ctx context.Context, p types.Promotion) (types.Promotion, error) {
						return promotion, nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				promotion: types.Promotion{},
			},
			expectedOutput: promotion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := promotions.New(tt.fields.persistentStore)
			res, err := c.CreatePromotions(context.Background(), tt.args.promotion)

			require.ErrorIs(t, err, tt.expectedError)
			if tt.expectedError == nil {
				require.Equal(t, tt.expectedOutput, res)
			}
		})
	}
}

func TestGetPromotions(t *testing.T) {
	type args struct {
		promotion types.Promotion
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	promotion := types.Promotion{
		ID:          ID,
		Title:       "Welcome",
		Description: "Description",
		IsActive:    true,
		Amount:      10,
		Type:        types.WelcomeBonus,
	}

	tests := []struct {
		name           string
		fields         fields
		vars           map[string]string
		args           args
		expectedError  error
		expectedOutput []types.Promotion
	}{
		{
			name: "it should get register",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					GetPromotionsStub: func(ctx context.Context) ([]types.Promotion, error) {
						return []types.Promotion{promotion}, err
					},
				},
				tester: &fakes.FakePromotionProvider{
					GetPromotionsStub: func(ctx context.Context) ([]types.Promotion, error) {
						return []types.Promotion{promotion}, err
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args:           args{},
			expectedOutput: []types.Promotion{promotion},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := promotions.New(tt.fields.persistentStore)
			res, err := c.GetPromotions(context.Background())

			require.ErrorIs(t, err, tt.expectedError)
			if tt.expectedError == nil {
				require.Equal(t, tt.expectedOutput, res)
			}
		})
	}
}

func TestGetPromotionByID(t *testing.T) {
	type args struct {
		promotionID uuid.UUID
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	promotion := types.Promotion{
		ID:          ID,
		Title:       "Welcome",
		Description: "Description",
		IsActive:    true,
		Amount:      10,
		Type:        types.WelcomeBonus,
	}

	tests := []struct {
		name           string
		fields         fields
		vars           map[string]string
		args           args
		expectedError  error
		expectedOutput types.Promotion
	}{
		{
			name: "it should get promotion",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					PromotionGetByIDStub: func(ctx context.Context, u uuid.UUID) (types.Promotion, error) {
						return promotion, nil
					},
				},
				tester: &fakes.FakePromotionProvider{
					GetPromotionByIDStub: func(ctx context.Context, u uuid.UUID) (types.Promotion, error) {
						return promotion, nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				promotionID: ID,
			},
			expectedOutput: promotion,
		},
		{
			name: "it should fail promotion not found",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					PromotionGetByIDStub: func(ctx context.Context, u uuid.UUID) (types.Promotion, error) {
						return types.Promotion{}, pgx.ErrNoRows
					},
				},
				tester: &fakes.FakePromotionProvider{
					GetPromotionByIDStub: func(ctx context.Context, u uuid.UUID) (types.Promotion, error) {
						return types.Promotion{}, pgx.ErrNoRows
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				promotionID: ID,
			},
			expectedOutput: types.Promotion{},
			expectedError:  pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := promotions.New(tt.fields.persistentStore)
			res, err := c.GetPromotionByID(context.Background(), tt.args.promotionID)

			require.ErrorIs(t, err, tt.expectedError)
			if tt.expectedError == nil {
				require.Equal(t, tt.expectedOutput, res)
			}
		})
	}
}

func TestUpdatePromotion(t *testing.T) {
	type args struct {
		promotion types.Promotion
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	promotion := types.Promotion{
		ID:          ID,
		Title:       "Welcome",
		Description: "Description",
		IsActive:    true,
		Amount:      10,
		Type:        types.WelcomeBonus,
		Created:     time.Time{},
		Updated:     time.Time{},
	}

	tests := []struct {
		name           string
		fields         fields
		vars           map[string]string
		args           args
		expectedError  error
		expectedOutput types.Promotion
	}{
		{
			name: "it should update promotion",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					PromotionUpdateStub: func(ctx context.Context, p types.Promotion) (types.Promotion, error) {
						return promotion, nil
					},
				},
				tester: &fakes.FakePromotionProvider{
					UpdatePromotionStub: func(ctx context.Context, p types.Promotion) (types.Promotion, error) {
						return promotion, nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				promotion: promotion,
			},
			expectedOutput: promotion,
		},
		{
			name: "it should fail update promotion not found",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					PromotionUpdateStub: func(ctx context.Context, p types.Promotion) (types.Promotion, error) {
						return types.Promotion{}, pgx.ErrNoRows
					},
				},
				tester: &fakes.FakePromotionProvider{
					UpdatePromotionStub: func(ctx context.Context, p types.Promotion) (types.Promotion, error) {
						return types.Promotion{}, pgx.ErrNoRows
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				promotion: promotion,
			},
			expectedOutput: types.Promotion{},
			expectedError:  pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := promotions.New(tt.fields.persistentStore)
			res, err := c.UpdatePromotion(context.Background(), tt.args.promotion)

			require.ErrorIs(t, err, tt.expectedError)
			if tt.expectedError == nil {
				require.Equal(t, tt.expectedOutput, res)
			}
		})
	}
}

func TestDeletePromotion(t *testing.T) {
	type args struct {
		promotionID uuid.UUID
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	tests := []struct {
		name          string
		fields        fields
		vars          map[string]string
		args          args
		expectedError error
	}{
		{
			name: "it should get delete promotion",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					PromotionDeleteStub: func(ctx context.Context, u uuid.UUID) error {
						return nil
					},
				},
				tester: &fakes.FakePromotionProvider{
					DeletePromotionStub: func(ctx context.Context, u uuid.UUID) error {
						return nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{promotionID: ID},
		},
		{
			name: "it should fail to delete not found",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					PromotionDeleteStub: func(ctx context.Context, u uuid.UUID) error {
						return pgx.ErrNoRows
					},
				},
				tester: &fakes.FakePromotionProvider{
					DeletePromotionStub: func(ctx context.Context, u uuid.UUID) error {
						return pgx.ErrNoRows
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args:          args{promotionID: ID},
			expectedError: pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := promotions.New(tt.fields.persistentStore)
			err := c.DeletePromotion(context.Background(), tt.args.promotionID)

			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}
