package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/inbound/dto"
	"my-web-app.com/smart-logistic-hub/internal/inbound/entity"
	inventoryentity "my-web-app.com/smart-logistic-hub/internal/inventory/entity"
)

type InboundRepository interface {
	Create(ctx context.Context, in *entity.Inbound) error
	GetByID(ctx context.Context, id int64) (*entity.Inbound, error)
	GetByCode(ctx context.Context, code string) (*entity.Inbound, error)
	List(ctx context.Context, offset, limit int) ([]entity.Inbound, error)
	Count(ctx context.Context) (int, error)
	Update(ctx context.Context, id int64, in *entity.Inbound) error
	Delete(ctx context.Context, id int64) error
	CreateItems(ctx context.Context, inboundID int64, items []entity.InboundItem) error
	GetItems(ctx context.Context, inboundID int64) ([]entity.InboundItem, error)
	DeleteItems(ctx context.Context, inboundID int64) error
}

type InventoryRepository interface {
	GetByProductWarehouse(ctx context.Context, productID, warehouseID int64) (*inventoryentity.Inventory, error)
	Create(ctx context.Context, inv *inventoryentity.Inventory) error
	Update(ctx context.Context, id int64, inv *inventoryentity.Inventory) error
}

type Service struct {
	repo      InboundRepository
	inventory InventoryRepository
}

func NewService(repo InboundRepository, inventory InventoryRepository) *Service {
	return &Service{repo: repo, inventory: inventory}
}

func NewServiceWithRepo(repo InboundRepository, inventory InventoryRepository) *Service {
	return NewService(repo, inventory)
}

func (s *Service) Create(ctx context.Context, req *dto.CreateInboundRequest) (*entity.Inbound, error) {
	if req.ReceiptCode == "" {
		return nil, fmt.Errorf("%w: receipt_code is required", apierrors.ErrBadRequest)
	}
	if req.SupplierName == "" {
		return nil, fmt.Errorf("%w: supplier_name is required", apierrors.ErrBadRequest)
	}
	if req.WarehouseID <= 0 {
		return nil, fmt.Errorf("%w: warehouse_id is required", apierrors.ErrBadRequest)
	}
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("%w: at least one inbound item is required", apierrors.ErrBadRequest)
	}

	existing, err := s.repo.GetByCode(ctx, req.ReceiptCode)
	if err != nil && err != apierrors.ErrNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: inbound with receipt_code %s already exists", apierrors.ErrConflict, req.ReceiptCode)
	}

	items, err := buildItems(req.Items)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	in := &entity.Inbound{
		ReceiptCode:  req.ReceiptCode,
		SupplierName: req.SupplierName,
		WarehouseID:  req.WarehouseID,
		Status:       req.Status,
		CreatedAt:    now,
		UpdatedAt:    now,
		CreatedBy:    req.CreatedBy,
	}
	if in.Status == "" {
		in.Status = "PENDING"
	}
	if err := s.repo.Create(ctx, in); err != nil {
		return nil, err
	}
	if err := s.repo.CreateItems(ctx, in.ID, items); err != nil {
		return nil, err
	}

	if in.Status == "COMPLETED" {
		in.CompletedAt = &now
		if err := s.addStock(ctx, in.WarehouseID, items); err != nil {
			return nil, err
		}
		if err := s.repo.Update(ctx, in.ID, in); err != nil {
			return nil, err
		}
	}
	return in, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*entity.Inbound, []entity.InboundItem, error) {
	in, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	items, err := s.repo.GetItems(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return in, items, nil
}

func (s *Service) List(ctx context.Context, offset, limit int) ([]entity.Inbound, int, error) {
	inbounds, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return inbounds, count, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *dto.UpdateInboundRequest) (*entity.Inbound, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	oldStatus := existing.Status
	if req.Status != nil {
		existing.Status = *req.Status
	}

	var items []entity.InboundItem
	if req.Items != nil {
		items, err = buildItems(*req.Items)
		if err != nil {
			return nil, err
		}
		if err := s.repo.DeleteItems(ctx, id); err != nil {
			return nil, err
		}
		if err := s.repo.CreateItems(ctx, id, items); err != nil {
			return nil, err
		}
	} else {
		items, err = s.repo.GetItems(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	if existing.Status == "COMPLETED" && oldStatus != "COMPLETED" {
		now := time.Now().UTC()
		existing.CompletedAt = &now
		if err := s.addStock(ctx, existing.WarehouseID, items); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Update(ctx, id, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.DeleteItems(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) addStock(ctx context.Context, warehouseID int64, items []entity.InboundItem) error {
	for _, it := range items {
		if it.ReceivedQty <= 0 {
			continue
		}
		inv, err := s.inventory.GetByProductWarehouse(ctx, it.ProductID, warehouseID)
		if err != nil {
			if !errors.Is(err, apierrors.ErrNotFound) {
				return err
			}
			inv = &inventoryentity.Inventory{
				ProductID:    it.ProductID,
				WarehouseID:  warehouseID,
				AvailableQty: it.ReceivedQty,
			}
			if err := s.inventory.Create(ctx, inv); err != nil {
				return err
			}
			continue
		}
		inv.AvailableQty += it.ReceivedQty
		if err := s.inventory.Update(ctx, inv.ID, inv); err != nil {
			return err
		}
	}
	return nil
}

func buildItems(req []dto.InboundItemRequest) ([]entity.InboundItem, error) {
	items := make([]entity.InboundItem, 0, len(req))
	for _, it := range req {
		if it.ProductID <= 0 {
			return nil, fmt.Errorf("%w: product_id is required for each inbound item", apierrors.ErrBadRequest)
		}
		if it.ExpectedQty < 0 || it.ReceivedQty < 0 || it.RejectedQty < 0 || it.QcPassed < 0 {
			return nil, fmt.Errorf("%w: quantities must not be negative", apierrors.ErrBadRequest)
		}
		items = append(items, entity.InboundItem{
			ProductID:   it.ProductID,
			ExpectedQty: it.ExpectedQty,
			ReceivedQty: it.ReceivedQty,
			RejectedQty: it.RejectedQty,
			QcPassed:    it.QcPassed,
		})
	}
	return items, nil
}
