package service_test

import (
	"context"
	"errors"
	"testing"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	driverentity "my-web-app.com/smart-logistic-hub/internal/driver/entity"
	"my-web-app.com/smart-logistic-hub/internal/trip/dto"
	"my-web-app.com/smart-logistic-hub/internal/trip/entity"
	"my-web-app.com/smart-logistic-hub/internal/trip/service"
)

type mockTripRepo struct {
	createFn      func(ctx context.Context, t *entity.Trip) error
	getByIDFn     func(ctx context.Context, id int64) (*entity.Trip, error)
	listFn        func(ctx context.Context, offset, limit int) ([]entity.Trip, error)
	countFn       func(ctx context.Context) (int, error)
	updateFn      func(ctx context.Context, id int64, t *entity.Trip) error
	deleteFn      func(ctx context.Context, id int64) error
	createStopsFn func(ctx context.Context, tripID int64, stops []entity.TripStop) error
	getStopsFn    func(ctx context.Context, tripID int64) ([]entity.TripStop, error)
	deleteStopsFn func(ctx context.Context, tripID int64) error
}

func (m *mockTripRepo) Create(ctx context.Context, t *entity.Trip) error {
	if m.createFn != nil {
		return m.createFn(ctx, t)
	}
	return nil
}

func (m *mockTripRepo) GetByID(ctx context.Context, id int64) (*entity.Trip, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockTripRepo) GetByCode(ctx context.Context, code string) (*entity.Trip, error) {
	return nil, apierrors.ErrNotFound
}

func (m *mockTripRepo) List(ctx context.Context, offset, limit int) ([]entity.Trip, error) {
	if m.listFn != nil {
		return m.listFn(ctx, offset, limit)
	}
	return []entity.Trip{}, nil
}

func (m *mockTripRepo) Count(ctx context.Context) (int, error) {
	if m.countFn != nil {
		return m.countFn(ctx)
	}
	return 0, nil
}

func (m *mockTripRepo) Update(ctx context.Context, id int64, t *entity.Trip) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, t)
	}
	return nil
}

