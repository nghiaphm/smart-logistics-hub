package service_test

import (
	"context"
	"errors"
	"testing"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/inbound/dto"
	"my-web-app.com/smart-logistic-hub/internal/inbound/entity"
	"my-web-app.com/smart-logistic-hub/internal/inbound/service"
	inventoryentity "my-web-app.com/smart-logistic-hub/internal/inventory/entity"
)

type mockInboundRepo struct {
	createFn      func(ctx context.Context, in *entity.Inbound) error
	getByIDFn     func(ctx context.Context, id int64) (*entity.Inbound, error)
	getByCodeFn   func(ctx context.Context, code string) (*entity.Inbound, error)
	listFn        func(ctx context.Context, offset, limit int) ([]entity.Inbound, error)
	countFn       func(ctx context.Context) (int, error)
	updateFn      func(ctx context.Context, id int64, in *entity.Inbound) error
	deleteFn      func(ctx context.Context, id int64) error
	createItemsFn func(ctx context.Context, inboundID int64, items []entity.InboundItem) error
	getItemsFn    func(ctx context.Context, inboundID int64) ([]entity.InboundItem, error)
	deleteItemsFn func(ctx context.Context, inboundID int64) error
}

func (m *mockInboundRepo) Create(ctx context.Context, in *entity.Inbound) error {
	if m.createFn != nil {
		return m.createFn(ctx, in)
	}
	return nil
}

func (m *mockInboundRepo) GetByID(ctx context.Context, id int64) (*entity.Inbound, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockInboundRepo) GetByCode(ctx context.Context, code string) (*entity.Inbound, error) {
	if m.getByCodeFn != nil {
		return m.getByCodeFn(ctx, code)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockInboundRepo) List(ctx context.Context, offset, limit int) ([]entity.Inbound, error) {
	if m.listFn != nil {
		return m.listFn(ctx, offset, limit)
	}
	return []entity.Inbound{}, nil
}

func (m *mockInboundRepo) Count(ctx context.Context) (int, error) {
	if m.countFn != nil {
		return m.countFn(ctx)
	}
	return 0, nil
}

func (m *mockInboundRepo) Update(ctx context.Context, id int64, in *entity.Inbound) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, in)
	}
	return nil
}

func (m *mockInboundRepo) Delete(ctx context.Context, id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockInboundRepo) CreateItems(ctx context.Context, inboundID int64, items []entity.InboundItem) error {
	if m.createItemsFn != nil {
		return m.createItemsFn(ctx, inboundID, items)
	}
	return nil
}

func (m *mockInboundRepo) GetItems(ctx context.Context, inboundID int64) ([]entity.InboundItem, error) {
	if m.getItemsFn != nil {
		return m.getItemsFn(ctx, inboundID)
	}
	return []entity.InboundItem{}, nil
}

func (m *mockInboundRepo) DeleteItems(ctx context.Context, inboundID int64) error {
	if m.deleteItemsFn != nil {
		return m.deleteItemsFn(ctx, inboundID)
	}
	return nil
}

type mockInventoryRepo struct {
	getByPWFn func(ctx context.Context, productID, warehouseID int64) (*inventoryentity.Inventory, error)
	createFn  func(ctx context.Context, inv *inventoryentity.Inventory) error
	updateFn  func(ctx context.Context, id int64, inv *inventoryentity.Inventory) error
}

func (m *mockInventoryRepo) GetByProductWarehouse(ctx context.Context, productID, warehouseID int64) (*inventoryentity.Inventory, error) {
	if m.getByPWFn != nil {
		return m.getByPWFn(ctx, productID, warehouseID)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockInventoryRepo) Create(ctx context.Context, inv *inventoryentity.Inventory) error {
	if m.createFn != nil {
		return m.createFn(ctx, inv)
	}
	return nil
}

func (m *mockInventoryRepo) Update(ctx context.Context, id int64, inv *inventoryentity.Inventory) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, inv)
	}
	return nil
}

func validCreateRequest() *dto.CreateInboundRequest {
	return &dto.CreateInboundRequest{
		ReceiptCode:  "RCP-001",
		SupplierName: "Supplier A",
		WarehouseID:  1,
		Items: []dto.InboundItemRequest{
			{ProductID: 100, ExpectedQty: 10},
		},
	}
}

