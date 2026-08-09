package service

import (
	"context"
	"database/sql"
	"time"

	"my-web-app.com/smart-logistic-hub/internal/driver/dto"
	"my-web-app.com/smart-logistic-hub/internal/driver/entity"
	"my-web-app.com/smart-logistic-hub/internal/driver/repository"
)

type DriverRepository interface {
	Create(ctx context.Context, d *entity.Driver) error
	GetByID(ctx context.Context, id int64) (*entity.Driver, error)
	GetByCode(ctx context.Context, code string) (*entity.Driver, error)
	List(ctx context.Context, status string, offset, limit int) ([]entity.Driver, error)
	Count(ctx context.Context, status string) (int, error)
	Update(ctx context.Context, id int64, d *entity.Driver) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	repo DriverRepository
}

func NewService(db *sql.DB) *Service {
	return &Service{repo: repository.NewRepository(db)}
}

func NewServiceWithRepo(repo DriverRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req *dto.CreateDriverRequest) (*entity.Driver, error) {
	now := time.Now().UTC()
	d := &entity.Driver{
		DriverCode:   req.DriverCode,
		FullName:     req.FullName,
		Phone:        req.Phone,
		VehicleType:  req.VehicleType,
		LicensePlate: req.LicensePlate,
		Status:       req.Status,
		CurrentLat:   req.CurrentLat,
		CurrentLng:   req.CurrentLng,
		WarehouseID:  req.WarehouseID,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if d.Status == "" {
		d.Status = "AVAILABLE"
	}
	if err := s.repo.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*entity.Driver, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, status string, offset, limit int) ([]entity.Driver, int, error) {
	drivers, err := s.repo.List(ctx, status, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.Count(ctx, status)
	if err != nil {
		return nil, 0, err
	}
	return drivers, count, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *dto.UpdateDriverRequest) (*entity.Driver, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.DriverCode != nil {
		existing.DriverCode = *req.DriverCode
	}
	if req.FullName != nil {
		existing.FullName = *req.FullName
	}
	if req.Phone != nil {
		existing.Phone = *req.Phone
	}
	if req.VehicleType != nil {
		existing.VehicleType = *req.VehicleType
	}
	if req.LicensePlate != nil {
		existing.LicensePlate = *req.LicensePlate
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.CurrentLat != nil {
		existing.CurrentLat = *req.CurrentLat
	}
	if req.CurrentLng != nil {
		existing.CurrentLng = *req.CurrentLng
	}
	if req.WarehouseID != nil {
		existing.WarehouseID = req.WarehouseID
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
