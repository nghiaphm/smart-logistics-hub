package service

import (
	"context"
	"fmt"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/warehouse/dto"
	"my-web-app.com/smart-logistic-hub/internal/warehouse/entity"
)

type WarehouseRepository interface {
	Create(ctx context.Context, w *entity.Warehouse) error
	GetByID(ctx context.Context, id int64) (*entity.Warehouse, error)
	GetByCode(ctx context.Context, code string) (*entity.Warehouse, error)
	List(ctx context.Context, offset, limit int) ([]entity.Warehouse, error)
	Count(ctx context.Context) (int, error)
	Update(ctx context.Context, id int64, w *entity.Warehouse) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	repo WarehouseRepository
}

func NewService(repo WarehouseRepository) *Service {
	return &Service{repo: repo}
}

func NewServiceWithRepo(repo WarehouseRepository) *Service {
	return NewService(repo)
}

func (s *Service) Create(ctx context.Context, req *dto.CreateWarehouseRequest) (*entity.Warehouse, error) {
	if req.WarehouseCode == "" {
		return nil, fmt.Errorf("%w: warehouse_code is required", apierrors.ErrBadRequest)
	}
	if req.Name == "" {
		return nil, fmt.Errorf("%w: name is required", apierrors.ErrBadRequest)
	}
	if req.Address == "" {
		return nil, fmt.Errorf("%w: address is required", apierrors.ErrBadRequest)
	}

	existing, err := s.repo.GetByCode(ctx, req.WarehouseCode)
	if err != nil && err != apierrors.ErrNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: warehouse with code %s already exists", apierrors.ErrConflict, req.WarehouseCode)
	}

	w := &entity.Warehouse{
		WarehouseCode: req.WarehouseCode,
		Name:          req.Name,
		Address:       req.Address,
		IsActive:      true,
	}
	if req.Location != nil {
		if req.Location.Lat != nil {
			w.Lat = *req.Location.Lat
		}
		if req.Location.Lng != nil {
			w.Lng = *req.Location.Lng
		}
	}
	if req.ContactPhone != nil {
		w.ContactPhone = *req.ContactPhone
	}
	if req.ManagerName != nil {
		w.ManagerName = *req.ManagerName
	}

	if err := s.repo.Create(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*entity.Warehouse, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, offset, limit int) ([]entity.Warehouse, int, error) {
	warehouses, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return warehouses, count, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *dto.UpdateWarehouseRequest) (*entity.Warehouse, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Address != nil {
		existing.Address = *req.Address
	}
	if req.Location != nil {
		if req.Location.Lat != nil {
			existing.Lat = *req.Location.Lat
		}
		if req.Location.Lng != nil {
			existing.Lng = *req.Location.Lng
		}
	}
	if req.ContactPhone != nil {
		existing.ContactPhone = *req.ContactPhone
	}
	if req.ManagerName != nil {
		existing.ManagerName = *req.ManagerName
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
