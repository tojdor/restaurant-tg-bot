package order

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

type OrderRequestBody struct {
	ID          int  `json:"id"`
	TableNumber int  `json:"table_number"`
	WaiterID    int  `json:"waiter_id"`
	IsServed    bool `json:"is_served"`
	IsPayed     bool `json:"is_payed"`
}

func (s *Storage) Create(ctx context.Context, order models.Order) (int, error) {

	var id int
	err := s.pool.QueryRow(
		ctx,
		`INSERT INTO order (table_number, waiter_id, is_served, is_payed, created_at, closed_at)
		VALUES($1, $2, $3, $4, $5, $6);`,
		&order.TableNumber,
		&order.WaiterID,
		&order.IsServed,
		&order.IsPayed,
		&order.CreatedAt,
		&order.ClosedAt,
	).Scan(&id)

	return id, err
}

func (s *Storage) GetAll(ctx context.Context) ([]models.Order, error) {

	orders := make([]models.Order, 0)

	rows, err := s.pool.Query(
		ctx,
		"SELECT id, table_number, waiter_id, is_served, is_payed, created_at, closed_at FROM order",
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.TableNumber,
			&order.WaiterID,
			&order.IsServed,
			&order.IsPayed,
			&order.CreatedAt,
			&order.ClosedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	defer rows.Close()
	return orders, nil
}

func (s *Storage) GetByTableNumber(ctx context.Context, number int) (models.Order, error) {

	var order models.Order

	res := s.pool.QueryRow(
		ctx,
		`SELECT id, waiter_id, is_served, is_payed, created_at, closed_at FROM order
		WHERE table_number = $1`, number,
	)
	err := res.Scan(
		&order.ID,
		&order.WaiterID,
		&order.IsServed,
		&order.IsPayed,
		&order.CreatedAt,
		&order.ClosedAt,
	)

	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (s *Storage) GetTablesByWaiterId(ctx context.Context, id int) ([]models.Order, error) {

	orders := make([]models.Order, 0)

	rows, err := s.pool.Query(
		ctx,
		`SELECT id, table_number, waiter_id, is_served, is_payed, created_at, closed_at FROM order
		WHERE waiter_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.TableNumber,
			&order.WaiterID,
			&order.IsServed,
			&order.IsPayed,
			&order.CreatedAt,
			&order.ClosedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	defer rows.Close()
	return orders, nil

}

func (s *Storage) Update(ctx context.Context, order OrderRequestBody) error {

	res, err := s.pool.Exec(
		ctx,
		"UPDATE order set id = $1, table_number = $2, waiter_id = $3, is_served = $4, is_payed = $5",
		order.ID,
		order.TableNumber,
		order.WaiterID,
		order.IsServed,
		order.IsPayed,
	)

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return myerrors.ErrNotFound
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id int) error {

	res, err := s.pool.Exec(
		ctx,
		"DELETE from order WHERE id = $1", id,
	)

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return myerrors.ErrNotFound
	}

	return nil
}
