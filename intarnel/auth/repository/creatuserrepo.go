package repository

import (
	"context"
	"ginframework/intarnel/auth"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateNewUser(ctx context.Context, pool *pgxpool.Pool, user auth.GoogleUserInfo) (*auth.GoogleUserInfo, error) {
	var err error

var query string = `
    INSERT INTO users (id, email, verified_email, name, picture)
    VALUES ($1, $2, $3, $4, $5)
    ON CONFLICT (id) 
    DO UPDATE SET 
        email = EXCLUDED.email, 
        verified_email = EXCLUDED.verified_email, 
        name = EXCLUDED.name, 
        picture = EXCLUDED.picture, 
        updated_at = NOW()
    RETURNING id, email, verified_email, name, picture, created_at, updated_at
`

	err = pool.QueryRow(ctx, query, user.ID, user.Email, user.VerifiedEmail, user.Name, user.Picture).Scan(
		&user.ID,
		&user.Email,
		&user.VerifiedEmail,
		&user.Name,
		&user.Picture,
		&user.Created_at,
		&user.Updated_at,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
