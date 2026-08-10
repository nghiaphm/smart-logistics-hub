package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"my-web-app.com/smart-logistic-hub/internal/billing/dto"
	"my-web-app.com/smart-logistic-hub/internal/billing/entity"
	"my-web-app.com/smart-logistic-hub/internal/billing/service"
	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	orderentity "my-web-app.com/smart-logistic-hub/internal/order/entity"
)

type mockBillingRepo struct {
	createFn         func(ctx context.Context, b *entity.Billing) error
	getByIDFn        func(ctx context.Context, id int64) (*entity.Billing, error)
	getByCodeFn      func(ctx context.Context, code string) (*entity.Billing, error)
	getByOrderCodeFn func(ctx context.Context, orderCode string) (*entity.Billing, error)
	listFn           func(ctx context.Context, offset, limit int) ([]entity.Billing, error)
	countFn          func(ctx context.Context) (int, error)
	updateFn         func(ctx context.Context, id int64, b *entity.Billing) error
	deleteFn         func(ctx context.Context, id int64) error
}

func (m *mockBillingRepo) Create(ctx context.Context, b *entity.Billing) error {
	if m.createFn != nil {
		return m.createFn(ctx, b)
	}
	return nil
}

func (m *mockBillingRepo) GetByID(ctx context.Context, id int64) (*entity.Billing, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockBillingRepo) GetByCode(ctx context.Context, code string) (*entity.Billing, error) {
	if m.getByCodeFn != nil {
		return m.getByCodeFn(ctx, code)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockBillingRepo) GetByOrderCode(ctx context.Context, orderCode string) (*entity.Billing, error) {
	if m.getByOrderCodeFn != nil {
		return m.getByOrderCodeFn(ctx, orderCode)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockBillingRepo) List(ctx context.Context, offset, limit int) ([]entity.Billing, error) {
	if m.listFn != nil {
		return m.listFn(ctx, offset, limit)
	}
	return []entity.Billing{}, nil
}

func (m *mockBillingRepo) Count(ctx context.Context) (int, error) {
	if m.countFn != nil {
		return m.countFn(ctx)
	}
	return 0, nil
}

func (m *mockBillingRepo) Update(ctx context.Context, id int64, b *entity.Billing) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, b)
	}
	return nil
}

