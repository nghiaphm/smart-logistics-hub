package service

import (
	"context"
	"fmt"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/workspace/dto"
	"my-web-app.com/smart-logistic-hub/internal/workspace/entity"
)

type WorkspaceRepository interface {
	Create(ctx context.Context, w *entity.Workspace) error
	GetByID(ctx context.Context, id int64) (*entity.Workspace, error)
	GetByCode(ctx context.Context, code string) (*entity.Workspace, error)
	List(ctx context.Context, offset, limit int) ([]entity.Workspace, error)
	Count(ctx context.Context) (int, error)
	Update(ctx context.Context, id int64, w *entity.Workspace) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	repo WorkspaceRepository
}

func NewService(repo WorkspaceRepository) *Service {
	return &Service{repo: repo}
}

func NewServiceWithRepo(repo WorkspaceRepository) *Service {
	return NewService(repo)
}

func (s *Service) Create(ctx context.Context, req *dto.CreateWorkspaceRequest) (*entity.Workspace, error) {
	if req.WorkspaceCode == "" {
		return nil, fmt.Errorf("%w: workspace_code is required", apierrors.ErrBadRequest)
	}
	if req.Name == "" {
		return nil, fmt.Errorf("%w: name is required", apierrors.ErrBadRequest)
	}

	existing, err := s.repo.GetByCode(ctx, req.WorkspaceCode)
	if err != nil && err != apierrors.ErrNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: workspace with code %s already exists", apierrors.ErrConflict, req.WorkspaceCode)
	}

	w := &entity.Workspace{
		WorkspaceCode: req.WorkspaceCode,
		Name:          req.Name,
		IsActive:      true,
	}
	if req.Description != nil {
		w.Description = *req.Description
	}

	if err := s.repo.Create(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*entity.Workspace, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, offset, limit int) ([]entity.Workspace, int, error) {
	workspaces, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return workspaces, count, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *dto.UpdateWorkspaceRequest) (*entity.Workspace, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, id, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