func TestInboundCreateValidDefaultsStatus(t *testing.T) {
	var saved *entity.Inbound
	var createdItems []entity.InboundItem
	inboundRepo := &mockInboundRepo{
		createFn: func(ctx context.Context, in *entity.Inbound) error {
			in.ID = 1
			saved = in
			return nil
		},
		createItemsFn: func(ctx context.Context, inboundID int64, items []entity.InboundItem) error {
			createdItems = items
			return nil
		},
	}
	svc := service.NewServiceWithRepo(inboundRepo, &mockInventoryRepo{})

	in, err := svc.Create(context.Background(), validCreateRequest())
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if in.ID != 1 {
		t.Errorf("Create() ID = %d, want 1", in.ID)
	}
	if saved.Status != "PENDING" {
		t.Errorf("Create() status = %q, want %q", saved.Status, "PENDING")
	}
	if len(createdItems) != 1 {
		t.Errorf("Create() created %d items, want 1", len(createdItems))
	}
}

func TestInboundCreateRequiresReceiptCode(t *testing.T) {
	req := validCreateRequest()
	req.ReceiptCode = ""
	svc := service.NewServiceWithRepo(&mockInboundRepo{}, &mockInventoryRepo{})
	_, err := svc.Create(context.Background(), req)
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestInboundCreateRequiresWarehouse(t *testing.T) {
	req := validCreateRequest()
	req.WarehouseID = 0
	svc := service.NewServiceWithRepo(&mockInboundRepo{}, &mockInventoryRepo{})
	_, err := svc.Create(context.Background(), req)
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestInboundCreateRequiresItems(t *testing.T) {
	req := validCreateRequest()
	req.Items = nil
	svc := service.NewServiceWithRepo(&mockInboundRepo{}, &mockInventoryRepo{})
	_, err := svc.Create(context.Background(), req)
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestInboundCreateRejectsDuplicateReceiptCode(t *testing.T) {
	inboundRepo := &mockInboundRepo{
		getByCodeFn: func(ctx context.Context, code string) (*entity.Inbound, error) {
			return &entity.Inbound{ID: 1, ReceiptCode: code}, nil
		},
	}
	svc := service.NewServiceWithRepo(inboundRepo, &mockInventoryRepo{})
	_, err := svc.Create(context.Background(), validCreateRequest())
	if !errors.Is(err, apierrors.ErrConflict) {
		t.Errorf("Create() error = %v, want ErrConflict", err)
	}
}

func TestInboundCreateRejectsNegativeQty(t *testing.T) {
	req := validCreateRequest()
	req.Items[0].ReceivedQty = -1
	svc := service.NewServiceWithRepo(&mockInboundRepo{}, &mockInventoryRepo{})
	_, err := svc.Create(context.Background(), req)
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestInboundGetReturnsItems(t *testing.T) {
	items := []entity.InboundItem{{ID: 1, InboundID: 1, ProductID: 100}}
	inboundRepo := &mockInboundRepo{
		getByIDFn: func(ctx context.Context, id int64) (*entity.Inbound, error) {
			return &entity.Inbound{ID: 1, ReceiptCode: "RCP-001"}, nil
		},
		getItemsFn: func(ctx context.Context, inboundID int64) ([]entity.InboundItem, error) {
			return items, nil
		},
	}
	svc := service.NewServiceWithRepo(inboundRepo, &mockInventoryRepo{})

	_, got, err := svc.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if len(got) != 1 {
		t.Errorf("Get() items len = %d, want 1", len(got))
	}
}

func TestInboundGetReturnsErrNotFound(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockInboundRepo{}, &mockInventoryRepo{})
	_, _, err := svc.Get(context.Background(), 99)
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestInboundListReturnsItemsAndCount(t *testing.T) {
	items := []entity.Inbound{{ID: 1, ReceiptCode: "RCP-001"}}
	inboundRepo := &mockInboundRepo{
		listFn:  func(ctx context.Context, offset, limit int) ([]entity.Inbound, error) { return items, nil },
		countFn: func(ctx context.Context) (int, error) { return 1, nil },
	}
	svc := service.NewServiceWithRepo(inboundRepo, &mockInventoryRepo{})

	got, total, err := svc.List(context.Background(), 0, 20)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if total != 1 || len(got) != 1 {
		t.Errorf("List() = (%d items, %d total), want (1, 1)", len(got), total)
	}
}

func TestInboundUpdateCompletesAndAddsStock(t *testing.T) {
	var updatedInv *inventoryentity.Inventory
	inboundRepo := &mockInboundRepo{
		getByIDFn: func(ctx context.Context, id int64) (*entity.Inbound, error) {
			return &entity.Inbound{ID: 1, ReceiptCode: "RCP-001", WarehouseID: 1, Status: "PENDING"}, nil
		},
		getItemsFn: func(ctx context.Context, inboundID int64) ([]entity.InboundItem, error) {
			return []entity.InboundItem{{ID: 1, ProductID: 100, ReceivedQty: 5}}, nil
		},
	}
	inventoryRepo := &mockInventoryRepo{
		getByPWFn: func(ctx context.Context, productID, warehouseID int64) (*inventoryentity.Inventory, error) {
			return &inventoryentity.Inventory{ID: 1, ProductID: productID, WarehouseID: warehouseID, AvailableQty: 10}, nil
		},
		updateFn: func(ctx context.Context, id int64, inv *inventoryentity.Inventory) error {
			updatedInv = inv
			return nil
		},
	}
	svc := service.NewServiceWithRepo(inboundRepo, inventoryRepo)

	newStatus := "COMPLETED"
	in, err := svc.Update(context.Background(), 1, &dto.UpdateInboundRequest{Status: &newStatus})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if in.Status != "COMPLETED" {
		t.Errorf("Update() status = %q, want COMPLETED", in.Status)
	}
	if in.CompletedAt == nil {
		t.Error("Update() CompletedAt not set")
	}
	if updatedInv == nil {
		t.Fatal("Update() did not add stock to inventory")
	}
	if updatedInv.AvailableQty != 15 {
		t.Errorf("Update() available_qty = %d, want 15 (10 + 5 received)", updatedInv.AvailableQty)
	}
}

func TestInboundUpdateCompletesCreatesInventoryRecordWhenMissing(t *testing.T) {
	var createdInv *inventoryentity.Inventory
	inboundRepo := &mockInboundRepo{
		getByIDFn: func(ctx context.Context, id int64) (*entity.Inbound, error) {
			return &entity.Inbound{ID: 1, ReceiptCode: "RCP-001", WarehouseID: 1, Status: "PENDING"}, nil
		},
		getItemsFn: func(ctx context.Context, inboundID int64) ([]entity.InboundItem, error) {
			return []entity.InboundItem{{ID: 1, ProductID: 100, ReceivedQty: 3}}, nil
		},
	}
	inventoryRepo := &mockInventoryRepo{
		createFn: func(ctx context.Context, inv *inventoryentity.Inventory) error {
			createdInv = inv
			return nil
		},
	}
	svc := service.NewServiceWithRepo(inboundRepo, inventoryRepo)

	newStatus := "COMPLETED"
	_, err := svc.Update(context.Background(), 1, &dto.UpdateInboundRequest{Status: &newStatus})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if createdInv == nil {
		t.Fatal("Update() did not create inventory record")
	}
	if createdInv.ProductID != 100 || createdInv.AvailableQty != 3 {
		t.Errorf("Update() created inventory = (product %d, available %d), want (100, 3)", createdInv.ProductID, createdInv.AvailableQty)
	}
}

func TestInboundUpdateDoesNotDoubleAddStock(t *testing.T) {
	var updateCalls int
	inboundRepo := &mockInboundRepo{
		getByIDFn: func(ctx context.Context, id int64) (*entity.Inbound, error) {
			return &entity.Inbound{ID: 1, ReceiptCode: "RCP-001", WarehouseID: 1, Status: "COMPLETED"}, nil
		},
		getItemsFn: func(ctx context.Context, inboundID int64) ([]entity.InboundItem, error) {
			return []entity.InboundItem{{ID: 1, ProductID: 100, ReceivedQty: 5}}, nil
		},
	}
	inventoryRepo := &mockInventoryRepo{
		getByPWFn: func(ctx context.Context, productID, warehouseID int64) (*inventoryentity.Inventory, error) {
			return &inventoryentity.Inventory{ID: 1, ProductID: productID, WarehouseID: warehouseID, AvailableQty: 10}, nil
		},
		updateFn: func(ctx context.Context, id int64, inv *inventoryentity.Inventory) error {
			updateCalls++
			return nil
		},
	}
	svc := service.NewServiceWithRepo(inboundRepo, inventoryRepo)

	newStatus := "RECEIVING"
	if _, err := svc.Update(context.Background(), 1, &dto.UpdateInboundRequest{Status: &newStatus}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updateCalls != 0 {
		t.Errorf("Update() inventory updated %d times for non-completion, want 0", updateCalls)
	}
}

func TestInboundUpdateReturnsErrNotFound(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockInboundRepo{}, &mockInventoryRepo{})
	_, err := svc.Update(context.Background(), 99, &dto.UpdateInboundRequest{})
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestInboundDeleteReturnsErrNotFound(t *testing.T) {
	inboundRepo := &mockInboundRepo{deleteFn: func(ctx context.Context, id int64) error { return apierrors.ErrNotFound }}
	svc := service.NewServiceWithRepo(inboundRepo, &mockInventoryRepo{})
	if err := svc.Delete(context.Background(), 99); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}