func (m *mockBillingRepo) Delete(ctx context.Context, id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

type mockOrderRepo struct {
	getByCodeFn func(ctx context.Context, code string) (*orderentity.Order, error)
}

func (m *mockOrderRepo) GetByCode(ctx context.Context, code string) (*orderentity.Order, error) {
	if m.getByCodeFn != nil {
		return m.getByCodeFn(ctx, code)
	}
	return nil, apierrors.ErrNotFound
}

func validCreateRequest() *dto.CreateBillingRequest {
	return &dto.CreateBillingRequest{
		BillingCode:   "BILL-001",
		OrderCode:     "ORD001",
		AmountTotal:   150000,
		PaymentMethod: "COD",
		PayerInfo: dto.PayerInfo{
			Name:  "Nguyen Van A",
			Phone: "0900000000",
		},
	}
}

func newSvc(billingRepo *mockBillingRepo, orderRepo *mockOrderRepo) *service.Service {
	if billingRepo == nil {
		billingRepo = &mockBillingRepo{}
	}
	if orderRepo == nil {
		orderRepo = &mockOrderRepo{}
	}
	return service.NewServiceWithRepo(billingRepo, orderRepo)
}

func TestBillingCreateValid(t *testing.T) {
	var saved *entity.Billing
	billingRepo := &mockBillingRepo{
		createFn: func(ctx context.Context, b *entity.Billing) error {
			b.ID = 1
			saved = b
			return nil
		},
	}
	orderRepo := &mockOrderRepo{getByCodeFn: func(ctx context.Context, code string) (*orderentity.Order, error) {
		return &orderentity.Order{ID: 1, OrderCode: code}, nil
	}}
	svc := newSvc(billingRepo, orderRepo)

	b, err := svc.Create(context.Background(), validCreateRequest())
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if b.ID != 1 {
		t.Errorf("Create() ID = %d, want 1", b.ID)
	}
	if saved.PaymentStatus != "PENDING" {
		t.Errorf("Create() payment_status = %q, want PENDING", saved.PaymentStatus)
	}
	if saved.Currency != "VND" {
		t.Errorf("Create() currency = %q, want VND", saved.Currency)
	}
	if saved.PayerName != "Nguyen Van A" {
		t.Errorf("Create() payer_name = %q, want %q", saved.PayerName, "Nguyen Van A")
	}
}

func TestBillingCreateSetsPaidAtWhenPaid(t *testing.T) {
	billingRepo := &mockBillingRepo{
		createFn: func(ctx context.Context, b *entity.Billing) error {
			b.ID = 1
			return nil
		},
	}
	orderRepo := &mockOrderRepo{getByCodeFn: func(ctx context.Context, code string) (*orderentity.Order, error) {
		return &orderentity.Order{ID: 1, OrderCode: code}, nil
	}}
	svc := newSvc(billingRepo, orderRepo)

	req := validCreateRequest()
	req.PaymentStatus = "PAID"
	b, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if b.PaidAt == nil {
		t.Error("Create() PaidAt not set for PAID status")
	}
}

func TestBillingCreateRejectsUnknownOrder(t *testing.T) {
	svc := newSvc(&mockBillingRepo{}, &mockOrderRepo{})
	_, err := svc.Create(context.Background(), validCreateRequest())
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestBillingCreateRejectsDuplicatePaid(t *testing.T) {
	billingRepo := &mockBillingRepo{
		getByOrderCodeFn: func(ctx context.Context, orderCode string) (*entity.Billing, error) {
			return &entity.Billing{ID: 1, OrderCode: orderCode, PaymentStatus: "PAID"}, nil
		},
	}
	orderRepo := &mockOrderRepo{getByCodeFn: func(ctx context.Context, code string) (*orderentity.Order, error) {
		return &orderentity.Order{ID: 1, OrderCode: code}, nil
	}}
	svc := newSvc(billingRepo, orderRepo)

	_, err := svc.Create(context.Background(), validCreateRequest())
	if !errors.Is(err, apierrors.ErrConflict) {
		t.Errorf("Create() error = %v, want ErrConflict", err)
	}
}

func TestBillingCreateAllowsPendingFollowUp(t *testing.T) {
	billingRepo := &mockBillingRepo{
		getByOrderCodeFn: func(ctx context.Context, orderCode string) (*entity.Billing, error) {
			return &entity.Billing{ID: 1, OrderCode: orderCode, PaymentStatus: "FAILED"}, nil
		},
		createFn: func(ctx context.Context, b *entity.Billing) error {
			b.ID = 2
			return nil
		},
	}
	orderRepo := &mockOrderRepo{getByCodeFn: func(ctx context.Context, code string) (*orderentity.Order, error) {
		return &orderentity.Order{ID: 1, OrderCode: code}, nil
	}}
	svc := newSvc(billingRepo, orderRepo)

	_, err := svc.Create(context.Background(), validCreateRequest())
	if err != nil {
		t.Fatalf("Create() error = %v, want success for non-PAID previous record", err)
	}
}

func TestBillingCreateRejectsNegativeAmount(t *testing.T) {
	req := validCreateRequest()
	req.AmountTotal = -1
	svc := newSvc(&mockBillingRepo{}, &mockOrderRepo{})
	_, err := svc.Create(context.Background(), req)
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestBillingGetReturnsErrNotFound(t *testing.T) {
	svc := newSvc(&mockBillingRepo{}, &mockOrderRepo{})
	_, err := svc.Get(context.Background(), 99)
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestBillingGetByOrderCode(t *testing.T) {
	expected := &entity.Billing{ID: 1, OrderCode: "ORD001"}
	billingRepo := &mockBillingRepo{
		getByOrderCodeFn: func(ctx context.Context, orderCode string) (*entity.Billing, error) {
			return expected, nil
		},
	}
	svc := newSvc(billingRepo, &mockOrderRepo{})

	got, err := svc.GetByOrderCode(context.Background(), "ORD001")
	if err != nil {
		t.Fatalf("GetByOrderCode() error = %v", err)
	}
	if got != expected {
		t.Errorf("GetByOrderCode() = %v, want %v", got, expected)
	}
}

func TestBillingListReturnsItemsAndCount(t *testing.T) {
	items := []entity.Billing{{ID: 1, BillingCode: "BILL-001"}}
	billingRepo := &mockBillingRepo{
		listFn:  func(ctx context.Context, offset, limit int) ([]entity.Billing, error) { return items, nil },
		countFn: func(ctx context.Context) (int, error) { return 1, nil },
	}
	svc := newSvc(billingRepo, &mockOrderRepo{})

	got, total, err := svc.List(context.Background(), 0, 20)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if total != 1 || len(got) != 1 {
		t.Errorf("List() = (%d items, %d total), want (1, 1)", len(got), total)
	}
}

func TestBillingUpdateToPaidSetsPaidAt(t *testing.T) {
	existing := &entity.Billing{ID: 1, BillingCode: "BILL-001", PaymentStatus: "PENDING"}
	billingRepo := &mockBillingRepo{
		getByIDFn: func(ctx context.Context, id int64) (*entity.Billing, error) { return existing, nil },
	}
	svc := newSvc(billingRepo, &mockOrderRepo{})

	status := "PAID"
	b, err := svc.Update(context.Background(), 1, &dto.UpdateBillingRequest{PaymentStatus: &status})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if b.PaymentStatus != "PAID" {
		t.Errorf("Update() status = %q, want PAID", b.PaymentStatus)
	}
	if b.PaidAt == nil || b.PaidAt.Before(time.Now().Add(-time.Minute)) {
		t.Errorf("Update() PaidAt = %v, want recent time", b.PaidAt)
	}
}

func TestBillingUpdateReturnsErrNotFound(t *testing.T) {
	svc := newSvc(&mockBillingRepo{}, &mockOrderRepo{})
	_, err := svc.Update(context.Background(), 99, &dto.UpdateBillingRequest{})
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestBillingDeleteReturnsErrNotFound(t *testing.T) {
	billingRepo := &mockBillingRepo{deleteFn: func(ctx context.Context, id int64) error { return apierrors.ErrNotFound }}
	svc := newSvc(billingRepo, &mockOrderRepo{})
	if err := svc.Delete(context.Background(), 99); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}
