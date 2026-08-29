package service_test

import (
	"context"
	"errors"
	"testing"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/vehicle/dto"
	"my-web-app.com/smart-logistic-hub/internal/vehicle/entity"
	"my-web-app.com/smart-logistic-hub/internal/vehicle/service"
)

type mockVehicleRepo struct {
	createFn  func(ctx context.Context, v *entity.Vehicle) error
	getByIDFn func(ctx context.Context, id int64) (*entity.Vehicle, error)
	listFn    func(ctx context.Context, status string, offset, limit int) ([]entity.Vehicle, error)
	countFn   func(ctx context.Context, status string) (int, error)
	updateFn  func(ctx context.Context, id int64, v *entity.Vehicle) error
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *mockVehicleRepo) Create(ctx context.Context, v *entity.Vehicle) error {
	if m.createFn != nil {
		return m.createFn(ctx, v)
	}
	return nil
}

func (m *mockVehicleRepo) GetByID(ctx context.Context, id int64) (*entity.Vehicle, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockVehicleRepo) List(ctx context.Context, status string, offset, limit int) ([]entity.Vehicle, error) {
	if m.listFn != nil {
		return m.listFn(ctx, status, offset, limit)
	}
	return []entity.Vehicle{}, nil
}

func (m *mockVehicleRepo) Count(ctx context.Context, status string) (int, error) {
	if m.countFn != nil {
		return m.countFn(ctx, status)
	}
	return 0, nil
}

func (m *mockVehicleRepo) Update(ctx context.Context, id int64, v *entity.Vehicle) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, v)
	}
	return nil
}

