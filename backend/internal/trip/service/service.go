package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	driverentity "my-web-app.com/smart-logistic-hub/internal/driver/entity"
	"my-web-app.com/smart-logistic-hub/internal/trip/dto"
	"my-web-app.com/smart-logistic-hub/internal/trip/entity"
)

const driverAvailableStatus = "AVAILABLE"

type TripRepository interface {
	Create(ctx context.Context, t *entity.Trip) error
	GetByID(ctx context.Context, id int64) (*entity.Trip, error)
	GetByCode(ctx context.Context, code string) (*entity.Trip, error)
	List(ctx context.Context, offset, limit int) ([]entity.Trip, error)
	Count(ctx context.Context) (int, error)
	Update(ctx context.Context, id int64, t *entity.Trip) error
	Delete(ctx context.Context, id int64) error
	CreateStops(ctx context.Context, tripID int64, stops []entity.TripStop) error
	GetStops(ctx context.Context, tripID int64) ([]entity.TripStop, error)
	DeleteStops(ctx context.Context, tripID int64) error
}

type DriverRepository interface {
	GetByCode(ctx context.Context, code string) (*driverentity.Driver, error)
}

type Service struct {
	repo   TripRepository
	driver DriverRepository
}

func NewService(repo TripRepository, driver DriverRepository) *Service {
	return &Service{repo: repo, driver: driver}
}

func NewServiceWithRepo(repo TripRepository, driver DriverRepository) *Service {
	return NewService(repo, driver)
}

func (s *Service) Create(ctx context.Context, req *dto.CreateTripRequest) (*entity.Trip, error) {
	if req.TripCode == "" {
		return nil, fmt.Errorf("%w: trip_code is required", apierrors.ErrBadRequest)
	}
	if req.DriverCode == "" {
		return nil, fmt.Errorf("%w: driver_code is required", apierrors.ErrBadRequest)
	}
	if len(req.Stops) == 0 {
		return nil, fmt.Errorf("%w: at least one stop is required", apierrors.ErrBadRequest)
	}

	driver, err := s.resolveAvailableDriver(ctx, req.DriverCode)
	if err != nil {
		return nil, err
	}

	licensePlate := driver.LicensePlate
	if req.VehicleLicensePlate != nil && *req.VehicleLicensePlate != "" {
		licensePlate = *req.VehicleLicensePlate
	}

	now := time.Now().UTC()
	t := &entity.Trip{
		TripCode:            req.TripCode,
		DriverID:            &driver.ID,
		VehicleLicensePlate: licensePlate,
		Status:              req.Status,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	if t.Status == "" {
		t.Status = "PLANNED"
	}
	if req.TotalDistanceKm != nil {
		t.TotalDistanceKm = *req.TotalDistanceKm
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	stops, err := buildStops(req.Stops)
	if err != nil {
		return nil, err
	}
	if err := s.repo.CreateStops(ctx, t.ID, stops); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) AssignDriver(ctx context.Context, tripID int64, driverCode string) (*entity.Trip, error) {
	if driverCode == "" {
		return nil, fmt.Errorf("%w: driver_code is required", apierrors.ErrBadRequest)
	}
	trip, err := s.repo.GetByID(ctx, tripID)
	if err != nil {
		return nil, err
	}
	driver, err := s.resolveAvailableDriver(ctx, driverCode)
	if err != nil {
		return nil, err
	}
	trip.DriverID = &driver.ID
	trip.VehicleLicensePlate = driver.LicensePlate
	if err := s.repo.Update(ctx, tripID, trip); err != nil {
		return nil, err
	}
	return trip, nil
}

func (s *Service) resolveAvailableDriver(ctx context.Context, driverCode string) (*driverentity.Driver, error) {
	driver, err := s.driver.GetByCode(ctx, driverCode)
	if err != nil {
		if errors.Is(err, apierrors.ErrNotFound) {
			return nil, fmt.Errorf("%w: driver %s does not exist", apierrors.ErrBadRequest, driverCode)
		}
		return nil, err
	}
	if driver.Status != driverAvailableStatus {
		return nil, fmt.Errorf("%w: driver %s is not available (status %s)", apierrors.ErrConflict, driverCode, driver.Status)
	}
	return driver, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*entity.Trip, []entity.TripStop, error) {
	trip, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	stops, err := s.repo.GetStops(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return trip, stops, nil
}

func (s *Service) List(ctx context.Context, offset, limit int) ([]entity.Trip, int, error) {
	trips, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	return trips, count, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *dto.UpdateTripRequest) (*entity.Trip, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.StartedAt != nil {
		existing.ActualStartAt = req.StartedAt
	}
	if req.CompletedAt != nil {
		existing.ActualEndAt = req.CompletedAt
	}
	if req.Stops != nil {
		stops, err := buildStops(*req.Stops)
		if err != nil {
			return nil, err
		}
		if err := s.repo.DeleteStops(ctx, id); err != nil {
			return nil, err
		}
		if err := s.repo.CreateStops(ctx, id, stops); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Update(ctx, id, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.DeleteStops(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func buildStops(req []dto.TripStopRequest) ([]entity.TripStop, error) {
	stops := make([]entity.TripStop, 0, len(req))
	for _, s := range req {
		if s.OrderCode == "" {
			return nil, fmt.Errorf("%w: order_code is required for each stop", apierrors.ErrBadRequest)
		}
		if s.Address == "" {
			return nil, fmt.Errorf("%w: address is required for each stop", apierrors.ErrBadRequest)
		}
		stop := entity.TripStop{
			OrderCode:   s.OrderCode,
			StopType:    s.StopType,
			Address:     s.Address,
			Status:      s.Status,
			PlannedAt:   s.PlannedAt,
			ArrivedAt:   s.ArrivedAt,
			DepartureAt: s.DepartureAt,
		}
		if stop.StopType == "" {
			stop.StopType = "PICKUP"
		}
		if stop.Status == "" {
			stop.Status = "PENDING"
		}
		if s.Location != nil {
			if s.Location.Lat != nil {
				stop.Lat = *s.Location.Lat
			}
			if s.Location.Lng != nil {
				stop.Lng = *s.Location.Lng
			}
		}
		stops = append(stops, stop)
	}
	return stops, nil
}
