package userpromotion_test

import (
	"context"
	"testing"
	"time"

	userpromotion "github.com/Jozzo6/casino_loyalty_reward_system/internal/component/user_promotion"
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
	tester          userpromotion.UserPromotionProvider
}

var (
	fixedTime    = time.Date(2025, time.March, 19, 8, 15, 55, 706491000, time.Local)
	fixedEndTime = fixedTime.Add(time.Hour)
)

func TestAddPromotion(t *testing.T) {
	type args struct {
		userPromotion types.UserPromotion
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)
	promotionID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)
	userID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	tests := []struct {
		name           string
		fields         fields
		vars           map[string]string
		args           args
		expectedError  error
		expectedOutput types.UserPromotion
	}{
		{
			name: "it should add promotion",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					PromotionGetByIDStub: func(ctx context.Context, u uuid.UUID) (types.Promotion, error) {
						return types.Promotion{IsActive: true}, nil
					},
					AddPromotionStub: func(ctx context.Context, up types.UserPromotion) (types.UserPromotion, error) {
						return types.UserPromotion{
							ID:          ID,
							UserID:      userID,
							PromotionID: promotionID,
							StartDate:   fixedTime,
							EndDate:     fixedEndTime,
							Claimed:     nil,
						}, nil
					},
				},
				tester: &fakes.FakeUserPromotionProvider{
					AddPromotionStub: func(ctx context.Context, up types.UserPromotion) (types.UserPromotion, error) {
						return types.UserPromotion{
							ID:          ID,
							UserID:      userID,
							PromotionID: promotionID,
							StartDate:   fixedTime,
							EndDate:     fixedEndTime,
							Claimed:     nil,
						}, nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				userPromotion: types.UserPromotion{
					ID:          ID,
					UserID:      userID,
					PromotionID: promotionID,
					StartDate:   fixedTime,
					EndDate:     fixedEndTime,
					Claimed:     nil,
				},
			},
			expectedOutput: types.UserPromotion{
				ID:          ID,
				UserID:      userID,
				PromotionID: promotionID,
				StartDate:   fixedTime,
				EndDate:     fixedEndTime,
				Claimed:     nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := userpromotion.New(tt.fields.persistentStore, tt.fields.pubsub)
			res, err := c.AddPromotion(context.Background(), tt.args.userPromotion)

			require.ErrorIs(t, err, tt.expectedError)
			if tt.expectedError == nil {
				require.Equal(t, tt.expectedOutput, res)
			}
		})
	}
}

func TestAddWelcomePromotion(t *testing.T) {
	type args struct {
		userPromotion types.UserPromotion
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)
	promotionID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)
	userID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	tests := []struct {
		name           string
		fields         fields
		vars           map[string]string
		args           args
		expectedError  error
		expectedOutput types.UserPromotion
	}{
		{
			name: "it should add welcome promotion",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					PromotionGetByIDStub: func(ctx context.Context, u uuid.UUID) (types.Promotion, error) {
						return types.Promotion{IsActive: true}, nil
					},
					AddPromotionStub: func(ctx context.Context, up types.UserPromotion) (types.UserPromotion, error) {
						return types.UserPromotion{
							ID:          ID,
							UserID:      userID,
							PromotionID: promotionID,
							StartDate:   fixedTime,
							EndDate:     fixedEndTime,
							Claimed:     nil,
						}, nil
					},
				},
				tester: &fakes.FakeUserPromotionProvider{
					AddPromotionStub: func(ctx context.Context, up types.UserPromotion) (types.UserPromotion, error) {
						return types.UserPromotion{
							ID:          ID,
							UserID:      userID,
							PromotionID: promotionID,
							StartDate:   fixedTime,
							EndDate:     fixedEndTime,
							Claimed:     nil,
						}, nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				userPromotion: types.UserPromotion{
					ID:          ID,
					UserID:      userID,
					PromotionID: promotionID,
					StartDate:   fixedTime,
					EndDate:     fixedEndTime,
					Claimed:     nil,
				},
			},
			expectedOutput: types.UserPromotion{
				ID:          ID,
				UserID:      userID,
				PromotionID: promotionID,
				StartDate:   fixedTime,
				EndDate:     fixedEndTime,
				Claimed:     nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := userpromotion.New(tt.fields.persistentStore, tt.fields.pubsub)
			res, err := c.AddPromotion(context.Background(), tt.args.userPromotion)

			require.ErrorIs(t, err, tt.expectedError)
			if tt.expectedError == nil {
				require.Equal(t, tt.expectedOutput, res)
			}
		})
	}
}

