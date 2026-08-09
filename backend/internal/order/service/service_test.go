package service_test

import (
	"context"
	"errors"
	"testing"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/order/dto"
	"my-web-app.com/smart-logistic-hub/internal/order/entity"
	"my-web-app.com/smart-logistic-hub/internal/order/service"
)

type mockOrderRepo struct {
	createFn      func(ctx context.Context, o *entity.Order) error
	getByIDFn     func(ctx context.Context, id int64) (*entity.Order, error)
	listFn        func(ctx context.Context, offset, limit int) ([]entity.Order, error)
	countFn       func(ctx context.Context) (int, error)
	updateFn      func(ctx context.Context, id int64, o *entity.Order) error
	deleteFn      func(ctx context.Context, id int64) error
	createItemFn  func(ctx context.Context, item *entity.OrderItem) error
	getItemsFn    func(ctx context.Context, orderID int64) ([]entity.OrderItem, error)
	deleteItemsFn func(ctx context.Context, orderID int64) error
}

func (m *mockOrderRepo) Create(ctx context.Context, o *entity.Order) error {
	if m.createFn != nil {
		return m.createFn(ctx, o)
	}
	return nil
}

func (m *mockOrderRepo) GetByID(ctx context.Context, id int64) (*entity.Order, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockOrderRepo) GetByCode(ctx context.Context, code string) (*entity.Order, error) {
	return nil, apierrors.ErrNotFound
}

func (m *mockOrderRepo) List(ctx context.Context, offset, limit int) ([]entity.Order, error) {
	if m.listFn != nil {
		return m.listFn(ctx, offset, limit)
	}
	return []entity.Order{}, nil
}

func (m *mockOrderRepo) Count(ctx context.Context) (int, error) {
	if m.countFn != nil {
		return m.countFn(ctx)
	}
	return 0, nil
}

func (m *mockOrderRepo) Update(ctx context.Context, id int64, o *entity.Order) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, o)
	}
	return nil
}

func (m *mockOrderRepo) Delete(ctx context.Context, id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockOrderRepo) CreateItem(ctx context.Context, item *entity.OrderItem) error {
	if m.createItemFn != nil {
		return m.createItemFn(ctx, item)
	}
	return nil
}

func (m *mockOrderRepo) GetItemsByOrder(ctx context.Context, orderID int64) ([]entity.OrderItem, error) {
	if m.getItemsFn != nil {
		return m.getItemsFn(ctx, orderID)
	}
	return []entity.OrderItem{}, nil
}

func (m *mockOrderRepo) DeleteItems(ctx context.Context, orderID int64) error {
	if m.deleteItemsFn != nil {
		return m.deleteItemsFn(ctx, orderID)
	}
	return nil
}

