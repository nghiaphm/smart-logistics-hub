package service

import (
	"context"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/profile/dto"
	"my-web-app.com/smart-logistic-hub/internal/profile/entity"
)

type ProfileRepository interface {
	GetByUserSub(ctx context.Context, userSub string) (*entity.Profile, error)
	Create(ctx context.Context, p *entity.Profile) error
	Update(ctx context.Context, p *entity.Profile) error
}

type Service struct {
	repo ProfileRepository
}

func NewService(repo ProfileRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, userSub string) (*entity.Profile, error) {
	p, err := s.repo.GetByUserSub(ctx, userSub)
	if err == apierrors.ErrNotFound {
		p = &entity.Profile{UserSub: userSub}
		if err := s.repo.Create(ctx, p); err != nil {
			return nil, err
		}
		return p, nil
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Update(ctx context.Context, userSub string, req *dto.UpdateProfileRequest) (*entity.Profile, error) {
	p, err := s.Get(ctx, userSub)
	if err != nil {
		return nil, err
	}

	if req.DisplayName != nil {
		p.DisplayName = *req.DisplayName
	}
	if req.Phone != nil {
		p.Phone = *req.Phone
	}
	if req.AvatarURL != nil {
		p.AvatarURL = *req.AvatarURL
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