func TestGetUserPromotions(t *testing.T) {
	type args struct {
		userID uuid.UUID
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)
	promotionID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)
	userID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	tests := []struct {
		name           string
		fields         fields
		vars           map[string]string
		args           args
		expectedError  error
		expectedOutput []types.UserPromotion
	}{
		{
			name: "it should user promotion by user id",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					GetUserPromotionsStub: func(ctx context.Context, u uuid.UUID) ([]types.UserPromotion, error) {
						return []types.UserPromotion{
							{
								ID:          ID,
								UserID:      userID,
								PromotionID: promotionID,
								StartDate:   fixedTime,
								EndDate:     fixedEndTime,
								Claimed:     nil,
							},
						}, nil
					},
				},
				tester: &fakes.FakeUserPromotionProvider{
					GetUserPromotionsStub: func(ctx context.Context, u uuid.UUID) ([]types.UserPromotion, error) {
						return []types.UserPromotion{
							{
								ID:          ID,
								UserID:      userID,
								PromotionID: promotionID,
								StartDate:   fixedTime,
								EndDate:     fixedEndTime,
								Claimed:     nil,
							},
						}, nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				userID: userID,
			},
			expectedOutput: []types.UserPromotion{
				{
					ID:          ID,
					UserID:      userID,
					PromotionID: promotionID,
					StartDate:   fixedTime,
					EndDate:     fixedEndTime,
					Claimed:     nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := userpromotion.New(tt.fields.persistentStore, tt.fields.pubsub)
			res, err := c.GetUserPromotions(context.Background(), tt.args.userID)

			require.ErrorIs(t, err, tt.expectedError)
			if tt.expectedError == nil {
				require.Equal(t, tt.expectedOutput, res)
			}
		})
	}
}

func TestGetUserPromotionByID(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)
	promotionID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)
	userID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	tests := []struct {
		name           string
		fields         fields
		vars           map[string]string
		args           args
		expectedError  error
		expectedOutput types.UserPromotion
	}{
		{
			name: "it should get user promotion by id",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					GetUserPromotionByIDStub: func(ctx context.Context, u uuid.UUID) (types.UserPromotion, error) {
						return types.UserPromotion{
							ID:          ID,
							UserID:      userID,
							PromotionID: promotionID,
							StartDate:   fixedTime,
							EndDate:     fixedEndTime,
							Claimed:     nil,
						}, nil
					},
				},
				tester: &fakes.FakeUserPromotionProvider{
					GetUserPromotionByIDStub: func(ctx context.Context, u uuid.UUID) (types.UserPromotion, error) {
						return types.UserPromotion{
							ID:          ID,
							UserID:      userID,
							PromotionID: promotionID,
							StartDate:   fixedTime,
							EndDate:     fixedEndTime,
							Claimed:     nil,
						}, nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				ID: ID,
			},
			expectedOutput: types.UserPromotion{
				ID:          ID,
				UserID:      userID,
				PromotionID: promotionID,
				StartDate:   fixedTime,
				EndDate:     fixedEndTime,
				Claimed:     nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := userpromotion.New(tt.fields.persistentStore, tt.fields.pubsub)
			res, err := c.GetUserPromotionByID(context.Background(), tt.args.ID)

			require.ErrorIs(t, err, tt.expectedError)
			if tt.expectedError == nil {
				require.Equal(t, tt.expectedOutput, res)
			}
		})
	}
}

