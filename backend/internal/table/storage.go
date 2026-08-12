package table

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

func (s *Storage) Create(ctx context.Context, table models.Table) (int, error) {

	_, err := s.pool.Exec(
		ctx,
		"INSERT INTO tables (number, status) VALUES ($1, $2)",
		table.Number,
		table.Status,
	)
	return table.Number, err
}

func (s *Storage) GetAll(ctx context.Context) ([]models.Table, error) {
	tables := make([]models.Table, 0)

	rows, err := s.pool.Query(ctx,
		"SELECT number, status FROM tables")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var table models.Table
		err := rows.Scan(
			&table.Number,
			&table.Status,
		)

		if err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func (s *Storage) Delete(ctx context.Context, number int) error {
	result, err := s.pool.Exec(ctx, "DELETE FROM tables WHERE number = $1", number)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return myerrors.ErrNotFound
	}

	return nil
}
