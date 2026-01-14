package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sh-latibov/telegram-bot-go/models"
)

type Repo struct {
	db *pgx.Conn
}

func (r *Repo) IsUserExists(ctx context.Context, userID int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func New(db *pgx.Conn) *Repo {
	return &Repo{db: db}
}

//SaveUser(ctx context.Context, userID int64) error
//SaveUserCity(ctx context.Context, userID int64, city string) error
//GetUserCity(ctx context.Context, userID int64) (string, error)
//UpdateUserCity(ctx context.Context, userID int64, city string) error

func (r *Repo) GetUserCity(ctx context.Context, userID int64) (string, error) {
	var city string

	row := r.db.QueryRow(ctx, "select coalesce(city, '') from users where id=$1", userID)
	err := row.Scan(&city)
	if err != nil {
		return "", fmt.Errorf("error row.Scan: %w", err)
	}

	return city, nil
}

func (r *Repo) SaveUser(ctx context.Context, userId int64) error {
	_, err := r.db.Exec(ctx, "insert into users (id) values($1)", userId)

	if err != nil {
		return fmt.Errorf("error saving user city: %w", err)
	}

	return nil
}

func (r *Repo) UpdateUserCity(ctx context.Context, userID int64, city string) error {
	_, err := r.db.Exec(ctx, "update users set city=$1 where id=$2", city, userID)
	if err != nil {
		return fmt.Errorf("error updating user city: %w", err)
	}
	return nil

}

func (r *Repo) GetUser(ctx context.Context, userID int64) (*models.User, error) {
	user := models.User{}
	row := r.db.QueryRow(ctx, "select id, coalesce(city, ''), created_at from users where id=$1", userID)
	err := row.Scan(&user.ID, &user.City, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error row.Scan: %w", err)
	}
	return &user, nil
}
