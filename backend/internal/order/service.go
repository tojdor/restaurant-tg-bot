package order

import (
	"backend/internal/models"
	myerrors "backend/internal/my_errors"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type Service struct {
	storage *Storage
}

func NewService(storage *Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) Create(ctx context.Context, order models.Order) (int, error) {

	if order.WaiterID <= 0 || order.TableNumber <= 0 {
		return 0, myerrors.ErrBadRequest
	}

	id, err := s.storage.Create(ctx, order)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) GetAll(ctx context.Context) ([]models.Order, error) {
	orders, err := s.storage.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return []models.Order{}, nil
	}

	return orders, nil
}

func (s *Service) GetByTableNumber(ctx context.Context, number int) (models.Order, error) {

	if number <= 0 {
		return models.Order{}, myerrors.ErrBadRequest
	}

	order, err := s.storage.GetByTableNumber(ctx, number)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Order{}, myerrors.ErrNotFound
		}
		return models.Order{}, err
	}

	return order, nil
}

func (s *Service) GetByWaiterId(ctx context.Context, id int) ([]models.Order, error) {

	if id <= 0 {
		return nil, myerrors.ErrBadRequest
	}

	res, err := s.storage.GetByWaiterId(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Service) IsOwnedByWaiter(ctx context.Context, orderID, waiterID int) (bool, error) {
	if orderID <= 0 || waiterID <= 0 {
		return false, myerrors.ErrBadRequest
	}
	return s.storage.IsOwnedByWaiter(ctx, orderID, waiterID)
}

func (s *Service) Update(ctx context.Context, order OrderRequestBody) error {

	if order.ID <= 0 || order.TableNumber == 0 || order.WaiterID == 0 {
		return myerrors.ErrBadRequest
	}

	if err := s.storage.Update(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id int) error {

	if id <= 0 {
		return myerrors.ErrBadRequest
	}

	if err := s.storage.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
