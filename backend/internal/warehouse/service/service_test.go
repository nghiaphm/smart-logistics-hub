package service_test

import (
	"context"
	"errors"
	"testing"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/warehouse/dto"
	"my-web-app.com/smart-logistic-hub/internal/warehouse/entity"
	"my-web-app.com/smart-logistic-hub/internal/warehouse/service"
)

type mockWarehouseRepo struct {
	createFn    func(ctx context.Context, w *entity.Warehouse) error
	getByIDFn   func(ctx context.Context, id int64) (*entity.Warehouse, error)
	getByCodeFn func(ctx context.Context, code string) (*entity.Warehouse, error)
	listFn      func(ctx context.Context, offset, limit int) ([]entity.Warehouse, error)
	countFn     func(ctx context.Context) (int, error)
	updateFn    func(ctx context.Context, id int64, w *entity.Warehouse) error
	deleteFn    func(ctx context.Context, id int64) error
}

func (m *mockWarehouseRepo) Create(ctx context.Context, w *entity.Warehouse) error {
	if m.createFn != nil {
		return m.createFn(ctx, w)
	}
	return nil
}

func (m *mockWarehouseRepo) GetByID(ctx context.Context, id int64) (*entity.Warehouse, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockWarehouseRepo) GetByCode(ctx context.Context, code string) (*entity.Warehouse, error) {
	if m.getByCodeFn != nil {
		return m.getByCodeFn(ctx, code)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockWarehouseRepo) List(ctx context.Context, offset, limit int) ([]entity.Warehouse, error) {
	if m.listFn != nil {
		return m.listFn(ctx, offset, limit)
	}
	return []entity.Warehouse{}, nil
}

func (m *mockWarehouseRepo) Count(ctx context.Context) (int, error) {
	if m.countFn != nil {
		return m.countFn(ctx)
	}
	return 0, nil
}

func (m *mockWarehouseRepo) Update(ctx context.Context, id int64, w *entity.Warehouse) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, w)
	}
	return nil
}

