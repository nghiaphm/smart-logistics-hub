package service

import (
	"context"
	"time"

	"my-web-app.com/smart-logistic-hub/internal/vehicle/dto"
	"my-web-app.com/smart-logistic-hub/internal/vehicle/entity"
)

type VehicleRepository interface {
	Create(ctx context.Context, v *entity.Vehicle) error
	GetByID(ctx context.Context, id int64) (*entity.Vehicle, error)
	List(ctx context.Context, status string, offset, limit int) ([]entity.Vehicle, error)
	Count(ctx context.Context, status string) (int, error)
	Update(ctx context.Context, id int64, v *entity.Vehicle) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	repo VehicleRepository
}

func NewService(repo VehicleRepository) *Service {
	return &Service{repo: repo}
}

func NewServiceWithRepo(repo VehicleRepository) *Service {
	return NewService(repo)
}

func (s *Service) Create(ctx context.Context, req *dto.CreateVehicleRequest) (*entity.Vehicle, error) {
	now := time.Now().UTC()
	v := &entity.Vehicle{
		LicensePlate: req.LicensePlate,
		Type:         req.Type,
		Capacity:     req.Capacity,
		Status:       req.Status,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if v.Status == "" {
		v.Status = "ACTIVE"
	}
	if err := s.repo.Create(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*entity.Vehicle, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, status string, offset, limit int) ([]entity.Vehicle, int, error) {
	vehicles, err := s.repo.List(ctx, status, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.Count(ctx, status)
	if err != nil {
		return nil, 0, err
	}
	return vehicles, count, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *dto.UpdateVehicleRequest) (*entity.Vehicle, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.LicensePlate != nil {
		existing.LicensePlate = *req.LicensePlate
	}
	if req.Type != nil {
		existing.Type = *req.Type
	}
	if req.Capacity != nil {
		existing.Capacity = *req.Capacity
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	existing.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, id, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
