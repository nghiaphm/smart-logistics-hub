package service

import (
	"context"
	"fmt"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/user/dto"
	"my-web-app.com/smart-logistic-hub/internal/user/entity"
)

type UserRepository interface {
	Create(ctx context.Context, u *entity.User) error
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	List(ctx context.Context, offset, limit int) ([]entity.User, error)
	Count(ctx context.Context) (int, error)
	Update(ctx context.Context, id int64, u *entity.User) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req *dto.CreateUserRequest) (*entity.User, error) {
	if req.Username == "" {
		return nil, fmt.Errorf("%w: username is required", apierrors.ErrBadRequest)
	}

	existing, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil && err != apierrors.ErrNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: user with username %s already exists", apierrors.ErrConflict, req.Username)
	}

	u := &entity.User{
		Username: req.Username,
		Role:     "user",
		IsActive: true,
	}
	if req.KeycloakSub != nil {
		u.KeycloakSub = *req.KeycloakSub
	}
	if req.FullName != nil {
		u.FullName = *req.FullName
	}
	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.Phone != nil {
		u.Phone = *req.Phone
	}
	if req.Role != nil {
		u.Role = *req.Role
	}
	if req.IsActive != nil {
		u.IsActive = *req.IsActive
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*entity.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, offset, limit int) ([]entity.User, int, error) {
	users, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *dto.UpdateUserRequest) (*entity.User, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.KeycloakSub != nil {
		existing.KeycloakSub = *req.KeycloakSub
	}
	if req.FullName != nil {
		existing.FullName = *req.FullName
	}
	if req.Email != nil {
		existing.Email = *req.Email
	}
	if req.Phone != nil {
		existing.Phone = *req.Phone
	}
	if req.Role != nil {
		existing.Role = *req.Role
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, id, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