func TestClaimPromotion(t *testing.T) {
	type args struct {
		ID uuid.UUID
	}

	ID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)
	promotionID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)
	userID, err := uuid.Parse("460aec7e-7d58-42fd-93b8-bca05a77bbf5")
	require.NoError(t, err)

	tests := []struct {
		name          string
		fields        fields
		vars          map[string]string
		args          args
		expectedError error
	}{
		{
			name: "it should claim promotion",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					WithTxStub: func(ctx context.Context) (store.Persistent, error) {
						return &fakes.FakePersistent{
							GetUserPromotionByIDStub: func(ctx context.Context, u uuid.UUID) (types.UserPromotion, error) {
								return types.UserPromotion{
									ID:          ID,
									UserID:      userID,
									PromotionID: promotionID,
									StartDate:   time.Now(),
									EndDate:     time.Date(2026, time.March, 19, 8, 15, 55, 706491000, time.Local),
									Claimed:     nil,
									Promotion: &types.Promotion{
										ID:       promotionID,
										Amount:   10,
										IsActive: true,
									},
									User: &types.User{ID: userID},
								}, nil
							},
							ClaimPromotionStub: func(ctx context.Context, u uuid.UUID) error {
								return nil
							},
							UserGetByStub: func(ctx context.Context, uf types.UserFilter) (types.User, error) {
								return types.User{ID: userID, Balance: 10}, nil
							},
							UserBalanceUpdateStub: func(ctx context.Context, u uuid.UUID, f float64) (types.User, error) {
								return types.User{
									ID:      userID,
									Balance: 20,
								}, nil
							},
						}, err
					},
				},
				tester: &fakes.FakeUserPromotionProvider{
					ClaimPromotionStub: func(ctx context.Context, u uuid.UUID) error {
						return nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				ID: ID,
			},
		},
		{
			name: "it should fail to claim promotion claimed already",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					WithTxStub: func(ctx context.Context) (store.Persistent, error) {
						return &fakes.FakePersistent{
							GetUserPromotionByIDStub: func(ctx context.Context, u uuid.UUID) (types.UserPromotion, error) {
								return types.UserPromotion{
									ID:          ID,
									UserID:      userID,
									PromotionID: promotionID,
									StartDate:   fixedTime,
									EndDate:     fixedEndTime,
									Claimed:     &fixedTime,
									Promotion: &types.Promotion{
										ID:       promotionID,
										Amount:   10,
										IsActive: true,
									},
								}, nil
							},
							ClaimPromotionStub: func(ctx context.Context, u uuid.UUID) error {
								return nil
							},
							UserGetByStub: func(ctx context.Context, uf types.UserFilter) (types.User, error) {
								return types.User{ID: userID, Balance: 10}, nil
							},
							UserBalanceUpdateStub: func(ctx context.Context, u uuid.UUID, f float64) (types.User, error) {
								return types.User{
									ID:      userID,
									Balance: 20,
								}, nil
							},
						}, nil
					},
				},
				tester: &fakes.FakeUserPromotionProvider{
					ClaimPromotionStub: func(ctx context.Context, u uuid.UUID) error {
						return types.ErrPromotionClaimed
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				ID: ID,
			},
			expectedError: types.ErrPromotionClaimed,
		},
		{
			name: "it should fail to claim promotion  promotion not active",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					WithTxStub: func(ctx context.Context) (store.Persistent, error) {
						return &fakes.FakePersistent{
							GetUserPromotionByIDStub: func(ctx context.Context, u uuid.UUID) (types.UserPromotion, error) {
								return types.UserPromotion{
									ID:          ID,
									UserID:      userID,
									PromotionID: promotionID,
									StartDate:   fixedTime,
									EndDate:     fixedEndTime,
									Claimed:     nil,
									Promotion: &types.Promotion{
										ID:       promotionID,
										Amount:   10,
										IsActive: false,
									},
								}, nil
							},
							ClaimPromotionStub: func(ctx context.Context, u uuid.UUID) error {
								return nil
							},
							UserGetByStub: func(ctx context.Context, uf types.UserFilter) (types.User, error) {
								return types.User{ID: userID, Balance: 10}, nil
							},
							UserBalanceUpdateStub: func(ctx context.Context, u uuid.UUID, f float64) (types.User, error) {
								return types.User{
									ID:      userID,
									Balance: 20,
								}, nil
							},
						}, nil
					},
				},
				tester: &fakes.FakeUserPromotionProvider{
					ClaimPromotionStub: func(ctx context.Context, u uuid.UUID) error {
						return types.ErrPromotionNoLongerActive
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				ID: ID,
			},
			expectedError: types.ErrPromotionNoLongerActive,
		},
		{
			name: "it should fail to claim promotion not started",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					WithTxStub: func(ctx context.Context) (store.Persistent, error) {
						return &fakes.FakePersistent{
							GetUserPromotionByIDStub: func(ctx context.Context, u uuid.UUID) (types.UserPromotion, error) {
								return types.UserPromotion{
									ID:          ID,
									UserID:      userID,
									PromotionID: promotionID,
									StartDate:   time.Date(3000, time.March, 19, 8, 15, 55, 706491000, time.Local),
									EndDate:     time.Date(3001, time.March, 19, 8, 15, 55, 706491000, time.Local),
									Claimed:     nil,
									Promotion: &types.Promotion{
										ID:       promotionID,
										Amount:   10,
										IsActive: true,
									},
								}, nil
							},
							ClaimPromotionStub: func(ctx context.Context, u uuid.UUID) error {
								return nil
							},
							UserGetByStub: func(ctx context.Context, uf types.UserFilter) (types.User, error) {
								return types.User{ID: userID, Balance: 10}, nil
							},
							UserBalanceUpdateStub: func(ctx context.Context, u uuid.UUID, f float64) (types.User, error) {
								return types.User{
									ID:      userID,
									Balance: 20,
								}, nil
							},
						}, nil
					},
				},
				tester: &fakes.FakeUserPromotionProvider{
					ClaimPromotionStub: func(ctx context.Context, u uuid.UUID) error {
						return types.ErrPromotionNotStarted
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				ID: ID,
			},
			expectedError: types.ErrPromotionNotStarted,
		},
		{
			name: "it should fail to claim promotion expired",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					WithTxStub: func(ctx context.Context) (store.Persistent, error) {
						return &fakes.FakePersistent{
							GetUserPromotionByIDStub: func(ctx context.Context, u uuid.UUID) (types.UserPromotion, error) {
								return types.UserPromotion{
									ID:          ID,
									UserID:      userID,
									PromotionID: promotionID,
									StartDate:   time.Date(2025, time.March, 17, 8, 15, 55, 706491000, time.Local),
									EndDate:     time.Date(2025, time.March, 18, 8, 15, 55, 706491000, time.Local),
									Claimed:     nil,
									Promotion: &types.Promotion{
										ID:       promotionID,
										Amount:   10,
										IsActive: true,
									},
									User: &types.User{
										ID: userID,
									},
								}, nil
							},
							ClaimPromotionStub: func(ctx context.Context, u uuid.UUID) error {
								return nil
							},
							UserGetByStub: func(ctx context.Context, uf types.UserFilter) (types.User, error) {
								return types.User{ID: userID, Balance: 10}, nil
							},
							UserBalanceUpdateStub: func(ctx context.Context, u uuid.UUID, f float64) (types.User, error) {
								return types.User{
									ID:      userID,
									Balance: 20,
								}, nil
							},
						}, nil
					},

				},
				tester: &fakes.FakeUserPromotionProvider{
					ClaimPromotionStub: func(ctx context.Context, u uuid.UUID) error {
						return types.ErrPromotionExpired
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				ID: ID,
			},
			expectedError: types.ErrPromotionExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := userpromotion.New(tt.fields.persistentStore, tt.fields.pubsub)
			err := c.ClaimPromotion(context.Background(), tt.args.ID)

			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestDeleteUserPromotion(t *testing.T) {
	type args struct {
		ID uuid.UUID
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
			name: "it should delete user promotion",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					DeleteUserPromotionStub: func(ctx context.Context, u uuid.UUID) error {
						return nil
					},
				},
				tester: &fakes.FakeUserPromotionProvider{
					DeleteUserPromotionStub: func(ctx context.Context, u uuid.UUID) error {
						return nil
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				ID: ID,
			},
		},
		{
			name: "it should fail to delete user promotion",
			fields: fields{
				persistentStore: &fakes.FakePersistent{
					DeleteUserPromotionStub: func(ctx context.Context, u uuid.UUID) error {
						return pgx.ErrNoRows
					},
				},
				tester: &fakes.FakeUserPromotionProvider{
					DeleteUserPromotionStub: func(ctx context.Context, u uuid.UUID) error {
						return pgx.ErrNoRows
					},
				},
				pubsub: &fakes.FakePubSub{},
			},
			args: args{
				ID: ID,
			},
			expectedError: pgx.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := userpromotion.New(tt.fields.persistentStore, tt.fields.pubsub)
			err := c.DeleteUserPromotion(context.Background(), tt.args.ID)

			require.ErrorIs(t, err, tt.expectedError)
		})
	}
}
