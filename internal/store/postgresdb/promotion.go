package postgresdb

import (
	"context"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (q *Queries) PromotionCreate(ctx context.Context, promotion types.Promotion) (types.Promotion, error) {
	query := `
		INSERT INTO promotions (
			id,
			title,
			description,
			amount,
			is_active
		) VALUES ($1, $2, $3, $4, $5)`

	_, err := q.db.Exec(ctx, query,
		promotion.ID,
		promotion.Title,
		promotion.Description,
		promotion.Amount,
		promotion.IsActive,
	)

	return promotion, err
}

func (q *Queries) PromotionGetByID(ctx context.Context, id uuid.UUID) (types.Promotion, error) {
	var (
		promotion types.Promotion
		query     = `
		SELECT 
			id,
			title,
			description,
			amount,
			is_active,
			created,
			updated
		FROM promotions 
		WHERE id = $1`
	)

	err := q.db.QueryRow(ctx, query, id).Scan(
		&promotion.ID,
		&promotion.Title,
		&promotion.Description,
		&promotion.Amount,
		&promotion.IsActive,
		&promotion.Created,
		&promotion.Updated,
	)

	return promotion, err
}

func (q *Queries) PromotionGetByType(ctx context.Context, promotionType types.PromotionType) (types.Promotion, error) {
	var (
		promotion types.Promotion
		query     = `
		SELECT 
			id,
			title,
			description,
			amount,
			is_active,
			created,
			updated
		FROM promotions 
		WHERE type = $1
		LIMIT 1`
	)

	err := q.db.QueryRow(ctx, query, promotionType).Scan(
		&promotion.ID,
		&promotion.Title,
		&promotion.Description,
		&promotion.Amount,
		&promotion.IsActive,
		&promotion.Created,
		&promotion.Updated,
	)

	return promotion, err
}

func (q *Queries) GetPromotions(ctx context.Context) ([]types.Promotion, error) {
	var (
		promotions []types.Promotion
		args       []any

		query = `
		SELECT 
			id,
			title,
			description,
			amount,
			is_active,
			created,
			updated
		FROM promotions`
	)

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var promotion types.Promotion
		err := rows.Scan(
			&promotion.ID,
			&promotion.Title,
			&promotion.Description,
			&promotion.Amount,
			&promotion.IsActive,
			&promotion.Created,
			&promotion.Updated,
		)

		if err != nil {
			return nil, err
		}

		promotions = append(promotions, promotion)
	}

	return promotions, rows.Err()
}

func (q *Queries) PromotionUpdate(ctx context.Context, promotion types.Promotion) (types.Promotion, error) {
	query := `
		UPDATE promotions SET
			title = $1,
			description = $2,
			amount = $3,
			is_active = $4
		WHERE id = $5`

	res, err := q.db.Exec(
		ctx,
		query,
		&promotion.Title,
		&promotion.Description,
		&promotion.Amount,
		&promotion.IsActive,
		&promotion.ID,
	)

	if err != nil {
		return promotion, err
	}

	if res.RowsAffected() == 0 {
		return promotion, pgx.ErrNoRows
	}

	return promotion, err
}

func (q *Queries) PromotionDelete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM promotions WHERE id = $1`

	res, err := q.db.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
