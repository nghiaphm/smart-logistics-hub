package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
	"my-web-app.com/smart-logistic-hub/internal/tracking/dto"
	"my-web-app.com/smart-logistic-hub/internal/tracking/entity"
	"my-web-app.com/smart-logistic-hub/internal/tracking/service"
)

type mockTrackingRepo struct {
	createFn     func(ctx context.Context, event *entity.TrackingEvent) error
	getByIDFn    func(ctx context.Context, id int64) (*entity.TrackingEvent, error)
	listFn       func(ctx context.Context, orderCode, driverCode string, offset, limit int) ([]entity.TrackingEvent, error)
	countFn      func(ctx context.Context, orderCode, driverCode string) (int, error)
	getByOrderFn func(ctx context.Context, orderCode string) ([]entity.TrackingEvent, error)
	updateFn     func(ctx context.Context, id int64, event *entity.TrackingEvent) error
	deleteFn     func(ctx context.Context, id int64) error
}

func (m *mockTrackingRepo) Create(ctx context.Context, event *entity.TrackingEvent) error {
	if m.createFn != nil {
		return m.createFn(ctx, event)
	}
	return nil
}
func (m *mockTrackingRepo) GetByID(ctx context.Context, id int64) (*entity.TrackingEvent, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, apierrors.ErrNotFound
}
func (m *mockTrackingRepo) List(ctx context.Context, orderCode, driverCode string, offset, limit int) ([]entity.TrackingEvent, error) {
	if m.listFn != nil {
		return m.listFn(ctx, orderCode, driverCode, offset, limit)
	}
	return []entity.TrackingEvent{}, nil
}
func (m *mockTrackingRepo) Count(ctx context.Context, orderCode, driverCode string) (int, error) {
	if m.countFn != nil {
		return m.countFn(ctx, orderCode, driverCode)
	}
	return 0, nil
}
func (m *mockTrackingRepo) GetByOrder(ctx context.Context, orderCode string) ([]entity.TrackingEvent, error) {
	if m.getByOrderFn != nil {
		return m.getByOrderFn(ctx, orderCode)
	}
	return []entity.TrackingEvent{}, nil
}
func (m *mockTrackingRepo) Update(ctx context.Context, id int64, event *entity.TrackingEvent) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, event)
	}
	return nil
}
func (m *mockTrackingRepo) Delete(ctx context.Context, id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func TestTrackingCreateValid(t *testing.T) {
	var saved *entity.TrackingEvent
	repo := &mockTrackingRepo{createFn: func(ctx context.Context, event *entity.TrackingEvent) error { event.ID = 1; saved = event; return nil }}
	svc := service.NewServiceWithRepo(repo)
	event, err := svc.Create(context.Background(), &dto.CreateTrackingEventRequest{OrderCode: "ORD001", DriverCode: "DRV001", StatusUpdate: "PICKED_UP"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if event.ID != 1 {
		t.Errorf("Create() ID = %d, want 1", event.ID)
	}
	if event.Timestamp.IsZero() {
		t.Error("Create() should set Timestamp")
	}
	if saved == nil {
		t.Fatal("Create() did not call repository.Create")
	}
}

func TestTrackingCreateRequiresOrderCode(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockTrackingRepo{})
	_, err := svc.Create(context.Background(), &dto.CreateTrackingEventRequest{DriverCode: "DRV001", StatusUpdate: "PICKED_UP"})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestTrackingCreateRequiresDriverCode(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockTrackingRepo{})
	_, err := svc.Create(context.Background(), &dto.CreateTrackingEventRequest{OrderCode: "ORD001", StatusUpdate: "PICKED_UP"})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestTrackingCreateRequiresStatusUpdate(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockTrackingRepo{})
	_, err := svc.Create(context.Background(), &dto.CreateTrackingEventRequest{OrderCode: "ORD001", DriverCode: "DRV001"})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestTrackingCreatePropagatesRepoError(t *testing.T) {
	repoErr := errors.New("db down")
	repo := &mockTrackingRepo{createFn: func(ctx context.Context, event *entity.TrackingEvent) error { return repoErr }}
	svc := service.NewServiceWithRepo(repo)
	_, err := svc.Create(context.Background(), &dto.CreateTrackingEventRequest{OrderCode: "ORD001", DriverCode: "DRV001", StatusUpdate: "PICKED_UP"})
	if !errors.Is(err, repoErr) {
		t.Errorf("Create() error = %v, want %v", err, repoErr)
	}
}

func TestTrackingGetReturnsErrNotFound(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockTrackingRepo{})
	_, err := svc.Get(context.Background(), 99)
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestTrackingListFiltersByCodes(t *testing.T) {
	var gotOrderCode, gotDriverCode string
	repo := &mockTrackingRepo{
		listFn: func(ctx context.Context, orderCode, driverCode string, offset, limit int) ([]entity.TrackingEvent, error) {
			gotOrderCode = orderCode
			gotDriverCode = driverCode
			return []entity.TrackingEvent{}, nil
		},
		countFn: func(ctx context.Context, orderCode, driverCode string) (int, error) { return 0, nil },
	}
	svc := service.NewServiceWithRepo(repo)
	_, _, err := svc.List(context.Background(), "ORD001", "DRV001", 0, 20)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if gotOrderCode != "ORD001" || gotDriverCode != "DRV001" {
		t.Errorf("List() filters = (%q, %q), want (ORD001, DRV001)", gotOrderCode, gotDriverCode)
	}
}

func TestTrackingGetByOrder(t *testing.T) {
	items := []entity.TrackingEvent{{ID: 1, OrderCode: "ORD001"}}
	repo := &mockTrackingRepo{getByOrderFn: func(ctx context.Context, orderCode string) ([]entity.TrackingEvent, error) { return items, nil }}
	svc := service.NewServiceWithRepo(repo)
	got, err := svc.GetByOrder(context.Background(), "ORD001")
	if err != nil {
		t.Fatalf("GetByOrder() error = %v", err)
	}
	if len(got) != 1 {
		t.Errorf("GetByOrder() len = %d, want 1", len(got))
	}
}

func TestTrackingUpdateRejectsEmptyOrderCode(t *testing.T) {
	repo := &mockTrackingRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.TrackingEvent, error) {
		return &entity.TrackingEvent{ID: 1, OrderCode: "ORD001"}, nil
	}}
	svc := service.NewServiceWithRepo(repo)
	empty := ""
	_, err := svc.Update(context.Background(), 1, &dto.UpdateTrackingEventRequest{OrderCode: &empty})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Update() error = %v, want ErrBadRequest", err)
	}
}

func TestTrackingUpdateAppliesPartialFields(t *testing.T) {
	existing := &entity.TrackingEvent{ID: 1, OrderCode: "ORD001", DriverCode: "DRV001", StatusUpdate: "PICKED_UP"}
	repo := &mockTrackingRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.TrackingEvent, error) { return existing, nil }}
	svc := service.NewServiceWithRepo(repo)
	note := "updated note"
	event, err := svc.Update(context.Background(), 1, &dto.UpdateTrackingEventRequest{Note: &note})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if event.Note != note {
		t.Errorf("Update() Note = %q, want %q", event.Note, note)
	}
	if !event.Timestamp.After(time.Time{}) {
		t.Error("Update() should refresh Timestamp")
	}
}

func TestTrackingUpdateReturnsErrNotFound(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockTrackingRepo{})
	_, err := svc.Update(context.Background(), 99, &dto.UpdateTrackingEventRequest{})
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestTrackingDeleteReturnsErrNotFound(t *testing.T) {
	repo := &mockTrackingRepo{deleteFn: func(ctx context.Context, id int64) error { return apierrors.ErrNotFound }}
	svc := service.NewServiceWithRepo(repo)
	if err := svc.Delete(context.Background(), 99); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}
