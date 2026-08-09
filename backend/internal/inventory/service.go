package inventory

import (
	"context"
	"fmt"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
)

type InventoryRepository interface {
	Create(ctx context.Context, inv *Inventory) error
	GetByID(ctx context.Context, id int64) (*Inventory, error)
	GetByProductWarehouse(ctx context.Context, productID, warehouseID int64) (*Inventory, error)
	List(ctx context.Context, offset, limit int) ([]Inventory, error)
	Count(ctx context.Context) (int, error)
	Update(ctx context.Context, id int64, inv *Inventory) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	Repo InventoryRepository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func NewServiceWithRepo(repo InventoryRepository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) Create(ctx context.Context, req *CreateInventoryRequest) (*Inventory, error) {
	if req.ProductID <= 0 {
		return nil, fmt.Errorf("%w: product_id is required", apierrors.ErrBadRequest)
	}
	if req.WarehouseID <= 0 {
		return nil, fmt.Errorf("%w: warehouse_id is required", apierrors.ErrBadRequest)
	}
	if req.AvailableQty < 0 || req.ReservedQty < 0 || req.DamagedQty < 0 || req.HoldQty < 0 {
		return nil, fmt.Errorf("%w: quantities must not be negative", apierrors.ErrBadRequest)
	}

	existing, err := s.Repo.GetByProductWarehouse(ctx, req.ProductID, req.WarehouseID)
	if err != nil && err != apierrors.ErrNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: inventory already exists for product %d at warehouse %d", apierrors.ErrConflict, req.ProductID, req.WarehouseID)
	}

	inv := &Inventory{
		ProductID:    req.ProductID,
		WarehouseID:  req.WarehouseID,
		AvailableQty: req.AvailableQty,
		ReservedQty:  req.ReservedQty,
		DamagedQty:   req.DamagedQty,
		HoldQty:      req.HoldQty,
		UpdatedBy:    req.UpdatedBy,
	}
	if err := s.Repo.Create(ctx, inv); err != nil {
		return nil, err
	}
	return inv, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*Inventory, error) {
	inv, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func (s *Service) List(ctx context.Context, offset, limit int) ([]Inventory, int, error) {
	invList, err := s.Repo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.Repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return invList, count, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpdateInventoryRequest) (*Inventory, error) {
	inv, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.AvailableQty != nil {
		if *req.AvailableQty < 0 {
			return nil, fmt.Errorf("%w: available_qty must not be negative", apierrors.ErrBadRequest)
		}
		inv.AvailableQty = *req.AvailableQty
	}
	if req.ReservedQty != nil {
		if *req.ReservedQty < 0 {
			return nil, fmt.Errorf("%w: reserved_qty must not be negative", apierrors.ErrBadRequest)
		}
		inv.ReservedQty = *req.ReservedQty
	}
	if req.DamagedQty != nil {
		if *req.DamagedQty < 0 {
			return nil, fmt.Errorf("%w: damaged_qty must not be negative", apierrors.ErrBadRequest)
		}
		inv.DamagedQty = *req.DamagedQty
	}
	if req.HoldQty != nil {
		if *req.HoldQty < 0 {
			return nil, fmt.Errorf("%w: hold_qty must not be negative", apierrors.ErrBadRequest)
		}
		inv.HoldQty = *req.HoldQty
	}
	if req.UpdatedBy != "" {
		inv.UpdatedBy = req.UpdatedBy
	}

	if err := s.Repo.Update(ctx, id, inv); err != nil {
		return nil, err
	}
	return inv, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.Repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
