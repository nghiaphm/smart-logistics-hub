package service

import (
	"context"
	"fmt"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/product/dto"
	"my-web-app.com/smart-logistic-hub/internal/product/entity"
)

type ProductRepository interface {
	Create(ctx context.Context, p *entity.Product) error
	GetByID(ctx context.Context, id int64) (*entity.Product, error)
	GetBySku(ctx context.Context, sku string) (*entity.Product, error)
	List(ctx context.Context, offset, limit int) ([]entity.Product, error)
	Count(ctx context.Context) (int, error)
	Update(ctx context.Context, id int64, p *entity.Product) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	repo ProductRepository
}

func NewService(repo ProductRepository) *Service {
	return &Service{repo: repo}
}

func NewServiceWithRepo(repo ProductRepository) *Service {
	return NewService(repo)
}

func (s *Service) Create(ctx context.Context, req *dto.CreateProductRequest) (*entity.Product, error) {
	if req.Sku == "" {
		return nil, fmt.Errorf("%w: sku is required", apierrors.ErrBadRequest)
	}
	if req.Name == "" {
		return nil, fmt.Errorf("%w: name is required", apierrors.ErrBadRequest)
	}
	if req.Price < 0 {
		return nil, fmt.Errorf("%w: price must not be negative", apierrors.ErrBadRequest)
	}
	if req.WeightGram < 0 {
		return nil, fmt.Errorf("%w: weight_gram must not be negative", apierrors.ErrBadRequest)
	}

	existing, err := s.repo.GetBySku(ctx, req.Sku)
	if err != nil && err != apierrors.ErrNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: product with sku %s already exists", apierrors.ErrConflict, req.Sku)
	}

	p := &entity.Product{
		Sku:        req.Sku,
		Name:       req.Name,
		Category:   req.Category,
		Price:      req.Price,
		WeightGram: req.WeightGram,
		CreatedBy:  req.CreatedBy,
	}
	if req.Dimensions != nil {
		if req.Dimensions.Length != nil {
			p.LengthCm = *req.Dimensions.Length
		}
		if req.Dimensions.Width != nil {
			p.WidthCm = *req.Dimensions.Width
		}
		if req.Dimensions.Height != nil {
			p.HeightCm = *req.Dimensions.Height
		}
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*entity.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, offset, limit int) ([]entity.Product, int, error) {
	products, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return products, count, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *dto.UpdateProductRequest) (*entity.Product, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Category != nil {
		existing.Category = *req.Category
	}
	if req.Price != nil {
		if *req.Price < 0 {
			return nil, fmt.Errorf("%w: price must not be negative", apierrors.ErrBadRequest)
		}
		existing.Price = *req.Price
	}
	if req.WeightGram != nil {
		if *req.WeightGram < 0 {
			return nil, fmt.Errorf("%w: weight_gram must not be negative", apierrors.ErrBadRequest)
		}
		existing.WeightGram = *req.WeightGram
	}
	if req.Dimensions != nil {
		if req.Dimensions.Length != nil {
			if *req.Dimensions.Length < 0 {
				return nil, fmt.Errorf("%w: length must not be negative", apierrors.ErrBadRequest)
			}
			existing.LengthCm = *req.Dimensions.Length
		}
		if req.Dimensions.Width != nil {
			if *req.Dimensions.Width < 0 {
				return nil, fmt.Errorf("%w: width must not be negative", apierrors.ErrBadRequest)
			}
			existing.WidthCm = *req.Dimensions.Width
		}
		if req.Dimensions.Height != nil {
			if *req.Dimensions.Height < 0 {
				return nil, fmt.Errorf("%w: height must not be negative", apierrors.ErrBadRequest)
			}
			existing.HeightCm = *req.Dimensions.Height
		}
	}
	if req.CreatedBy != "" {
		existing.CreatedBy = req.CreatedBy
	}

	if err := s.repo.Update(ctx, id, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