func (m *mockTripRepo) Delete(ctx context.Context, id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func (m *mockTripRepo) CreateStops(ctx context.Context, tripID int64, stops []entity.TripStop) error {
	if m.createStopsFn != nil {
		return m.createStopsFn(ctx, tripID, stops)
	}
	return nil
}

func (m *mockTripRepo) GetStops(ctx context.Context, tripID int64) ([]entity.TripStop, error) {
	if m.getStopsFn != nil {
		return m.getStopsFn(ctx, tripID)
	}
	return []entity.TripStop{}, nil
}

func (m *mockTripRepo) DeleteStops(ctx context.Context, tripID int64) error {
	if m.deleteStopsFn != nil {
		return m.deleteStopsFn(ctx, tripID)
	}
	return nil
}

type mockDriverRepo struct {
	getByCodeFn func(ctx context.Context, code string) (*driverentity.Driver, error)
}

func (m *mockDriverRepo) GetByCode(ctx context.Context, code string) (*driverentity.Driver, error) {
	if m.getByCodeFn != nil {
		return m.getByCodeFn(ctx, code)
	}
	return nil, apierrors.ErrNotFound
}

func availableDriver() *driverentity.Driver {
	return &driverentity.Driver{ID: 7, DriverCode: "DRV001", Status: "AVAILABLE", LicensePlate: "51A-12345"}
}

func validCreateRequest() *dto.CreateTripRequest {
	return &dto.CreateTripRequest{
		TripCode:   "TRIP-001",
		DriverCode: "DRV001",
		Stops: []dto.TripStopRequest{
			{OrderCode: "ORD001", Address: "1 Nguyen Hue"},
		},
	}
}

func TestTripCreateValidAssignsDriverAndStops(t *testing.T) {
	var saved *entity.Trip
	var createdStops []entity.TripStop
	tripRepo := &mockTripRepo{
		createFn: func(ctx context.Context, t *entity.Trip) error {
			t.ID = 1
			saved = t
			return nil
		},
		createStopsFn: func(ctx context.Context, tripID int64, stops []entity.TripStop) error {
			createdStops = stops
			return nil
		},
	}
	driverRepo := &mockDriverRepo{getByCodeFn: func(ctx context.Context, code string) (*driverentity.Driver, error) {
		return availableDriver(), nil
	}}
	svc := service.NewServiceWithRepo(tripRepo, driverRepo)

	trip, err := svc.Create(context.Background(), validCreateRequest())
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if trip.ID != 1 {
		t.Errorf("Create() ID = %d, want 1", trip.ID)
	}
	if saved == nil || saved.DriverID == nil || *saved.DriverID != 7 {
		t.Errorf("Create() driver_id = %v, want 7", saved.DriverID)
	}
	if saved.VehicleLicensePlate != "51A-12345" {
		t.Errorf("Create() license plate = %q, want inherited from driver %q", saved.VehicleLicensePlate, "51A-12345")
	}
	if saved.Status != "PLANNED" {
		t.Errorf("Create() status = %q, want %q", saved.Status, "PLANNED")
	}
	if len(createdStops) != 1 {
		t.Fatalf("Create() created %d stops, want 1", len(createdStops))
	}
	if createdStops[0].StopType != "PICKUP" || createdStops[0].Status != "PENDING" {
		t.Errorf("Create() stop defaults = (%q, %q), want (PICKUP, PENDING)", createdStops[0].StopType, createdStops[0].Status)
	}
}

func TestTripCreateRequiresTripCode(t *testing.T) {
	req := validCreateRequest()
	req.TripCode = ""
	svc := service.NewServiceWithRepo(&mockTripRepo{}, &mockDriverRepo{})
	_, err := svc.Create(context.Background(), req)
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestTripCreateRequiresDriverCode(t *testing.T) {
	req := validCreateRequest()
	req.DriverCode = ""
	svc := service.NewServiceWithRepo(&mockTripRepo{}, &mockDriverRepo{})
	_, err := svc.Create(context.Background(), req)
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestTripCreateRejectsUnknownDriver(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockTripRepo{}, &mockDriverRepo{})
	_, err := svc.Create(context.Background(), validCreateRequest())
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestTripCreateRejectsUnavailableDriver(t *testing.T) {
	driverRepo := &mockDriverRepo{getByCodeFn: func(ctx context.Context, code string) (*driverentity.Driver, error) {
		return &driverentity.Driver{ID: 7, DriverCode: code, Status: "BUSY"}, nil
	}}
	svc := service.NewServiceWithRepo(&mockTripRepo{}, driverRepo)
	_, err := svc.Create(context.Background(), validCreateRequest())
	if !errors.Is(err, apierrors.ErrConflict) {
		t.Errorf("Create() error = %v, want ErrConflict", err)
	}
}

func TestTripCreatePropagatesRepoError(t *testing.T) {
	repoErr := errors.New("db down")
	tripRepo := &mockTripRepo{createFn: func(ctx context.Context, t *entity.Trip) error { return repoErr }}
	driverRepo := &mockDriverRepo{getByCodeFn: func(ctx context.Context, code string) (*driverentity.Driver, error) {
		return availableDriver(), nil
	}}
	svc := service.NewServiceWithRepo(tripRepo, driverRepo)
	_, err := svc.Create(context.Background(), validCreateRequest())
	if !errors.Is(err, repoErr) {
		t.Errorf("Create() error = %v, want %v", err, repoErr)
	}
}

func TestTripAssignDriverSuccess(t *testing.T) {
	tripRepo := &mockTripRepo{
		getByIDFn: func(ctx context.Context, id int64) (*entity.Trip, error) {
			return &entity.Trip{ID: 1, TripCode: "TRIP-001", Status: "PLANNED"}, nil
		},
	}
	driverRepo := &mockDriverRepo{getByCodeFn: func(ctx context.Context, code string) (*driverentity.Driver, error) {
		return availableDriver(), nil
	}}
	svc := service.NewServiceWithRepo(tripRepo, driverRepo)

	trip, err := svc.AssignDriver(context.Background(), 1, "DRV001")
	if err != nil {
		t.Fatalf("AssignDriver() error = %v", err)
	}
	if trip.DriverID == nil || *trip.DriverID != 7 {
		t.Errorf("AssignDriver() driver_id = %v, want 7", trip.DriverID)
	}
}

func TestTripAssignDriverRejectsUnavailableDriver(t *testing.T) {
	tripRepo := &mockTripRepo{
		getByIDFn: func(ctx context.Context, id int64) (*entity.Trip, error) {
			return &entity.Trip{ID: 1, TripCode: "TRIP-001", Status: "PLANNED"}, nil
		},
	}
	driverRepo := &mockDriverRepo{getByCodeFn: func(ctx context.Context, code string) (*driverentity.Driver, error) {
		return &driverentity.Driver{ID: 7, DriverCode: code, Status: "OFFLINE"}, nil
	}}
	svc := service.NewServiceWithRepo(tripRepo, driverRepo)

	_, err := svc.AssignDriver(context.Background(), 1, "DRV001")
	if !errors.Is(err, apierrors.ErrConflict) {
		t.Errorf("AssignDriver() error = %v, want ErrConflict", err)
	}
}

func TestTripAssignDriverReturnsErrNotFoundForTrip(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockTripRepo{}, &mockDriverRepo{})
	_, err := svc.AssignDriver(context.Background(), 99, "DRV001")
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("AssignDriver() error = %v, want ErrNotFound", err)
	}
}

func TestTripGetReturnsStops(t *testing.T) {
	stops := []entity.TripStop{{ID: 1, TripID: 1, OrderCode: "ORD001"}}
	tripRepo := &mockTripRepo{
		getByIDFn: func(ctx context.Context, id int64) (*entity.Trip, error) {
			return &entity.Trip{ID: 1, TripCode: "TRIP-001"}, nil
		},
		getStopsFn: func(ctx context.Context, tripID int64) ([]entity.TripStop, error) {
			return stops, nil
		},
	}
	svc := service.NewServiceWithRepo(tripRepo, &mockDriverRepo{})

	_, got, err := svc.Get(context.Background(), 1)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if len(got) != 1 {
		t.Errorf("Get() stops len = %d, want 1", len(got))
	}
}

func TestTripGetReturnsErrNotFound(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockTripRepo{}, &mockDriverRepo{})
	_, _, err := svc.Get(context.Background(), 99)
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestTripListReturnsItemsAndCount(t *testing.T) {
	items := []entity.Trip{{ID: 1, TripCode: "TRIP-001"}}
	tripRepo := &mockTripRepo{
		listFn:  func(ctx context.Context, offset, limit int) ([]entity.Trip, error) { return items, nil },
		countFn: func(ctx context.Context) (int, error) { return 1, nil },
	}
	svc := service.NewServiceWithRepo(tripRepo, &mockDriverRepo{})

	got, total, err := svc.List(context.Background(), 0, 20)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if total != 1 || len(got) != 1 {
		t.Errorf("List() = (%d items, %d total), want (1, 1)", len(got), total)
	}
}

func TestTripUpdateAppliesStatusAndReplacesStops(t *testing.T) {
	var deletedStops, createdStops []int64
	tripRepo := &mockTripRepo{
		getByIDFn: func(ctx context.Context, id int64) (*entity.Trip, error) {
			return &entity.Trip{ID: 1, TripCode: "TRIP-001", Status: "PLANNED"}, nil
		},
		deleteStopsFn: func(ctx context.Context, tripID int64) error {
			deletedStops = append(deletedStops, tripID)
			return nil
		},
		createStopsFn: func(ctx context.Context, tripID int64, stops []entity.TripStop) error {
			createdStops = append(createdStops, tripID)
			return nil
		},
	}
	svc := service.NewServiceWithRepo(tripRepo, &mockDriverRepo{})

	newStatus := "IN_TRANSIT"
	trip, err := svc.Update(context.Background(), 1, &dto.UpdateTripRequest{
		Status: &newStatus,
		Stops: &[]dto.TripStopRequest{
			{OrderCode: "ORD002", Address: "2 Ly Tu Trong"},
		},
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if trip.Status != newStatus {
		t.Errorf("Update() Status = %q, want %q", trip.Status, newStatus)
	}
	if len(deletedStops) != 1 || len(createdStops) != 1 {
		t.Errorf("Update() stops replaced = (del %d, create %d), want (1, 1)", len(deletedStops), len(createdStops))
	}
}

func TestTripUpdateReturnsErrNotFound(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockTripRepo{}, &mockDriverRepo{})
	_, err := svc.Update(context.Background(), 99, &dto.UpdateTripRequest{})
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestTripDeleteReturnsErrNotFound(t *testing.T) {
	tripRepo := &mockTripRepo{deleteFn: func(ctx context.Context, id int64) error { return apierrors.ErrNotFound }}
	svc := service.NewServiceWithRepo(tripRepo, &mockDriverRepo{})
	if err := svc.Delete(context.Background(), 99); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}
