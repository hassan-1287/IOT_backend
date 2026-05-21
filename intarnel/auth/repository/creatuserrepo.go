package repository

import (
	"context"
	"ginframework/intarnel/auth"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateNewUser(pool *pgxpool.Pool, user auth.GoogleUserInfo) (*auth.GoogleUserInfo, error) {
	var ctx context.Context
	var cancle context.CancelFunc
	var err error

	ctx, cancle = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	exsiquery := `
		SELECT id  FROM users WHERE id = $1
	`
	err = pool.QueryRow(ctx, exsiquery, user.ID).Scan(&user.ID)
	if err == nil {
		return &user, nil
	}

	query := `
		INSERT INTO users (id , email , verified_email, name , picture)
		VALUES ($1 , $2 , $3 , $4,$5)
		RETURNING ID , id,email,verified_email,name, picture , created_at , updated_at

	`

	err = pool.QueryRow(ctx, query, user.ID, user.Email, user.VerifiedEmail, user.Name, user.Picture).Scan(&user.ID, &user.Email, &user.Name, &user.VerifiedEmail, &user.Picture)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
