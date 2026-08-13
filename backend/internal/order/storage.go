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
		`INSERT INTO orders (table_number, waiter_id, is_served, is_payed, created_at, closed_at)
		VALUES($1, $2, $3, $4, $5, $6) RETURNING id;`,
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
		"SELECT id, table_number, waiter_id, is_served, is_payed, created_at, closed_at FROM orders",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

	return orders, nil
}

func (s *Storage) GetByTableNumber(ctx context.Context, number int) (models.Order, error) {

	var order models.Order

	res := s.pool.QueryRow(
		ctx,
		`SELECT id, waiter_id, is_served, is_payed, created_at, closed_at FROM orders
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

func (s *Storage) GetByWaiterId(ctx context.Context, id int) ([]models.Order, error) {

	orders := make([]models.Order, 0)

	rows, err := s.pool.Query(
		ctx,
		`SELECT id, table_number, waiter_id, is_served, is_payed, created_at, closed_at FROM orders
		WHERE waiter_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

	return orders, nil

}

func (s *Storage) Update(ctx context.Context, order OrderRequestBody) error {

	res, err := s.pool.Exec(
		ctx,
		"UPDATE orders set table_number = $1, waiter_id = $2, is_served = $3, is_payed = $4 WHERE id = $5",
		order.TableNumber,
		order.WaiterID,
		order.IsServed,
		order.IsPayed,
		order.ID,
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
		"DELETE from orders WHERE id = $1", id,
	)

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return myerrors.ErrNotFound
	}

	return nil
}
