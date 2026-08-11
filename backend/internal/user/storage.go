package user

import (
	"backend/internal/models"
	myerrors "backend/internal/my_errors"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{
		pool: pool,
	}
}

func (s *Storage) Create(ctx context.Context, user models.User) (int, error) {
	var id int
	err := s.pool.QueryRow(
		ctx,
		"INSERT INTO users (nickname, phone_number, role) VALUES($1, $2, $3) RETURNING id",
		user.Nickname,
		user.PhoneNumber,
		user.Role,
	).Scan(&id)
	return id, err
}

func (s *Storage) GetUsers(ctx context.Context) ([]models.User, error) {
	users := make([]models.User, 0)

	rows, err := s.pool.Query(ctx,
		"SELECT id, nickname, phone_number, role FROM users")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Nickname,
			&user.PhoneNumber,
			&user.Role,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Storage) IsRegistered(ctx context.Context, nickname string, phone string) (string, error) {
	var role string

	err := s.pool.QueryRow(ctx,
		"SELECT role FROM users WHERE nickname = $1 OR phone_number = $2 )",
		nickname,
		phone,
	).Scan(&role)

	if err != nil {
		return "", err
	}

	return role, nil
}

func (s *Storage) Delete(ctx context.Context, id int) error {
	result, err := s.pool.Exec(ctx,
		"DELETE FROM users WHERE i = $1", id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return myerrors.ErrNotFound
	}

	return nil
}
