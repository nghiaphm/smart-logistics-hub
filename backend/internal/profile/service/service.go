package service

import (
	"context"

	"my-web-app.com/smart-logistic-hub/internal/profile/dto"
	"my-web-app.com/smart-logistic-hub/internal/profile/entity"
)

type ProfileRepository interface {
	GetByKeycloakUserID(ctx context.Context, keycloakUserID string) (*entity.Profile, error)
	Create(ctx context.Context, p *entity.Profile) error
	Update(ctx context.Context, p *entity.Profile) error
}

type Service struct {
	repo ProfileRepository
}

func NewService(repo ProfileRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, keycloakUserID string) (*entity.Profile, error) {
	return s.repo.GetByKeycloakUserID(ctx, keycloakUserID)
}

func (s *Service) Create(ctx context.Context, keycloakUserID string, req *dto.CreateProfileRequest) (*entity.Profile, error) {
	p := &entity.Profile{
		KeycloakUserID: keycloakUserID,
		Name:           req.Name,
		Phone:          req.Phone,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Update(ctx context.Context, keycloakUserID string, req *dto.UpdateProfileRequest) (*entity.Profile, error) {
	p, err := s.Get(ctx, keycloakUserID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Phone != nil {
		p.Phone = *req.Phone
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
