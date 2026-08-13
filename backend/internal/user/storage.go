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
		"INSERT INTO users (telegram_user_id, nickname, phone_number, role) VALUES($1, $2, $3, $4) RETURNING id",
		user.TelegramUserID,
		user.Nickname,
		user.PhoneNumber,
		user.Role,
	).Scan(&id)
	return id, err
}

func (s *Storage) GetUsers(ctx context.Context) ([]models.User, error) {
	users := make([]models.User, 0)

	rows, err := s.pool.Query(ctx,
		"SELECT id, telegram_user_id, nickname, phone_number, role FROM users")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.TelegramUserID,
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

func (s *Storage) Delete(ctx context.Context, id int) error {
	result, err := s.pool.Exec(ctx,
		"DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return myerrors.ErrNotFound
	}

	return nil
}