func (m *mockWarehouseRepo) Delete(ctx context.Context, id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func TestWarehouseCreateValid(t *testing.T) {
	var saved *entity.Warehouse
	repo := &mockWarehouseRepo{createFn: func(ctx context.Context, w *entity.Warehouse) error {
		w.ID = 1
		saved = w
		return nil
	}}
	svc := service.NewServiceWithRepo(repo)
	w, err := svc.Create(context.Background(), &dto.CreateWarehouseRequest{
		WarehouseCode: "WH-001",
		Name:          "Main Hub",
		Address:       "1 Nguyen Hue",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if w.ID != 1 {
		t.Errorf("Create() ID = %d, want 1", w.ID)
	}
	if !w.IsActive {
		t.Error("Create() IsActive = false, want true (default)")
	}
	if saved == nil {
		t.Fatal("Create() did not call repository.Create")
	}
}

func TestWarehouseCreateRequiresCode(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockWarehouseRepo{})
	_, err := svc.Create(context.Background(), &dto.CreateWarehouseRequest{Name: "Main Hub", Address: "1 Nguyen Hue"})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestWarehouseCreateRequiresName(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockWarehouseRepo{})
	_, err := svc.Create(context.Background(), &dto.CreateWarehouseRequest{WarehouseCode: "WH-001", Address: "1 Nguyen Hue"})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestWarehouseCreateRequiresAddress(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockWarehouseRepo{})
	_, err := svc.Create(context.Background(), &dto.CreateWarehouseRequest{WarehouseCode: "WH-001", Name: "Main Hub"})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestWarehouseCreateReturnsConflictWhenCodeDuplicate(t *testing.T) {
	repo := &mockWarehouseRepo{getByCodeFn: func(ctx context.Context, code string) (*entity.Warehouse, error) {
		return &entity.Warehouse{ID: 1, WarehouseCode: code}, nil
	}}
	svc := service.NewServiceWithRepo(repo)
	_, err := svc.Create(context.Background(), &dto.CreateWarehouseRequest{
		WarehouseCode: "WH-001", Name: "Main Hub", Address: "1 Nguyen Hue",
	})
	if !errors.Is(err, apierrors.ErrConflict) {
		t.Errorf("Create() error = %v, want ErrConflict", err)
	}
}

func TestWarehouseCreateAppliesOptionalFields(t *testing.T) {
	var saved *entity.Warehouse
	repo := &mockWarehouseRepo{createFn: func(ctx context.Context, w *entity.Warehouse) error {
		saved = w
		return nil
	}}
	svc := service.NewServiceWithRepo(repo)
	lat, lng := 10.8231, 106.6297
	phone := "0900000000"
	_, err := svc.Create(context.Background(), &dto.CreateWarehouseRequest{
		WarehouseCode: "WH-002",
		Name:          "South Hub",
		Address:       "2 Ly Tu Trong",
		Location:      &dto.Location{Lat: &lat, Lng: &lng},
		ContactPhone:  &phone,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if saved.Lat != lat || saved.Lng != lng {
		t.Errorf("Create() location = (%.4f, %.4f), want (%.4f, %.4f)", saved.Lat, saved.Lng, lat, lng)
	}
	if saved.ContactPhone != phone {
		t.Errorf("Create() ContactPhone = %q, want %q", saved.ContactPhone, phone)
	}
}

func TestWarehouseCreatePropagatesRepoError(t *testing.T) {
	repoErr := errors.New("db down")
	repo := &mockWarehouseRepo{createFn: func(ctx context.Context, w *entity.Warehouse) error { return repoErr }}
	svc := service.NewServiceWithRepo(repo)
	_, err := svc.Create(context.Background(), &dto.CreateWarehouseRequest{
		WarehouseCode: "WH-001", Name: "Main Hub", Address: "1 Nguyen Hue",
	})
	if !errors.Is(err, repoErr) {
		t.Errorf("Create() error = %v, want %v", err, repoErr)
	}
}

func TestWarehouseGetReturnsErrNotFound(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockWarehouseRepo{})
	_, err := svc.Get(context.Background(), 99)
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestWarehouseListReturnsItemsAndCount(t *testing.T) {
	items := []entity.Warehouse{{ID: 1, WarehouseCode: "WH-001"}}
	repo := &mockWarehouseRepo{
		listFn:  func(ctx context.Context, offset, limit int) ([]entity.Warehouse, error) { return items, nil },
		countFn: func(ctx context.Context) (int, error) { return 1, nil },
	}
	svc := service.NewServiceWithRepo(repo)
	got, total, err := svc.List(context.Background(), 0, 20)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if total != 1 || len(got) != 1 {
		t.Errorf("List() = (%d items, %d total), want (1, 1)", len(got), total)
	}
}

func TestWarehouseUpdateAppliesPartialFields(t *testing.T) {
	existing := &entity.Warehouse{ID: 1, WarehouseCode: "WH-001", Name: "Old", Address: "Old St", IsActive: true}
	repo := &mockWarehouseRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.Warehouse, error) { return existing, nil }}
	svc := service.NewServiceWithRepo(repo)

	newName := "New Hub"
	inactive := false
	w, err := svc.Update(context.Background(), 1, &dto.UpdateWarehouseRequest{Name: &newName, IsActive: &inactive})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if w.Name != newName {
		t.Errorf("Update() Name = %q, want %q", w.Name, newName)
	}
	if w.IsActive {
		t.Error("Update() IsActive = true, want false")
	}
	if w.Address != "Old St" {
		t.Errorf("Update() Address = %q, want unchanged %q", w.Address, "Old St")
	}
}

func TestWarehouseUpdateReturnsErrNotFound(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockWarehouseRepo{})
	_, err := svc.Update(context.Background(), 99, &dto.UpdateWarehouseRequest{})
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestWarehouseDeleteReturnsErrNotFound(t *testing.T) {
	repo := &mockWarehouseRepo{deleteFn: func(ctx context.Context, id int64) error { return apierrors.ErrNotFound }}
	svc := service.NewServiceWithRepo(repo)
	if err := svc.Delete(context.Background(), 99); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}
