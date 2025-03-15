package postgresdb

import (
	"casino_loyalty_reward_system/internal/types"
	"context"
	"fmt"
	"strings"
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

	return users, err
}

