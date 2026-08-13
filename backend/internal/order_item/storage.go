package orderitem

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

func (s *Storage) AddItem(
	ctx context.Context,
	item models.OrderItem,
) error {
	_, err := s.pool.Exec(
		ctx,
		`INSERT INTO order_items
            (order_id, menu_item_id, count, is_ready)
         VALUES ($1, $2, $3, $4)`,
		item.OrderID,
		item.MenuItemID,
		item.Count,
		item.IsReady,
	)

	return err
}

func (s *Storage) GetByOrderID(
	ctx context.Context,
	orderID int,
) ([]models.OrderItem, error) {

	items := make([]models.OrderItem, 0)

	rows, err := s.pool.Query(
		ctx,
		`SELECT order_id, menu_item_id, count, is_ready
         FROM order_items
         WHERE order_id = $1`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem

		err := rows.Scan(
			&item.OrderID,
			&item.MenuItemID,
			&item.Count,
			&item.IsReady,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Storage) IsOwnedByWaiter(ctx context.Context, orderID, waiterID int) (bool, error) {
	var owned bool
	err := s.pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1 AND waiter_id = $2)",
		orderID, waiterID,
	).Scan(&owned)
	return owned, err
}

func (s *Storage) SetReady(
	ctx context.Context,
	orderID int,
	menuItemID int,
	ready bool,
) error {
	res, err := s.pool.Exec(
		ctx,
		`UPDATE order_items
         SET is_ready = $1
         WHERE order_id = $2
           AND menu_item_id = $3`,
		ready,
		orderID,
		menuItemID,
	)

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return myerrors.ErrNotFound
	}

	return nil
}

func (s *Storage) UpdateCount(
	ctx context.Context,
	orderID int,
	menuItemID int,
	count int,
) error {
	res, err := s.pool.Exec(
		ctx,
		`UPDATE order_items
         SET count = $1
         WHERE order_id = $2
           AND menu_item_id = $3`,
		count,
		orderID,
		menuItemID,
	)

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return myerrors.ErrNotFound
	}

	return nil
}

func (s *Storage) DeleteItem(
	ctx context.Context,
	orderID int,
	menuItemID int,
) error {
	res, err := s.pool.Exec(
		ctx,
		`DELETE FROM order_items
         WHERE order_id = $1
           AND menu_item_id = $2`,
		orderID,
		menuItemID,
	)

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return myerrors.ErrNotFound
	}

	return nil
}
