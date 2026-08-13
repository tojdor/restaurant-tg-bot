package orderitem

import (
	"backend/internal/models"
	"context"
	"errors"
)

type Service struct {
	storage *Storage
}

func NewService(storage *Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) AddItem(
	ctx context.Context,
	item models.OrderItem,
) error {
	if item.OrderID <= 0 {
		return errors.New("invalid order id")
	}

	if item.MenuItemID <= 0 {
		return errors.New("invalid menu item id")
	}

	if item.Count <= 0 {
		return errors.New("count must be greater than zero")
	}

	return s.storage.AddItem(ctx, item)
}

func (s *Service) GetByOrderID(
	ctx context.Context,
	orderID int,
) ([]models.OrderItem, error) {
	if orderID <= 0 {
		return nil, errors.New("invalid order id")
	}

	return s.storage.GetByOrderID(ctx, orderID)
}

func (s *Service) IsOwnedByWaiter(ctx context.Context, orderID, waiterID int) (bool, error) {
	if orderID <= 0 || waiterID <= 0 {
		return false, errors.New("invalid order or waiter id")
	}
	return s.storage.IsOwnedByWaiter(ctx, orderID, waiterID)
}

func (s *Service) SetReady(
	ctx context.Context,
	orderID int,
	menuItemID int,
	ready bool,
) error {
	if orderID <= 0 {
		return errors.New("invalid order id")
	}

	if menuItemID <= 0 {
		return errors.New("invalid menu item id")
	}

	return s.storage.SetReady(
		ctx,
		orderID,
		menuItemID,
		ready,
	)
}

func (s *Service) UpdateCount(
	ctx context.Context,
	orderID int,
	menuItemID int,
	count int,
) error {
	if orderID <= 0 {
		return errors.New("invalid order id")
	}

	if menuItemID <= 0 {
		return errors.New("invalid menu item id")
	}

	if count <= 0 {
		return errors.New("count must be greater than zero")
	}

	return s.storage.UpdateCount(
		ctx,
		orderID,
		menuItemID,
		count,
	)
}

func (s *Service) DeleteItem(
	ctx context.Context,
	orderID int,
	menuItemID int,
) error {
	if orderID <= 0 {
		return errors.New("invalid order id")
	}

	if menuItemID <= 0 {
		return errors.New("invalid menu item id")
	}

	return s.storage.DeleteItem(
		ctx,
		orderID,
		menuItemID,
	)
}