func (m *mockVehicleRepo) Delete(ctx context.Context, id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func TestCreateDefaultsStatusToActive(t *testing.T) {
	var saved *entity.Vehicle
	repo := &mockVehicleRepo{
		createFn: func(ctx context.Context, v *entity.Vehicle) error {
			v.ID = 1
			saved = v
			return nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	v, err := svc.Create(context.Background(), &dto.CreateVehicleRequest{
		LicensePlate: "51F-123.45",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if v.Status != "ACTIVE" {
		t.Errorf("Create() status = %q, want %q", v.Status, "ACTIVE")
	}
	if saved == nil {
		t.Fatal("Create() did not call repository.Create")
	}
	if saved.CreatedAt.IsZero() || saved.UpdatedAt.IsZero() {
		t.Error("Create() should set CreatedAt and UpdatedAt")
	}
}

func TestCreateKeepsProvidedFields(t *testing.T) {
	repo := &mockVehicleRepo{
		createFn: func(ctx context.Context, v *entity.Vehicle) error {
			v.ID = 1
			return nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	v, err := svc.Create(context.Background(), &dto.CreateVehicleRequest{
		LicensePlate: "51F-123.45",
		Type:         "TRUCK",
		Capacity:     1500,
		Status:       "MAINTENANCE",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if v.Type != "TRUCK" {
		t.Errorf("Create() Type = %q, want %q", v.Type, "TRUCK")
	}
	if v.Capacity != 1500 {
		t.Errorf("Create() Capacity = %v, want 1500", v.Capacity)
	}
	if v.Status != "MAINTENANCE" {
		t.Errorf("Create() Status = %q, want %q", v.Status, "MAINTENANCE")
	}
}

func TestCreatePropagatesRepositoryError(t *testing.T) {
	repoErr := errors.New("database down")
	repo := &mockVehicleRepo{
		createFn: func(ctx context.Context, v *entity.Vehicle) error {
			return repoErr
		},
	}
	svc := service.NewServiceWithRepo(repo)

	_, err := svc.Create(context.Background(), &dto.CreateVehicleRequest{LicensePlate: "51F-123.45"})
	if !errors.Is(err, repoErr) {
		t.Errorf("Create() error = %v, want %v", err, repoErr)
	}
}

func TestGetReturnsVehicle(t *testing.T) {
	expected := &entity.Vehicle{ID: 7, LicensePlate: "51F-999.99"}
	repo := &mockVehicleRepo{
		getByIDFn: func(ctx context.Context, id int64) (*entity.Vehicle, error) {
			return expected, nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	v, err := svc.Get(context.Background(), 7)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if v != expected {
		t.Errorf("Get() = %v, want %v", v, expected)
	}
}

func TestGetReturnsErrNotFound(t *testing.T) {
	repo := &mockVehicleRepo{}
	svc := service.NewServiceWithRepo(repo)

	_, err := svc.Get(context.Background(), 99)
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestListReturnsItemsAndCount(t *testing.T) {
	items := []entity.Vehicle{{ID: 1, LicensePlate: "51F-111.11"}, {ID: 2, LicensePlate: "51F-222.22"}}
	repo := &mockVehicleRepo{
		listFn: func(ctx context.Context, status string, offset, limit int) ([]entity.Vehicle, error) {
			return items, nil
		},
		countFn: func(ctx context.Context, status string) (int, error) {
			return 2, nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	got, total, err := svc.List(context.Background(), "", 0, 20)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if total != 2 {
		t.Errorf("List() total = %d, want 2", total)
	}
	if len(got) != len(items) {
		t.Errorf("List() len = %d, want %d", len(got), len(items))
	}
}

func TestListPassesStatusFilter(t *testing.T) {
	var gotStatus string
	repo := &mockVehicleRepo{
		listFn: func(ctx context.Context, status string, offset, limit int) ([]entity.Vehicle, error) {
			gotStatus = status
			return []entity.Vehicle{}, nil
		},
		countFn: func(ctx context.Context, status string) (int, error) {
			return 0, nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	_, _, err := svc.List(context.Background(), "ACTIVE", 0, 20)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if gotStatus != "ACTIVE" {
		t.Errorf("List() status = %q, want %q", gotStatus, "ACTIVE")
	}
}

func TestUpdateAppliesPartialFields(t *testing.T) {
	existing := &entity.Vehicle{ID: 5, LicensePlate: "51F-123.45", Type: "TRUCK", Status: "ACTIVE"}
	repo := &mockVehicleRepo{
		getByIDFn: func(ctx context.Context, id int64) (*entity.Vehicle, error) {
			return existing, nil
		},
		updateFn: func(ctx context.Context, id int64, v *entity.Vehicle) error {
			return nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	newType := "VAN"
	newStatus := "INACTIVE"
	newCapacity := 800.5
	v, err := svc.Update(context.Background(), 5, &dto.UpdateVehicleRequest{
		Type:     &newType,
		Status:   &newStatus,
		Capacity: &newCapacity,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if v.Type != newType {
		t.Errorf("Update() Type = %q, want %q", v.Type, newType)
	}
	if v.Status != newStatus {
		t.Errorf("Update() Status = %q, want %q", v.Status, newStatus)
	}
	if v.Capacity != newCapacity {
		t.Errorf("Update() Capacity = %v, want %v", v.Capacity, newCapacity)
	}
	if v.LicensePlate != "51F-123.45" {
		t.Errorf("Update() LicensePlate = %q, want unchanged %q", v.LicensePlate, "51F-123.45")
	}
}

func TestUpdateReturnsErrNotFound(t *testing.T) {
	repo := &mockVehicleRepo{}
	svc := service.NewServiceWithRepo(repo)

	_, err := svc.Update(context.Background(), 99, &dto.UpdateVehicleRequest{})
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestDeleteCallsRepository(t *testing.T) {
	var deletedID int64
	repo := &mockVehicleRepo{
		deleteFn: func(ctx context.Context, id int64) error {
			deletedID = id
			return nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	if err := svc.Delete(context.Background(), 5); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if deletedID != 5 {
		t.Errorf("Delete() deletedID = %d, want 5", deletedID)
	}
}

func TestDeleteReturnsErrNotFound(t *testing.T) {
	repo := &mockVehicleRepo{
		deleteFn: func(ctx context.Context, id int64) error {
			return apierrors.ErrNotFound
		},
	}
	svc := service.NewServiceWithRepo(repo)

	if err := svc.Delete(context.Background(), 99); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}
