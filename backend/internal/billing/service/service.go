package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"my-web-app.com/smart-logistic-hub/internal/billing/dto"
	"my-web-app.com/smart-logistic-hub/internal/billing/entity"
	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	orderentity "my-web-app.com/smart-logistic-hub/internal/order/entity"
)

type BillingRepository interface {
	Create(ctx context.Context, b *entity.Billing) error
	GetByID(ctx context.Context, id int64) (*entity.Billing, error)
	GetByCode(ctx context.Context, code string) (*entity.Billing, error)
	GetByOrderCode(ctx context.Context, orderCode string) (*entity.Billing, error)
	List(ctx context.Context, offset, limit int) ([]entity.Billing, error)
	Count(ctx context.Context) (int, error)
	Update(ctx context.Context, id int64, b *entity.Billing) error
	Delete(ctx context.Context, id int64) error
}

type OrderRepository interface {
	GetByCode(ctx context.Context, code string) (*orderentity.Order, error)
}

type Service struct {
	repo   BillingRepository
	orders OrderRepository
}

func NewService(repo BillingRepository, orders OrderRepository) *Service {
	return &Service{repo: repo, orders: orders}
}

func NewServiceWithRepo(repo BillingRepository, orders OrderRepository) *Service {
	return NewService(repo, orders)
}

func (s *Service) Create(ctx context.Context, req *dto.CreateBillingRequest) (*entity.Billing, error) {
	if req.BillingCode == "" {
		return nil, fmt.Errorf("%w: billing_code is required", apierrors.ErrBadRequest)
	}
	if req.OrderCode == "" {
		return nil, fmt.Errorf("%w: order_code is required", apierrors.ErrBadRequest)
	}
	if req.AmountTotal < 0 {
		return nil, fmt.Errorf("%w: amount_total must not be negative", apierrors.ErrBadRequest)
	}

	if _, err := s.orders.GetByCode(ctx, req.OrderCode); err != nil {
		if errors.Is(err, apierrors.ErrNotFound) {
			return nil, fmt.Errorf("%w: order %s does not exist", apierrors.ErrBadRequest, req.OrderCode)
		}
		return nil, err
	}

	existing, err := s.repo.GetByOrderCode(ctx, req.OrderCode)
	if err != nil && err != apierrors.ErrNotFound {
		return nil, err
	}
	if existing != nil && existing.PaymentStatus == "PAID" {
		return nil, fmt.Errorf("%w: order %s already has a paid billing record", apierrors.ErrConflict, req.OrderCode)
	}

	existingByCode, err := s.repo.GetByCode(ctx, req.BillingCode)
	if err != nil && err != apierrors.ErrNotFound {
		return nil, err
	}
	if existingByCode != nil {
		return nil, fmt.Errorf("%w: billing_code %s already exists", apierrors.ErrConflict, req.BillingCode)
	}

	b := &entity.Billing{
		BillingCode:   req.BillingCode,
		OrderCode:     req.OrderCode,
		AmountTotal:   req.AmountTotal,
		Currency:      req.Currency,
		PaymentMethod: req.PaymentMethod,
		PaymentStatus: req.PaymentStatus,
		PayerName:     req.PayerInfo.Name,
		PayerPhone:    req.PayerInfo.Phone,
		CreatedBy:     req.CreatedBy,
	}
	if req.PayerInfo.Email != nil {
		b.PayerEmail = *req.PayerInfo.Email
	}
	if b.Currency == "" {
		b.Currency = "VND"
	}
	if b.PaymentMethod == "" {
		b.PaymentMethod = "COD"
	}
	if b.PaymentStatus == "" {
		b.PaymentStatus = "PENDING"
	}
	if b.PaymentStatus == "PAID" {
		now := time.Now().UTC()
		b.PaidAt = &now
	}

	if err := s.repo.Create(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*entity.Billing, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByCode(ctx context.Context, billingCode string) (*entity.Billing, error) {
	return s.repo.GetByCode(ctx, billingCode)
}

func (s *Service) GetByOrderCode(ctx context.Context, orderCode string) (*entity.Billing, error) {
	return s.repo.GetByOrderCode(ctx, orderCode)
}

func (s *Service) List(ctx context.Context, offset, limit int) ([]entity.Billing, int, error) {
	billings, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return billings, count, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *dto.UpdateBillingRequest) (*entity.Billing, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.PaymentStatus != nil {
		old := existing.PaymentStatus
		existing.PaymentStatus = *req.PaymentStatus
		if *req.PaymentStatus == "PAID" && old != "PAID" {
			now := time.Now().UTC()
			existing.PaidAt = &now
		}
		if *req.PaymentStatus != "PAID" && old == "PAID" {
			existing.PaidAt = nil
		}
	}
	if req.TransactionID != nil {
		existing.TransactionID = *req.TransactionID
	}
	if req.PaidAt != nil {
		existing.PaidAt = req.PaidAt
	}

	if err := s.repo.Update(ctx, id, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
