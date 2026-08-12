package table

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

func (s *Service) Create(ctx context.Context, table models.Table) (int, error) {
	if table.Number <= 0 {
		return 0, myerrors.ErrBadRequest
	}

	switch table.Status {
	case "free", "occupied":
	default:
		return 0, myerrors.ErrBadRequest
	}

	number, err := s.storage.Create(ctx, table)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func (s *Service) GetAll(ctx context.Context) ([]models.Table, error) {
	return s.storage.GetAll(ctx)
}

func (s *Service) Delete(ctx context.Context, number int) error {
	if number <= 0 {
		return myerrors.ErrBadRequest
	}

	return s.storage.Delete(ctx, number)

}
