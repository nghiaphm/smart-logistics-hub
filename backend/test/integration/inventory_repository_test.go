package integration

import (
	"context"
	"errors"
	"testing"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/inventory/entity"
	invrepo "my-web-app.com/smart-logistic-hub/internal/inventory/repository"
)

func TestInventoryRepositoryCRUD(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()
	repo := invrepo.NewRepository(testDB)

	inv := &entity.Inventory{
		ProductID:    100,
		WarehouseID:  200,
		AvailableQty: 10,
		ReservedQty:  2,
		DamagedQty:   0,
		HoldQty:      0,
	}
	if err := repo.Create(ctx, inv); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if inv.ID == 0 {
		t.Fatal("Create() did not set ID")
	}

	got, err := repo.GetByID(ctx, inv.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.ProductID != 100 || got.AvailableQty != 10 {
		t.Errorf("GetByID() = %+v", got)
	}

	got, err = repo.GetByProductWarehouse(ctx, 100, 200)
	if err != nil {
		t.Fatalf("GetByProductWarehouse() error = %v", err)
	}
	if got.ID != inv.ID {
		t.Errorf("GetByProductWarehouse() ID = %d, want %d", got.ID, inv.ID)
	}

	inv.AvailableQty = 8
	if err := repo.Update(ctx, inv.ID, inv); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	got, err = repo.GetByID(ctx, inv.ID)
	if err != nil {
		t.Fatalf("GetByID() after update error = %v", err)
	}
	if got.AvailableQty != 8 {
		t.Errorf("Update() AvailableQty = %d, want 8", got.AvailableQty)
	}

	if err := repo.Delete(ctx, inv.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := repo.GetByID(ctx, inv.ID); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("GetByID() after delete error = %v, want ErrNotFound", err)
	}
}

func TestInventoryRepositoryNotFound(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()
	repo := invrepo.NewRepository(testDB)

	if _, err := repo.GetByID(ctx, 999999); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("GetByID() error = %v, want ErrNotFound", err)
	}
	if _, err := repo.GetByProductWarehouse(ctx, 1, 2); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("GetByProductWarehouse() error = %v, want ErrNotFound", err)
	}
	if err := repo.Update(ctx, 999999, &entity.Inventory{}); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
	if err := repo.Delete(ctx, 999999); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}

func TestInventoryRepositoryListAndCount(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()
	repo := invrepo.NewRepository(testDB)

	for i := int64(1); i <= 3; i++ {
		inv := &entity.Inventory{
			ProductID:    1000 + i,
			WarehouseID:  2000 + i,
			AvailableQty: int(i),
		}
		if err := repo.Create(ctx, inv); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	items, err := repo.List(ctx, 0, 10)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(items) != 3 {
		t.Errorf("List() len = %d, want 3", len(items))
	}
	count, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}
	if count != 3 {
		t.Errorf("Count() = %d, want 3", count)
	}
}
