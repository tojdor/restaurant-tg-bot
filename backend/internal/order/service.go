package order

import (
	"backend/internal/models"
	myerrors "backend/internal/my_errors"
	"context"
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

	if order.WaiterID != 0 && order.TableNumber != 0 {
		return 0, myerrors.ErrBadRequest
	}

	id, err := s.storage.Create(ctx, order)
	if err != nil {
		return 0, nil
	}

	return id, nil
}
