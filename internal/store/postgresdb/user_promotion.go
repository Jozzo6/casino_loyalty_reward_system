package postgresdb

import (
	"context"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"
	"github.com/google/uuid"
)

func (q *Queries) AddPromotion(ctx context.Context, userPromotion types.UserPromotion) (types.UserPromotion, error) {
	query := `
		INSERT INTO users_promotions (
			id,
			user_id,
			promotion_id,
			claimed,
			start_date,
			end_date
		) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := q.db.Exec(ctx, query,
		&userPromotion.ID,
		&userPromotion.UserID,
		&userPromotion.PromotionID,
		&userPromotion.Claimed,
		&userPromotion.StartDate,
		&userPromotion.EndDate,
	)

	return userPromotion, err
}

func (q *Queries) ClaimPromotion(ctx context.Context, userPromotionID uuid.UUID) error {
	query := `UPDATE users_promotions SET claimed = NOW() WHERE id = $1 `

	_, err := q.db.Exec(ctx, query, &userPromotionID)

	return err
}

func (q *Queries) DeleteUserPromotion(ctx context.Context, userPromotionID uuid.UUID) error {
	query := `DELETE FROM users_promotions WHERE id = $1 LIMIT 1`

	_, err := q.db.Exec(ctx, query, &userPromotionID)

	return err
}

func (q *Queries) GetUserPromotionByID(ctx context.Context, userPromotionID uuid.UUID) (types.UserPromotion, error) {
	var (
		userPromotion types.UserPromotion
		query         = `SELECT 
			up.id,
			up.user_id,
			up.promotion_id,
			up.claimed,
			up.start_date,
			up.end_date,
			json_build_object(
				'id', p.id,
				'title', p.title,
				'description', p.description,
				'amount', p.amount,
				'is_active', p.is_active,
				'created', p.created,
				'updated', p.updated
			) as promotion
			FROM users_promotions up
			INNER JOIN promotions p on p.id = up.promotion_id 
			WHERE up.id = $1 
			LIMIT 1`
	)

	err := q.db.QueryRow(ctx, query, userPromotionID).Scan(
		&userPromotion.ID,
		&userPromotion.UserID,
		&userPromotion.PromotionID,
		&userPromotion.Claimed,
		&userPromotion.StartDate,
		&userPromotion.EndDate,
		&userPromotion.Promotion,
	)

	return userPromotion, err
}

func (q *Queries) GetUserPromotions(ctx context.Context, userPromotionID uuid.UUID) ([]types.UserPromotion, error) {
	var (
		userPromotions []types.UserPromotion
		query          = `SELECT 
			up.id,
			up.user_id,
			up.promotion_id,
			up.claimed,
			up.start_date,
			up.end_date,
			json_build_object(
				'id', p.id,
				'title', p.title,
				'description', p.description,
				'amount', p.amount,
				'is_active', p.is_active,
				'created', p.created,
				'updated', p.updated
			) as promotion
			FROM users_promotions up
			INNER JOIN promotions p on p.id = up.promotion_id
			WHERE up.user_id = $1`
	)

	rows, err := q.db.Query(ctx, query, userPromotionID)
	if err != nil {
		return []types.UserPromotion{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var userPromotion types.UserPromotion
		err := rows.Scan(
			&userPromotion.ID,
			&userPromotion.UserID,
			&userPromotion.PromotionID,
			&userPromotion.Claimed,
			&userPromotion.StartDate,
			&userPromotion.EndDate,
			&userPromotion.Promotion,
		)

		if err != nil {
			return nil, err
		}

		userPromotions = append(userPromotions, userPromotion)
	}

	return userPromotions, rows.Err()
}
