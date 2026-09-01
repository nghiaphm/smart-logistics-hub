package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	inventoryentity "my-web-app.com/smart-logistic-hub/internal/inventory/entity"
	"my-web-app.com/smart-logistic-hub/internal/order/dto"
	"my-web-app.com/smart-logistic-hub/internal/order/entity"
	productentity "my-web-app.com/smart-logistic-hub/internal/product/entity"
)

type OrderRepository interface {
	Create(ctx context.Context, o *entity.Order) error
	GetByID(ctx context.Context, id int64) (*entity.Order, error)
	GetByCode(ctx context.Context, code string) (*entity.Order, error)
	List(ctx context.Context, offset, limit int, workspaceID *int64) ([]entity.Order, error)
	Count(ctx context.Context, workspaceID *int64) (int, error)
	Update(ctx context.Context, id int64, o *entity.Order) error
	Delete(ctx context.Context, id int64) error
	CreateItem(ctx context.Context, item *entity.OrderItem) error
	GetItemsByOrder(ctx context.Context, orderID int64) ([]entity.OrderItem, error)
	DeleteItems(ctx context.Context, orderID int64) error
}

type ProductRepository interface {
	GetByID(ctx context.Context, id int64) (*productentity.Product, error)
}

type InventoryRepository interface {
	GetByProductWarehouse(ctx context.Context, productID, warehouseID int64) (*inventoryentity.Inventory, error)
	Update(ctx context.Context, id int64, inv *inventoryentity.Inventory) error
}

type Service struct {
	repo      OrderRepository
	products  ProductRepository
	inventory InventoryRepository
}

func NewService(repo OrderRepository, products ProductRepository, inventory InventoryRepository) *Service {
	return &Service{repo: repo, products: products, inventory: inventory}
}

func NewServiceWithRepo(repo OrderRepository, products ProductRepository, inventory InventoryRepository) *Service {
	return NewService(repo, products, inventory)
}

