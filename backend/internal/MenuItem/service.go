package menuitem

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

func (s *Service) Create(ctx context.Context, menu models.MenuItem) (int, error) {
	if menu.Name == "" || menu.Price == 0 || menu.Category == "" {
		return 0, myerrors.ErrBadRequest
	}

	id, err := s.storage.Create(ctx, menu)
	if err != nil {
		return 0, err
	}

	return id, nil
}
