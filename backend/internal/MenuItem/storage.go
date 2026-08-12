package menuitem

import (
	"backend/internal/models"
	errors "backend/internal/my_errors"
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

func (s *Storage) Create(ctx context.Context, menu models.MenuItem) (int, error) {

	_, err := s.GetMenuItemByName(ctx, menu.Name)
	if err == nil {
		return 0, myerrors.ErrAlreadyExists
	}

	var id int
	err = s.pool.QueryRow(
		ctx,
		"INSERT INTO menu_item (name, price, category) VALUES($1, $2, $3) RETURNING id",
		menu.Name,
		menu.Price,
		menu.Category,
	).Scan(id)

	return id, err
}

func (s *Storage) GetMenuItems(ctx context.Context) ([]models.MenuItem, error) {

	menu_items := make([]models.MenuItem, 0)

	rows, err := s.pool.Query(
		ctx,
		"SELECT id, name, price, category FROM menu_item")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var menu_item models.MenuItem
		err := rows.Scan(
			&menu_item.ID,
			&menu_item.Name,
			&menu_item.Price,
			&menu_item.Category,
		)
		if err != nil {
			return nil, err
		}
		menu_items = append(menu_items, menu_item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return menu_items, nil
}

func (s *Storage) GetMenuItemByID(ctx context.Context, id int) (models.MenuItem, error) {

	var menu_item models.MenuItem

	err := s.pool.QueryRow(
		ctx,
		"SELECT id, name, price, category FROM menu_item WHERE id=$1", id,
	).Scan(
		&menu_item.ID,
		&menu_item.Name,
		&menu_item.Price,
		&menu_item.Category,
	)

	if err != nil {
		return models.MenuItem{}, err
	}

	return menu_item, err
}

func (s *Storage) GetMenuItemByName(ctx context.Context, name string) (models.MenuItem, error) {

	var menu_item models.MenuItem

	err := s.pool.QueryRow(
		ctx,
		"SELECT id, name, price, category FROM menu_item WHERE name=$1", name,
	).Scan(
		&menu_item.ID,
		&menu_item.Name,
		&menu_item.Price,
		&menu_item.Category,
	)

	if err != nil {
		return models.MenuItem{}, err
	}

	return menu_item, err
}

func (s *Storage) DeleteMenuItem(ctx context.Context, id int) error {

	res, err := s.pool.Exec(
		ctx,
		"DELETE FROM menu_item WHERE id = $1", id,
	)

	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.ErrNotFound
	}

	return nil
}