func (s *Service) Create(ctx context.Context, req *dto.CreateOrderRequest) (*entity.Order, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("%w: at least one order item is required", apierrors.ErrBadRequest)
	}
	if req.WarehouseID <= 0 {
		return nil, fmt.Errorf("%w: warehouse_id is required", apierrors.ErrBadRequest)
	}

	type stockRef struct {
		productID int64
		quantity  int
	}
	stock := make([]stockRef, 0, len(req.Items))

	for _, it := range req.Items {
		if it.ProductID == nil {
			return nil, fmt.Errorf("%w: product_id is required for each order item", apierrors.ErrBadRequest)
		}
		if _, err := s.products.GetByID(ctx, *it.ProductID); err != nil {
			if errors.Is(err, apierrors.ErrNotFound) {
				return nil, fmt.Errorf("%w: product %d does not exist", apierrors.ErrBadRequest, *it.ProductID)
			}
			return nil, err
		}
		inv, err := s.inventory.GetByProductWarehouse(ctx, *it.ProductID, req.WarehouseID)
		if err != nil {
			if errors.Is(err, apierrors.ErrNotFound) {
				return nil, fmt.Errorf("%w: insufficient inventory for product %d at warehouse %d", apierrors.ErrConflict, *it.ProductID, req.WarehouseID)
			}
			return nil, err
		}
		if inv.AvailableQty < it.Quantity {
			return nil, fmt.Errorf("%w: insufficient inventory for product %d: available %d, required %d",
				apierrors.ErrConflict, *it.ProductID, inv.AvailableQty, it.Quantity)
		}
		stock = append(stock, stockRef{productID: *it.ProductID, quantity: it.Quantity})
	}

	now := time.Now().UTC()
	o := &entity.Order{
		OrderCode:          req.OrderCode,
		SenderWorkspaceID:  req.SenderWorkspaceID,
		SenderName:         req.SenderName,
		SenderPhone:        req.SenderPhone,
		SenderAddress:      req.SenderAddress,
		SenderProvince:     req.SenderProvince,
		SenderDistrict:     req.SenderDistrict,
		SenderWard:         req.SenderWard,
		SenderPostalCode:   req.SenderPostalCode,
		ReceiverName:       req.ReceiverName,
		ReceiverPhone:      req.ReceiverPhone,
		ReceiverAddress:    req.ReceiverAddress,
		ReceiverProvince:   req.ReceiverProvince,
		ReceiverDistrict:   req.ReceiverDistrict,
		ReceiverWard:       req.ReceiverWard,
		ReceiverPostalCode: req.ReceiverPostalCode,
		Status:             req.Status,
		AssignedDriverID:   req.AssignedDriverID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if o.Status == "" {
		o.Status = "PENDING"
	}
	if err := s.repo.Create(ctx, o); err != nil {
		return nil, err
	}
	for i := range req.Items {
		item := &entity.OrderItem{
			OrderID:     o.ID,
			ProductID:   req.Items[i].ProductID,
			ProductName: req.Items[i].ProductName,
			Quantity:    req.Items[i].Quantity,
			WeightGram:  req.Items[i].WeightGram,
		}
		if err := s.repo.CreateItem(ctx, item); err != nil {
			return nil, err
		}
	}

	for _, st := range stock {
		inv, err := s.inventory.GetByProductWarehouse(ctx, st.productID, req.WarehouseID)
		if err != nil {
			return nil, err
		}
		inv.AvailableQty -= st.quantity
		inv.ReservedQty += st.quantity
		if err := s.inventory.Update(ctx, inv.ID, inv); err != nil {
			return nil, err
		}
	}

	return o, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*entity.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetItems(ctx context.Context, orderID int64) ([]entity.OrderItem, error) {
	return s.repo.GetItemsByOrder(ctx, orderID)
}

func (s *Service) List(ctx context.Context, offset, limit int, workspaceID *int64) ([]entity.Order, int, error) {
	orders, err := s.repo.List(ctx, offset, limit, workspaceID)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.Count(ctx, workspaceID)
	if err != nil {
		return nil, 0, err
	}
	return orders, count, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *dto.UpdateOrderRequest) (*entity.Order, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.OrderCode != nil {
		existing.OrderCode = *req.OrderCode
	}
	if req.SenderName != nil {
		existing.SenderName = *req.SenderName
	}
	if req.SenderPhone != nil {
		existing.SenderPhone = *req.SenderPhone
	}
	if req.SenderAddress != nil {
		existing.SenderAddress = *req.SenderAddress
	}
	if req.SenderProvince != nil {
		existing.SenderProvince = *req.SenderProvince
	}
	if req.SenderDistrict != nil {
		existing.SenderDistrict = *req.SenderDistrict
	}
	if req.SenderWard != nil {
		existing.SenderWard = *req.SenderWard
	}
	if req.SenderPostalCode != nil {
		existing.SenderPostalCode = *req.SenderPostalCode
	}
	if req.ReceiverName != nil {
		existing.ReceiverName = *req.ReceiverName
	}
	if req.ReceiverPhone != nil {
		existing.ReceiverPhone = *req.ReceiverPhone
	}
	if req.ReceiverAddress != nil {
		existing.ReceiverAddress = *req.ReceiverAddress
	}
	if req.ReceiverProvince != nil {
		existing.ReceiverProvince = *req.ReceiverProvince
	}
	if req.ReceiverDistrict != nil {
		existing.ReceiverDistrict = *req.ReceiverDistrict
	}
	if req.ReceiverWard != nil {
		existing.ReceiverWard = *req.ReceiverWard
	}
	if req.ReceiverPostalCode != nil {
		existing.ReceiverPostalCode = *req.ReceiverPostalCode
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.AssignedDriverID != nil {
		existing.AssignedDriverID = req.AssignedDriverID
	}
	existing.UpdatedAt = time.Now().UTC()
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
