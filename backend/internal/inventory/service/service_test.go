package service_test

import (
	"context"
	"errors"
	"testing"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/inventory/dto"
	"my-web-app.com/smart-logistic-hub/internal/inventory/entity"
	"my-web-app.com/smart-logistic-hub/internal/inventory/service"
)

type mockInventoryRepo struct {
	createFn  func(ctx context.Context, inv *entity.Inventory) error
	getByIDFn func(ctx context.Context, id int64) (*entity.Inventory, error)
	getByPWFn func(ctx context.Context, productID, warehouseID int64) (*entity.Inventory, error)
	listFn    func(ctx context.Context, offset, limit int) ([]entity.Inventory, error)
	countFn   func(ctx context.Context) (int, error)
	updateFn  func(ctx context.Context, id int64, inv *entity.Inventory) error
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *mockInventoryRepo) Create(ctx context.Context, inv *entity.Inventory) error {
	if m.createFn != nil {
		return m.createFn(ctx, inv)
	}
	return nil
}
func (m *mockInventoryRepo) GetByID(ctx context.Context, id int64) (*entity.Inventory, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, apierrors.ErrNotFound
}
func (m *mockInventoryRepo) GetByProductWarehouse(ctx context.Context, productID, warehouseID int64) (*entity.Inventory, error) {
	if m.getByPWFn != nil {
		return m.getByPWFn(ctx, productID, warehouseID)
	}
	return nil, apierrors.ErrNotFound
}
func (m *mockInventoryRepo) List(ctx context.Context, offset, limit int) ([]entity.Inventory, error) {
	if m.listFn != nil {
		return m.listFn(ctx, offset, limit)
	}
	return []entity.Inventory{}, nil
}
func (m *mockInventoryRepo) Count(ctx context.Context) (int, error) {
	if m.countFn != nil {
		return m.countFn(ctx)
	}
	return 0, nil
}
func (m *mockInventoryRepo) Update(ctx context.Context, id int64, inv *entity.Inventory) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, inv)
	}
	return nil
}
func (m *mockInventoryRepo) Delete(ctx context.Context, id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func TestInventoryCreateValid(t *testing.T) {
	repo := &mockInventoryRepo{createFn: func(ctx context.Context, inv *entity.Inventory) error { inv.ID = 1; return nil }}
	svc := service.NewServiceWithRepo(repo)
	inv, err := svc.Create(context.Background(), &dto.CreateInventoryRequest{ProductID: 10, WarehouseID: 20, AvailableQty: 5})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if inv.ID != 1 {
		t.Errorf("Create() ID = %d, want 1", inv.ID)
	}
}

func TestInventoryCreateRequiresProductID(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockInventoryRepo{})
	_, err := svc.Create(context.Background(), &dto.CreateInventoryRequest{ProductID: 0, WarehouseID: 20})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestInventoryCreateRequiresWarehouseID(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockInventoryRepo{})
	_, err := svc.Create(context.Background(), &dto.CreateInventoryRequest{ProductID: 10, WarehouseID: 0})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestInventoryCreateRejectsNegativeQty(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockInventoryRepo{})
	_, err := svc.Create(context.Background(), &dto.CreateInventoryRequest{ProductID: 10, WarehouseID: 20, AvailableQty: -1})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestInventoryCreateReturnsConflictWhenDuplicate(t *testing.T) {
	repo := &mockInventoryRepo{getByPWFn: func(ctx context.Context, productID, warehouseID int64) (*entity.Inventory, error) {
		return &entity.Inventory{ID: 1, ProductID: productID, WarehouseID: warehouseID}, nil
	}}
	svc := service.NewServiceWithRepo(repo)
	_, err := svc.Create(context.Background(), &dto.CreateInventoryRequest{ProductID: 10, WarehouseID: 20})
	if !errors.Is(err, apierrors.ErrConflict) {
		t.Errorf("Create() error = %v, want ErrConflict", err)
	}
}

func TestInventoryCreatePropagatesRepoError(t *testing.T) {
	repoErr := errors.New("db down")
	repo := &mockInventoryRepo{createFn: func(ctx context.Context, inv *entity.Inventory) error { return repoErr }}
	svc := service.NewServiceWithRepo(repo)
	_, err := svc.Create(context.Background(), &dto.CreateInventoryRequest{ProductID: 10, WarehouseID: 20})
	if !errors.Is(err, repoErr) {
		t.Errorf("Create() error = %v, want %v", err, repoErr)
	}
}

func TestInventoryGetReturnsErrNotFound(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockInventoryRepo{})
	_, err := svc.Get(context.Background(), 99)
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestInventoryListReturnsItemsAndCount(t *testing.T) {
	items := []entity.Inventory{{ID: 1, ProductID: 10}}
	repo := &mockInventoryRepo{listFn: func(ctx context.Context, offset, limit int) ([]entity.Inventory, error) { return items, nil }, countFn: func(ctx context.Context) (int, error) { return 1, nil }}
	svc := service.NewServiceWithRepo(repo)
	got, total, err := svc.List(context.Background(), 0, 20)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if total != 1 || len(got) != 1 {
		t.Errorf("List() = (%d items, %d total), want (1, 1)", len(got), total)
	}
}

func TestInventoryUpdateRejectsNegativeAvailable(t *testing.T) {
	repo := &mockInventoryRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.Inventory, error) {
		return &entity.Inventory{ID: 1, AvailableQty: 10}, nil
	}}
	svc := service.NewServiceWithRepo(repo)
	neg := -1
	_, err := svc.Update(context.Background(), 1, &dto.UpdateInventoryRequest{AvailableQty: &neg})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Update() error = %v, want ErrBadRequest", err)
	}
}

func TestInventoryUpdateAppliesPartialFields(t *testing.T) {
	existing := &entity.Inventory{ID: 1, AvailableQty: 10, ReservedQty: 2}
	repo := &mockInventoryRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.Inventory, error) { return existing, nil }}
	svc := service.NewServiceWithRepo(repo)
	newAvail := 7
	inv, err := svc.Update(context.Background(), 1, &dto.UpdateInventoryRequest{AvailableQty: &newAvail})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if inv.AvailableQty != 7 {
		t.Errorf("Update() AvailableQty = %d, want 7", inv.AvailableQty)
	}
	if inv.ReservedQty != 2 {
		t.Errorf("Update() ReservedQty = %d, want unchanged 2", inv.ReservedQty)
	}
}

func TestInventoryUpdateReturnsErrNotFound(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockInventoryRepo{})
	_, err := svc.Update(context.Background(), 99, &dto.UpdateInventoryRequest{})
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestInventoryDeleteReturnsErrNotFound(t *testing.T) {
	repo := &mockInventoryRepo{deleteFn: func(ctx context.Context, id int64) error { return apierrors.ErrNotFound }}
	svc := service.NewServiceWithRepo(repo)
	if err := svc.Delete(context.Background(), 99); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}
