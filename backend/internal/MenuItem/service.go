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

	_, err := s.storage.GetByName(ctx, menu.Name)
	if err == nil {
		return 0, myerrors.ErrAlreadyExists
	}

	id, err := s.storage.Create(ctx, menu)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) GetAll(ctx context.Context) ([]models.MenuItem, error) {
	return s.storage.GetAll(ctx)
}

func (s *Service) Update(ctx context.Context, menu models.MenuItem) error {

	if menu.Name == "" || menu.Price == 0 || menu.Category == "" {
		return myerrors.ErrBadRequest
	}

	_, err := s.storage.GetByID(ctx, menu.ID)
	if err == nil {
		return myerrors.ErrAlreadyExists
	}

	if err = s.storage.Update(ctx, menu); err != nil {
		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return myerrors.ErrNotFound
	}

	return s.storage.Delete(ctx, id)
}
