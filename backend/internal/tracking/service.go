package tracking

import (
	"context"
	"fmt"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
)

type TrackingRepository interface {
	Create(ctx context.Context, event *TrackingEvent) error
	GetByID(ctx context.Context, id int64) (*TrackingEvent, error)
	List(ctx context.Context, orderCode, driverCode string, offset, limit int) ([]TrackingEvent, error)
	Count(ctx context.Context, orderCode, driverCode string) (int, error)
	GetByOrder(ctx context.Context, orderCode string) ([]TrackingEvent, error)
	Update(ctx context.Context, id int64, event *TrackingEvent) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	Repo TrackingRepository
}

func NewService(repo *Repository) *Service {
	return &Service{Repo: repo}
}

func NewServiceWithRepo(repo TrackingRepository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) Create(ctx context.Context, req *CreateTrackingEventRequest) (*TrackingEvent, error) {
	if req.OrderCode == "" {
		return nil, fmt.Errorf("%w: order_code is required", apierrors.ErrBadRequest)
	}
	if req.DriverCode == "" {
		return nil, fmt.Errorf("%w: driver_code is required", apierrors.ErrBadRequest)
	}
	if req.StatusUpdate == "" {
		return nil, fmt.Errorf("%w: status_update is required", apierrors.ErrBadRequest)
	}

	event := &TrackingEvent{
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

func (s *Service) Get(ctx context.Context, id int64) (*TrackingEvent, error) {
	event, err := s.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *Service) List(ctx context.Context, orderCode, driverCode string, offset, limit int) ([]TrackingEvent, int, error) {
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

func (s *Service) GetByOrder(ctx context.Context, orderCode string) ([]TrackingEvent, error) {
	events, err := s.Repo.GetByOrder(ctx, orderCode)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *UpdateTrackingEventRequest) (*TrackingEvent, error) {
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
	if err := s.Repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
