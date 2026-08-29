package integration

import (
	"context"
	"errors"
	"testing"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/vehicle/entity"
	"my-web-app.com/smart-logistic-hub/internal/vehicle/repository"
)

func TestVehicleRepositoryCRUD(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()
	repo := repository.NewRepository(testDB)

	vehicle := &entity.Vehicle{
		LicensePlate: "51F-IT-001",
		Type:         "TRUCK",
		Capacity:     1500,
		Status:       "ACTIVE",
	}
	if err := repo.Create(ctx, vehicle); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if vehicle.ID == 0 {
		t.Fatal("Create() did not set ID")
	}

	got, err := repo.GetByID(ctx, vehicle.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.LicensePlate != vehicle.LicensePlate || got.Status != "ACTIVE" {
		t.Errorf("GetByID() = %+v, want plate %q", got, vehicle.LicensePlate)
	}

	vehicle.Status = "MAINTENANCE"
	if err := repo.Update(ctx, vehicle.ID, vehicle); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	got, err = repo.GetByID(ctx, vehicle.ID)
	if err != nil {
		t.Fatalf("GetByID() after update error = %v", err)
	}
	if got.Status != "MAINTENANCE" {
		t.Errorf("Update() status = %q, want %q", got.Status, "MAINTENANCE")
	}

	if err := repo.Delete(ctx, vehicle.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := repo.GetByID(ctx, vehicle.ID); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("GetByID() after delete error = %v, want ErrNotFound", err)
	}
}

func TestVehicleRepositoryNotFound(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()
	repo := repository.NewRepository(testDB)

	if _, err := repo.GetByID(ctx, 999999); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("GetByID() error = %v, want ErrNotFound", err)
	}
	if err := repo.Update(ctx, 999999, &entity.Vehicle{}); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
	if err := repo.Delete(ctx, 999999); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}

func TestVehicleRepositoryListAndCount(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()
	repo := repository.NewRepository(testDB)

	for i := 1; i <= 3; i++ {
		status := "ACTIVE"
		if i == 3 {
			status = "INACTIVE"
		}
		v := &entity.Vehicle{
			LicensePlate: "51F-IT-LST" + string(rune('0'+i)),
			Status:       status,
		}
		if err := repo.Create(ctx, v); err != nil {
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

	inactive, err := repo.List(ctx, "INACTIVE", 0, 10)
	if err != nil {
		t.Fatalf("List(INACTIVE) error = %v", err)
	}
	if len(inactive) != 1 {
		t.Errorf("List(INACTIVE) len = %d, want 1", len(inactive))
	}
	inactiveCount, err := repo.Count(ctx, "INACTIVE")
	if err != nil {
		t.Fatalf("Count(INACTIVE) error = %v", err)
	}
	if inactiveCount != 1 {
		t.Errorf("Count(INACTIVE) = %d, want 1", inactiveCount)
	}
}
