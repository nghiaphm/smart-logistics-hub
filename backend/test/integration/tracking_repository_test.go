package integration

import (
	"context"
	"errors"
	"testing"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/tracking/entity"
	trkrepo "my-web-app.com/smart-logistic-hub/internal/tracking/repository"
)

func TestTrackingRepositoryCRUD(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()
	repo := trkrepo.NewRepository(testDB)

	event := &entity.TrackingEvent{
		OrderCode:    "ORD-IT-001",
		DriverCode:   "DRV-IT-001",
		StatusUpdate: "PICKED_UP",
		Lat:          10.5,
		Lng:          106.7,
		Note:         "picked up",
		Timestamp:    time.Now().UTC(),
	}
	if err := repo.Create(ctx, event); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if event.ID == 0 {
		t.Fatal("Create() did not set ID")
	}

	got, err := repo.GetByID(ctx, event.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.OrderCode != event.OrderCode || got.StatusUpdate != "PICKED_UP" {
		t.Errorf("GetByID() = %+v", got)
	}

	event.StatusUpdate = "IN_TRANSIT"
	if err := repo.Update(ctx, event.ID, event); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	got, err = repo.GetByID(ctx, event.ID)
	if err != nil {
		t.Fatalf("GetByID() after update error = %v", err)
	}
	if got.StatusUpdate != "IN_TRANSIT" {
		t.Errorf("Update() StatusUpdate = %q, want %q", got.StatusUpdate, "IN_TRANSIT")
	}

	if err := repo.Delete(ctx, event.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := repo.GetByID(ctx, event.ID); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("GetByID() after delete error = %v, want ErrNotFound", err)
	}
}

func TestTrackingRepositoryNotFound(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()
	repo := trkrepo.NewRepository(testDB)

	if _, err := repo.GetByID(ctx, 999999); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("GetByID() error = %v, want ErrNotFound", err)
	}
	if err := repo.Update(ctx, 999999, &entity.TrackingEvent{}); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
	if err := repo.Delete(ctx, 999999); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}

func TestTrackingRepositoryListGetByOrder(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()
	repo := trkrepo.NewRepository(testDB)

	for i := 1; i <= 3; i++ {
		e := &entity.TrackingEvent{
			OrderCode:    "ORD-IT-001",
			DriverCode:   "DRV-IT-001",
			StatusUpdate: "STATUS_" + string(rune('0'+i)),
			Timestamp:    time.Now().UTC().Add(time.Duration(i) * time.Minute),
		}
		if err := repo.Create(ctx, e); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	items, err := repo.List(ctx, "ORD-IT-001", "", 0, 10)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(items) != 3 {
		t.Errorf("List() len = %d, want 3", len(items))
	}
	count, err := repo.Count(ctx, "ORD-IT-001", "")
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}
	if count != 3 {
		t.Errorf("Count() = %d, want 3", count)
	}

	timeline, err := repo.GetByOrder(ctx, "ORD-IT-001")
	if err != nil {
		t.Fatalf("GetByOrder() error = %v", err)
	}
	if len(timeline) != 3 {
		t.Errorf("GetByOrder() len = %d, want 3", len(timeline))
	}
	if timeline[0].Timestamp.After(timeline[2].Timestamp) {
		t.Error("GetByOrder() should return events sorted by timestamp ASC")
	}
}
