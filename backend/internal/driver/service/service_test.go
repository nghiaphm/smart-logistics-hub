package service_test

import (
	"context"
	"errors"
	"testing"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/driver/dto"
	"my-web-app.com/smart-logistic-hub/internal/driver/entity"
	"my-web-app.com/smart-logistic-hub/internal/driver/service"
)

type mockDriverRepo struct {
	createFn  func(ctx context.Context, d *entity.Driver) error
	getByIDFn func(ctx context.Context, id int64) (*entity.Driver, error)
	listFn    func(ctx context.Context, status string, offset, limit int) ([]entity.Driver, error)
	countFn   func(ctx context.Context, status string) (int, error)
	updateFn  func(ctx context.Context, id int64, d *entity.Driver) error
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *mockDriverRepo) Create(ctx context.Context, d *entity.Driver) error {
	if m.createFn != nil {
		return m.createFn(ctx, d)
	}
	return nil
}

func (m *mockDriverRepo) GetByID(ctx context.Context, id int64) (*entity.Driver, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockDriverRepo) GetByCode(ctx context.Context, code string) (*entity.Driver, error) {
	return nil, apierrors.ErrNotFound
}

func (m *mockDriverRepo) List(ctx context.Context, status string, offset, limit int) ([]entity.Driver, error) {
	if m.listFn != nil {
		return m.listFn(ctx, status, offset, limit)
	}
	return []entity.Driver{}, nil
}

func (m *mockDriverRepo) Count(ctx context.Context, status string) (int, error) {
	if m.countFn != nil {
		return m.countFn(ctx, status)
	}
	return 0, nil
}

func (m *mockDriverRepo) Update(ctx context.Context, id int64, d *entity.Driver) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, d)
	}
	return nil
}

func (m *mockDriverRepo) Delete(ctx context.Context, id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func TestCreateDefaultsStatusToAvailable(t *testing.T) {
	var saved *entity.Driver
	repo := &mockDriverRepo{
		createFn: func(ctx context.Context, d *entity.Driver) error {
			d.ID = 1
			saved = d
			return nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	d, err := svc.Create(context.Background(), &dto.CreateDriverRequest{
		DriverCode: "DRV001",
		FullName:   "Nguyen Van A",
		Phone:      "0900000000",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if d.Status != "AVAILABLE" {
		t.Errorf("Create() status = %q, want %q", d.Status, "AVAILABLE")
	}
	if saved == nil {
		t.Fatal("Create() did not call repository.Create")
	}
	if saved.CreatedAt.IsZero() || saved.UpdatedAt.IsZero() {
		t.Error("Create() should set CreatedAt and UpdatedAt")
	}
}

func TestCreateKeepsProvidedStatus(t *testing.T) {
	repo := &mockDriverRepo{
		createFn: func(ctx context.Context, d *entity.Driver) error {
			d.ID = 1
			return nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	status := "BUSY"
	d, err := svc.Create(context.Background(), &dto.CreateDriverRequest{
		DriverCode: "DRV002",
		Status:     status,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if d.Status != status {
		t.Errorf("Create() status = %q, want %q", d.Status, status)
	}
}

func TestCreatePropagatesRepositoryError(t *testing.T) {
	repoErr := errors.New("database down")
	repo := &mockDriverRepo{
		createFn: func(ctx context.Context, d *entity.Driver) error {
			return repoErr
		},
	}
	svc := service.NewServiceWithRepo(repo)

	_, err := svc.Create(context.Background(), &dto.CreateDriverRequest{DriverCode: "DRV003"})
	if !errors.Is(err, repoErr) {
		t.Errorf("Create() error = %v, want %v", err, repoErr)
	}
}

func TestGetReturnsDriver(t *testing.T) {
	expected := &entity.Driver{ID: 7, DriverCode: "DRV007"}
	repo := &mockDriverRepo{
		getByIDFn: func(ctx context.Context, id int64) (*entity.Driver, error) {
			return expected, nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	d, err := svc.Get(context.Background(), 7)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if d != expected {
		t.Errorf("Get() = %v, want %v", d, expected)
	}
}

func TestGetReturnsErrNotFound(t *testing.T) {
	repo := &mockDriverRepo{}
	svc := service.NewServiceWithRepo(repo)

	_, err := svc.Get(context.Background(), 99)
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestListReturnsItemsAndCount(t *testing.T) {
	items := []entity.Driver{{ID: 1, DriverCode: "DRV001"}, {ID: 2, DriverCode: "DRV002"}}
	repo := &mockDriverRepo{
		listFn: func(ctx context.Context, status string, offset, limit int) ([]entity.Driver, error) {
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

func TestUpdateAppliesPartialFields(t *testing.T) {
	existing := &entity.Driver{ID: 5, DriverCode: "DRV005", FullName: "Old Name", Status: "AVAILABLE"}
	repo := &mockDriverRepo{
		getByIDFn: func(ctx context.Context, id int64) (*entity.Driver, error) {
			return existing, nil
		},
		updateFn: func(ctx context.Context, id int64, d *entity.Driver) error {
			return nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	newName := "New Name"
	newStatus := "OFFLINE"
	d, err := svc.Update(context.Background(), 5, &dto.UpdateDriverRequest{
		FullName: &newName,
		Status:   &newStatus,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if d.FullName != newName {
		t.Errorf("Update() FullName = %q, want %q", d.FullName, newName)
	}
	if d.Status != newStatus {
		t.Errorf("Update() Status = %q, want %q", d.Status, newStatus)
	}
	if d.DriverCode != "DRV005" {
		t.Errorf("Update() DriverCode = %q, want unchanged %q", d.DriverCode, "DRV005")
	}
}

func TestUpdateReturnsErrNotFound(t *testing.T) {
	repo := &mockDriverRepo{}
	svc := service.NewServiceWithRepo(repo)

	_, err := svc.Update(context.Background(), 99, &dto.UpdateDriverRequest{})
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestDeleteCallsRepository(t *testing.T) {
	var deletedID int64
	repo := &mockDriverRepo{
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
	repo := &mockDriverRepo{
		deleteFn: func(ctx context.Context, id int64) error {
			return apierrors.ErrNotFound
		},
	}
	svc := service.NewServiceWithRepo(repo)

	if err := svc.Delete(context.Background(), 99); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}
