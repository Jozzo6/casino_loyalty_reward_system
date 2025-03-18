package postgresdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (q *Queries) UserCreate(ctx context.Context, user types.User) (types.User, error) {
	query := `
INSERT INTO users (
	id,
	name,
	email,
	password,
	role,
	created,
	updated
) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := q.db.Exec(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		user.Created,
		user.Updated,
	)

	return user, err
}

func (q *Queries) UserGetBy(ctx context.Context, filter types.UserFilter) (types.User, error) {
	var (
		user        types.User
		whereClause []string
		args        []any

		query = `
		SELECT
			id,
			email,
			name,
			password,
			role,
			created,
			updated
		FROM users
		WHERE
			%s
		LIMIT 1`
	)

	if filter.ByID.Valid {
		whereClause = append(whereClause, fmt.Sprintf("id = $%d", len(args)+1))
		args = append(args, filter.ByID)
	}

	if filter.ByEmail != nil {
		whereClause = append(whereClause, fmt.Sprintf("email = $%d", len(args)+1))
		args = append(args, filter.ByEmail)
	}

	if len(args) == 0 {
		return types.User{}, ErrorNoFiltersProvided
	}

	query = fmt.Sprintf(query, strings.Join(whereClause, " AND "))

	err := q.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.Role,
		&user.Created,
		&user.Updated,
	)

	return user, err
}

func (q *Queries) GetUsers(ctx context.Context) ([]types.User, error) {
	var (
		users []types.User
		args  []any

		query = `
		SELECT
			id,
			email,
			name,
			role,
			created,
			updated
		FROM users`
	)

	rows, err := q.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user types.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Role,
			&user.Created,
			&user.Updated,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, rows.Err()
}

func (q *Queries) UserUpdate(ctx context.Context, user types.User) (types.User, error) {
	query := `
		UPDATE users SET
			email = $1,
			name = $2,
			role = $3
		WHERE id = $4`

	res, err := q.db.Exec(
		ctx,
		query,
		user.Email,
		user.Name,
		user.Role,
		user.ID,
	)

	if err != nil {
		return user, err
	}

	if res.RowsAffected() == 0 {
		return user, pgx.ErrNoRows
	}

	return user, err
}

func (q *Queries) UserDelete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

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

func (q *Queries) UserBalanceUpdate(ctx context.Context, id uuid.UUID, newBalance float64) (types.User, error) {
	query := `UPDATE users
			SET balance = $1
			WHERE id = $2 
			RETURNING 
				id, 
				email, 
				name, 
				role, 
				balance, 
				created, 
				updated`

	var user types.User
	err := q.db.QueryRow(
		ctx,
		query,
		newBalance,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.Balance,
		&user.Created,
		&user.Updated,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return user, pgx.ErrNoRows
		}
		return user, err
	}

	return user, nil
}

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
