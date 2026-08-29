package service

import (
	"context"
	"fmt"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/tracking/dto"
	"my-web-app.com/smart-logistic-hub/internal/tracking/entity"
	trkrepo "my-web-app.com/smart-logistic-hub/internal/tracking/repository"
)

type TrackingRepository interface {
	Create(ctx context.Context, event *entity.TrackingEvent) error
	GetByID(ctx context.Context, id int64) (*entity.TrackingEvent, error)
	List(ctx context.Context, orderCode, driverCode string, offset, limit int) ([]entity.TrackingEvent, error)
	Count(ctx context.Context, orderCode, driverCode string) (int, error)
	GetByOrder(ctx context.Context, orderCode string) ([]entity.TrackingEvent, error)
	Update(ctx context.Context, id int64, event *entity.TrackingEvent) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	Repo TrackingRepository
}

func NewService(repo *trkrepo.Repository) *Service {
	return &Service{Repo: repo}
}

func NewServiceWithRepo(repo TrackingRepository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) Create(ctx context.Context, req *dto.CreateTrackingEventRequest) (*entity.TrackingEvent, error) {
	if req.OrderCode == "" {
		return nil, fmt.Errorf("%w: order_code is required", apierrors.ErrBadRequest)
	}
	if req.DriverCode == "" {
		return nil, fmt.Errorf("%w: driver_code is required", apierrors.ErrBadRequest)
	}
	if req.StatusUpdate == "" {
		return nil, fmt.Errorf("%w: status_update is required", apierrors.ErrBadRequest)
	}
	event := &entity.TrackingEvent{
		OrderID:      req.OrderID,
		OrderCode:    req.OrderCode,
		DriverCode:   req.DriverCode,
		StatusUpdate: req.StatusUpdate,
		Lat:          req.Lat,
		Lng:          req.Lng,
		Note:         req.Note,
		Timestamp:    time.Now().UTC(),
	}
	if err := s.Repo.Create(ctx, event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*entity.TrackingEvent, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, orderCode, driverCode string, offset, limit int) ([]entity.TrackingEvent, int, error) {
	events, err := s.Repo.List(ctx, orderCode, driverCode, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.Repo.Count(ctx, orderCode, driverCode)
	if err != nil {
		return nil, 0, err
	}
	return events, count, nil
}

func (s *Service) GetByOrder(ctx context.Context, orderCode string) ([]entity.TrackingEvent, error) {
	return s.Repo.GetByOrder(ctx, orderCode)
}

func (s *Service) Update(ctx context.Context, id int64, req *dto.UpdateTrackingEventRequest) (*entity.TrackingEvent, error) {
	event, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.OrderCode != nil {
		if *req.OrderCode == "" {
			return nil, fmt.Errorf("%w: order_code must not be empty", apierrors.ErrBadRequest)
		}
		event.OrderCode = *req.OrderCode
	}
	if req.OrderID != nil {
		event.OrderID = req.OrderID
	}
	if req.DriverCode != nil {
		if *req.DriverCode == "" {
			return nil, fmt.Errorf("%w: driver_code must not be empty", apierrors.ErrBadRequest)
		}
		event.DriverCode = *req.DriverCode
	}
	if req.StatusUpdate != nil {
		if *req.StatusUpdate == "" {
			return nil, fmt.Errorf("%w: status_update must not be empty", apierrors.ErrBadRequest)
		}
		event.StatusUpdate = *req.StatusUpdate
	}
	if req.Lat != nil {
		event.Lat = *req.Lat
	}
	if req.Lng != nil {
		event.Lng = *req.Lng
	}
	if req.Note != nil {
		event.Note = *req.Note
	}
	event.Timestamp = time.Now().UTC()
	if err := s.Repo.Update(ctx, id, event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.Repo.Delete(ctx, id)
}
