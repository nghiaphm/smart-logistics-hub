package service_test

import (
	"context"
	"errors"
	"testing"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/product/dto"
	"my-web-app.com/smart-logistic-hub/internal/product/entity"
	"my-web-app.com/smart-logistic-hub/internal/product/service"
)

type mockProductRepo struct {
	createFn   func(ctx context.Context, p *entity.Product) error
	getByIDFn  func(ctx context.Context, id int64) (*entity.Product, error)
	getBySkuFn func(ctx context.Context, sku string) (*entity.Product, error)
	listFn     func(ctx context.Context, offset, limit int) ([]entity.Product, error)
	countFn    func(ctx context.Context) (int, error)
	updateFn   func(ctx context.Context, id int64, p *entity.Product) error
	deleteFn   func(ctx context.Context, id int64) error
}

func (m *mockProductRepo) Create(ctx context.Context, p *entity.Product) error {
	if m.createFn != nil {
		return m.createFn(ctx, p)
	}
	return nil
}

func (m *mockProductRepo) GetByID(ctx context.Context, id int64) (*entity.Product, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockProductRepo) GetBySku(ctx context.Context, sku string) (*entity.Product, error) {
	if m.getBySkuFn != nil {
		return m.getBySkuFn(ctx, sku)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockProductRepo) List(ctx context.Context, offset, limit int) ([]entity.Product, error) {
	if m.listFn != nil {
		return m.listFn(ctx, offset, limit)
	}
	return []entity.Product{}, nil
}

func (m *mockProductRepo) Count(ctx context.Context) (int, error) {
	if m.countFn != nil {
		return m.countFn(ctx)
	}
	return 0, nil
}

func (m *mockProductRepo) Update(ctx context.Context, id int64, p *entity.Product) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, p)
	}
	return nil
}

func (m *mockProductRepo) Delete(ctx context.Context, id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func TestProductCreateValid(t *testing.T) {
	repo := &mockProductRepo{createFn: func(ctx context.Context, p *entity.Product) error { p.ID = 1; return nil }}
	svc := service.NewServiceWithRepo(repo)
	p, err := svc.Create(context.Background(), &dto.CreateProductRequest{Sku: "SKU-001", Name: "Laptop", Price: 999.5, WeightGram: 2000})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if p.ID != 1 {
		t.Errorf("Create() ID = %d, want 1", p.ID)
	}
	if p.Sku != "SKU-001" {
		t.Errorf("Create() Sku = %q, want %q", p.Sku, "SKU-001")
	}
}

func TestProductCreateRequiresSku(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockProductRepo{})
	_, err := svc.Create(context.Background(), &dto.CreateProductRequest{Name: "Laptop"})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestProductCreateRequiresName(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockProductRepo{})
	_, err := svc.Create(context.Background(), &dto.CreateProductRequest{Sku: "SKU-001"})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestProductCreateRejectsNegativePrice(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockProductRepo{})
	_, err := svc.Create(context.Background(), &dto.CreateProductRequest{Sku: "SKU-001", Name: "Laptop", Price: -1})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestProductCreateRejectsNegativeWeight(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockProductRepo{})
	_, err := svc.Create(context.Background(), &dto.CreateProductRequest{Sku: "SKU-001", Name: "Laptop", WeightGram: -1})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestProductCreateReturnsConflictWhenSkuDuplicate(t *testing.T) {
	repo := &mockProductRepo{getBySkuFn: func(ctx context.Context, sku string) (*entity.Product, error) {
		return &entity.Product{ID: 1, Sku: sku}, nil
	}}
	svc := service.NewServiceWithRepo(repo)
	_, err := svc.Create(context.Background(), &dto.CreateProductRequest{Sku: "SKU-001", Name: "Laptop"})
	if !errors.Is(err, apierrors.ErrConflict) {
		t.Errorf("Create() error = %v, want ErrConflict", err)
	}
}

func TestProductCreateAppliesDimensions(t *testing.T) {
	var saved *entity.Product
	repo := &mockProductRepo{createFn: func(ctx context.Context, p *entity.Product) error {
		saved = p
		return nil
	}}
	svc := service.NewServiceWithRepo(repo)
	length, width, height := 30.0, 20.0, 1.5
	_, err := svc.Create(context.Background(), &dto.CreateProductRequest{
		Sku: "SKU-002", Name: "Monitor",
		Dimensions: &dto.Dimensions{Length: &length, Width: &width, Height: &height},
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if saved.LengthCm != 30 || saved.WidthCm != 20 || saved.HeightCm != 1.5 {
		t.Errorf("Create() dimensions = (%.1f, %.1f, %.1f), want (30.0, 20.0, 1.5)", saved.LengthCm, saved.WidthCm, saved.HeightCm)
	}
}

func TestProductCreatePropagatesRepoError(t *testing.T) {
	repoErr := errors.New("db down")
	repo := &mockProductRepo{createFn: func(ctx context.Context, p *entity.Product) error { return repoErr }}
	svc := service.NewServiceWithRepo(repo)
	_, err := svc.Create(context.Background(), &dto.CreateProductRequest{Sku: "SKU-001", Name: "Laptop"})
	if !errors.Is(err, repoErr) {
		t.Errorf("Create() error = %v, want %v", err, repoErr)
	}
}

func TestProductGetReturnsErrNotFound(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockProductRepo{})
	_, err := svc.Get(context.Background(), 99)
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestProductListReturnsItemsAndCount(t *testing.T) {
	items := []entity.Product{{ID: 1, Sku: "SKU-001"}}
	repo := &mockProductRepo{
		listFn:  func(ctx context.Context, offset, limit int) ([]entity.Product, error) { return items, nil },
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

func TestProductUpdateAppliesPartialFields(t *testing.T) {
	existing := &entity.Product{ID: 1, Sku: "SKU-001", Name: "Old", Price: 100}
	repo := &mockProductRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.Product, error) { return existing, nil }}
	svc := service.NewServiceWithRepo(repo)

	newName := "New"
	newPrice := 150.0
	p, err := svc.Update(context.Background(), 1, &dto.UpdateProductRequest{Name: &newName, Price: &newPrice})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if p.Name != newName {
		t.Errorf("Update() Name = %q, want %q", p.Name, newName)
	}
	if p.Price != newPrice {
		t.Errorf("Update() Price = %v, want %v", p.Price, newPrice)
	}
	if p.Sku != "SKU-001" {
		t.Errorf("Update() Sku = %q, want unchanged %q", p.Sku, "SKU-001")
	}
}

func TestProductUpdateRejectsNegativePrice(t *testing.T) {
	repo := &mockProductRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.Product, error) {
		return &entity.Product{ID: 1}, nil
	}}
	svc := service.NewServiceWithRepo(repo)
	neg := -1.0
	_, err := svc.Update(context.Background(), 1, &dto.UpdateProductRequest{Price: &neg})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Update() error = %v, want ErrBadRequest", err)
	}
}

func TestProductUpdateReturnsErrNotFound(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockProductRepo{})
	_, err := svc.Update(context.Background(), 99, &dto.UpdateProductRequest{})
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestProductDeleteReturnsErrNotFound(t *testing.T) {
	repo := &mockProductRepo{deleteFn: func(ctx context.Context, id int64) error { return apierrors.ErrNotFound }}
	svc := service.NewServiceWithRepo(repo)
	if err := svc.Delete(context.Background(), 99); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}
