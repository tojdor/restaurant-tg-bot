package user

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

func (s *Service) Create(ctx context.Context, user models.User) (int, error) {
	if user.TelegramUserID <= 0 || user.Nickname == "" || user.PhoneNumber == "" {
		return 0, myerrors.ErrBadRequest
	}

	if user.Role == "admin" || user.Role == "waiter" || user.Role == "kitchen" {
		id, err := s.storage.Create(ctx, user)
		if err != nil {
			return 0, err
		}
		return id, nil
	}

	return 0, myerrors.ErrBadRequest
}

func (s *Service) GetUsers(ctx context.Context) ([]models.User, error) {
	return s.storage.GetUsers(ctx)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return myerrors.ErrBadRequest
	}

	return s.storage.Delete(ctx, id)
}
