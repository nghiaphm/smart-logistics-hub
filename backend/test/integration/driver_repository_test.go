package integration

import (
	"context"
	"errors"
	"testing"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/driver/entity"
	"my-web-app.com/smart-logistic-hub/internal/driver/repository"
)

func TestDriverRepositoryCRUD(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()
	repo := repository.NewRepository(testDB)

	driver := &entity.Driver{
		DriverCode:   "DRV-IT-001",
		FullName:     "Test Driver",
		Phone:        "0900000001",
		VehicleType:  "truck",
		LicensePlate: "51A-12345",
		Status:       "AVAILABLE",
	}
	if err := repo.Create(ctx, driver); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if driver.ID == 0 {
		t.Fatal("Create() did not set ID")
	}

	got, err := repo.GetByID(ctx, driver.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.DriverCode != driver.DriverCode || got.FullName != driver.FullName {
		t.Errorf("GetByID() = %+v, want code %q", got, driver.DriverCode)
	}

	got, err = repo.GetByCode(ctx, "DRV-IT-001")
	if err != nil {
		t.Fatalf("GetByCode() error = %v", err)
	}
	if got.ID != driver.ID {
		t.Errorf("GetByCode() ID = %d, want %d", got.ID, driver.ID)
	}

	driver.Status = "BUSY"
	if err := repo.Update(ctx, driver.ID, driver); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	got, err = repo.GetByID(ctx, driver.ID)
	if err != nil {
		t.Fatalf("GetByID() after update error = %v", err)
	}
	if got.Status != "BUSY" {
		t.Errorf("Update() status = %q, want %q", got.Status, "BUSY")
	}

	if err := repo.Delete(ctx, driver.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := repo.GetByID(ctx, driver.ID); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("GetByID() after delete error = %v, want ErrNotFound", err)
	}
}

func TestDriverRepositoryNotFound(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()
	repo := repository.NewRepository(testDB)

	if _, err := repo.GetByID(ctx, 999999); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("GetByID() error = %v, want ErrNotFound", err)
	}
	if _, err := repo.GetByCode(ctx, "NOPE"); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("GetByCode() error = %v, want ErrNotFound", err)
	}
	if err := repo.Update(ctx, 999999, &entity.Driver{}); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
	if err := repo.Delete(ctx, 999999); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}

func TestDriverRepositoryListAndCount(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()
	repo := repository.NewRepository(testDB)

	for i := 1; i <= 3; i++ {
		status := "AVAILABLE"
		if i == 3 {
			status = "BUSY"
		}
		d := &entity.Driver{
			DriverCode: "DRV-IT-LST" + string(rune('0'+i)),
			FullName:   "Driver " + string(rune('0'+i)),
			Status:     status,
		}
		if err := repo.Create(ctx, d); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	all, err := repo.List(ctx, "", 0, 10)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(all) != 3 {
		t.Errorf("List() len = %d, want 3", len(all))
	}
	count, err := repo.Count(ctx, "")
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}
	if count != 3 {
		t.Errorf("Count() = %d, want 3", count)
	}

	busy, err := repo.List(ctx, "BUSY", 0, 10)
	if err != nil {
		t.Fatalf("List(BUSY) error = %v", err)
	}
	if len(busy) != 1 {
		t.Errorf("List(BUSY) len = %d, want 1", len(busy))
	}
	busyCount, err := repo.Count(ctx, "BUSY")
	if err != nil {
		t.Fatalf("Count(BUSY) error = %v", err)
	}
	if busyCount != 1 {
		t.Errorf("Count(BUSY) = %d, want 1", busyCount)
	}
}
