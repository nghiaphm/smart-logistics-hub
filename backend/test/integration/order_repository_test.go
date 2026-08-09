package integration

import (
	"context"
	"errors"
	"testing"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/order/entity"
	"my-web-app.com/smart-logistic-hub/internal/order/repository"
)

func TestOrderRepositoryCRUD(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()
	repo := repository.NewRepository(testDB)

	order := &entity.Order{
		OrderCode:       "ORD-IT-001",
		SenderName:      "Sender",
		SenderPhone:     "0900000001",
		SenderAddress:   "123 Street",
		ReceiverName:    "Receiver",
		ReceiverPhone:   "0900000002",
		ReceiverAddress: "456 Avenue",
		Status:          "PENDING",
	}
	if err := repo.Create(ctx, order); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if order.ID == 0 {
		t.Fatal("Create() did not set ID")
	}

	item := &entity.OrderItem{
		OrderID:     order.ID,
		ProductName: "Product A",
		Quantity:    2,
		WeightGram:  500,
	}
	if err := repo.CreateItem(ctx, item); err != nil {
		t.Fatalf("CreateItem() error = %v", err)
	}
	if item.ID == 0 {
		t.Fatal("CreateItem() did not set ID")
	}

	items, err := repo.GetItemsByOrder(ctx, order.ID)
	if err != nil {
		t.Fatalf("GetItemsByOrder() error = %v", err)
	}
	if len(items) != 1 {
		t.Errorf("GetItemsByOrder() len = %d, want 1", len(items))
	}

	got, err := repo.GetByID(ctx, order.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.OrderCode != order.OrderCode {
		t.Errorf("GetByID() OrderCode = %q, want %q", got.OrderCode, order.OrderCode)
	}

	got, err = repo.GetByCode(ctx, "ORD-IT-001")
	if err != nil {
		t.Fatalf("GetByCode() error = %v", err)
	}
	if got.ID != order.ID {
		t.Errorf("GetByCode() ID = %d, want %d", got.ID, order.ID)
	}

	order.Status = "SHIPPING"
	if err := repo.Update(ctx, order.ID, order); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	got, err = repo.GetByID(ctx, order.ID)
	if err != nil {
		t.Fatalf("GetByID() after update error = %v", err)
	}
	if got.Status != "SHIPPING" {
		t.Errorf("Update() Status = %q, want %q", got.Status, "SHIPPING")
	}

	if err := repo.DeleteItems(ctx, order.ID); err != nil {
		t.Fatalf("DeleteItems() error = %v", err)
	}
	if err := repo.Delete(ctx, order.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := repo.GetByID(ctx, order.ID); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("GetByID() after delete error = %v, want ErrNotFound", err)
	}
}

func TestOrderRepositoryNotFound(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()
	repo := repository.NewRepository(testDB)

	if _, err := repo.GetByID(ctx, 999999); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("GetByID() error = %v, want ErrNotFound", err)
	}
	if _, err := repo.GetByCode(ctx, "NOPE"); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("GetByCode() error = %v, want ErrNotFound", err)
	}
	if err := repo.Update(ctx, 999999, &entity.Order{}); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
	if err := repo.Delete(ctx, 999999); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}

func TestOrderRepositoryListAndCount(t *testing.T) {
	truncateTables(t)
	ctx := context.Background()
	repo := repository.NewRepository(testDB)

	for i := 1; i <= 3; i++ {
		o := &entity.Order{
			OrderCode:       "ORD-IT-LST" + string(rune('0'+i)),
			SenderName:      "Sender",
			SenderPhone:     "0900000001",
			SenderAddress:   "Addr",
			ReceiverName:    "Receiver",
			ReceiverPhone:   "0900000002",
			ReceiverAddress: "Addr2",
			Status:          "PENDING",
		}
		if err := repo.Create(ctx, o); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	orders, err := repo.List(ctx, 0, 10)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(orders) != 3 {
		t.Errorf("List() len = %d, want 3", len(orders))
	}
	count, err := repo.Count(ctx)
	if err != nil {
		t.Fatalf("Count() error = %v", err)
	}
	if count != 3 {
		t.Errorf("Count() = %d, want 3", count)
	}
}