func TestOrderCreateDefaultsStatusToPending(t *testing.T) {
	var saved *entity.Order
	var createdItems []*entity.OrderItem
	repo := &mockOrderRepo{
		createFn: func(ctx context.Context, o *entity.Order) error {
			o.ID = 1
			saved = o
			return nil
		},
		createItemFn: func(ctx context.Context, item *entity.OrderItem) error {
			createdItems = append(createdItems, item)
			return nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	o, err := svc.Create(context.Background(), &dto.CreateOrderRequest{
		OrderCode:   "ORD001",
		SenderName:  "Sender",
		SenderPhone: "0900000000",
		Items: []dto.OrderItemRequest{
			{ProductName: "Product A", Quantity: 2},
		},
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if o.Status != "PENDING" {
		t.Errorf("Create() status = %q, want %q", o.Status, "PENDING")
	}
	if saved == nil {
		t.Fatal("Create() did not call repository.Create")
	}
	if len(createdItems) != 1 {
		t.Errorf("Create() created %d items, want 1", len(createdItems))
	}
	if createdItems[0].OrderID != 1 {
		t.Errorf("Create() item.OrderID = %d, want 1", createdItems[0].OrderID)
	}
}

func TestOrderCreateKeepsProvidedStatus(t *testing.T) {
	repo := &mockOrderRepo{
		createFn: func(ctx context.Context, o *entity.Order) error {
			o.ID = 1
			return nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	status := "SHIPPING"
	o, err := svc.Create(context.Background(), &dto.CreateOrderRequest{
		OrderCode: "ORD002",
		Status:    status,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if o.Status != status {
		t.Errorf("Create() status = %q, want %q", o.Status, status)
	}
}

func TestOrderCreatePropagatesRepositoryError(t *testing.T) {
	repoErr := errors.New("database down")
	repo := &mockOrderRepo{
		createFn: func(ctx context.Context, o *entity.Order) error {
			return repoErr
		},
	}
	svc := service.NewServiceWithRepo(repo)

	_, err := svc.Create(context.Background(), &dto.CreateOrderRequest{OrderCode: "ORD003"})
	if !errors.Is(err, repoErr) {
		t.Errorf("Create() error = %v, want %v", err, repoErr)
	}
}

func TestOrderGetReturnsErrNotFound(t *testing.T) {
	repo := &mockOrderRepo{}
	svc := service.NewServiceWithRepo(repo)

	_, err := svc.Get(context.Background(), 99)
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestOrderGetItems(t *testing.T) {
	items := []entity.OrderItem{{ID: 1, OrderID: 1, ProductName: "P1", Quantity: 1}}
	repo := &mockOrderRepo{
		getItemsFn: func(ctx context.Context, orderID int64) ([]entity.OrderItem, error) {
			return items, nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	got, err := svc.GetItems(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetItems() error = %v", err)
	}
	if len(got) != 1 {
		t.Errorf("GetItems() len = %d, want 1", len(got))
	}
}

func TestOrderListReturnsItemsAndCount(t *testing.T) {
	items := []entity.Order{{ID: 1, OrderCode: "ORD001"}}
	repo := &mockOrderRepo{
		listFn: func(ctx context.Context, offset, limit int) ([]entity.Order, error) {
			return items, nil
		},
		countFn: func(ctx context.Context) (int, error) {
			return 1, nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	got, total, err := svc.List(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if total != 1 || len(got) != 1 {
		t.Errorf("List() = (%d items, %d total), want (1, 1)", len(got), total)
	}
}

func TestOrderUpdateAppliesPartialFields(t *testing.T) {
	existing := &entity.Order{ID: 1, OrderCode: "ORD001", Status: "PENDING", ReceiverName: "Old"}
	repo := &mockOrderRepo{
		getByIDFn: func(ctx context.Context, id int64) (*entity.Order, error) {
			return existing, nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	newStatus := "COMPLETED"
	o, err := svc.Update(context.Background(), 1, &dto.UpdateOrderRequest{
		Status: &newStatus,
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if o.Status != newStatus {
		t.Errorf("Update() Status = %q, want %q", o.Status, newStatus)
	}
	if o.ReceiverName != "Old" {
		t.Errorf("Update() ReceiverName = %q, want unchanged %q", o.ReceiverName, "Old")
	}
}

func TestOrderUpdateReturnsErrNotFound(t *testing.T) {
	repo := &mockOrderRepo{}
	svc := service.NewServiceWithRepo(repo)

	_, err := svc.Update(context.Background(), 99, &dto.UpdateOrderRequest{})
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestOrderDeleteDeletesItemsThenOrder(t *testing.T) {
	var deletedItemsID, deletedID int64
	order := 1
	repo := &mockOrderRepo{
		deleteItemsFn: func(ctx context.Context, orderID int64) error {
			deletedItemsID = orderID
			return nil
		},
		deleteFn: func(ctx context.Context, id int64) error {
			deletedID = id
			return nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	if err := svc.Delete(context.Background(), int64(order)); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if deletedItemsID != 1 || deletedID != 1 {
		t.Errorf("Delete() itemsID=%d orderID=%d, want 1,1", deletedItemsID, deletedID)
	}
}

func TestOrderDeleteReturnsErrNotFound(t *testing.T) {
	repo := &mockOrderRepo{
		deleteFn: func(ctx context.Context, id int64) error {
			return apierrors.ErrNotFound
		},
	}
	svc := service.NewServiceWithRepo(repo)

	if err := svc.Delete(context.Background(), 99); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}
