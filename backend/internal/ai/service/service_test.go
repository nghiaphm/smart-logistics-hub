package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"my-web-app.com/smart-logistic-hub/internal/ai/dto"
	"my-web-app.com/smart-logistic-hub/internal/ai/entity"
	"my-web-app.com/smart-logistic-hub/internal/ai/service"
	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
)

type mockAIRepo struct {
	createFn    func(ctx context.Context, e *entity.AIEvent) error
	getByIDFn   func(ctx context.Context, id int64) (*entity.AIEvent, error)
	getByCodeFn func(ctx context.Context, code string) (*entity.AIEvent, error)
	listFn      func(ctx context.Context, licensePlate, gateID, eventType string, offset, limit int) ([]entity.AIEvent, error)
	countFn     func(ctx context.Context, licensePlate, gateID, eventType string) (int, error)
	updateFn    func(ctx context.Context, id int64, e *entity.AIEvent) error
	deleteFn    func(ctx context.Context, id int64) error
}

func (m *mockAIRepo) Create(ctx context.Context, e *entity.AIEvent) error {
	if m.createFn != nil {
		return m.createFn(ctx, e)
	}
	return nil
}

func (m *mockAIRepo) GetByID(ctx context.Context, id int64) (*entity.AIEvent, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockAIRepo) GetByCode(ctx context.Context, code string) (*entity.AIEvent, error) {
	if m.getByCodeFn != nil {
		return m.getByCodeFn(ctx, code)
	}
	return nil, apierrors.ErrNotFound
}

func (m *mockAIRepo) List(ctx context.Context, licensePlate, gateID, eventType string, offset, limit int) ([]entity.AIEvent, error) {
	if m.listFn != nil {
		return m.listFn(ctx, licensePlate, gateID, eventType, offset, limit)
	}
	return []entity.AIEvent{}, nil
}

func (m *mockAIRepo) Count(ctx context.Context, licensePlate, gateID, eventType string) (int, error) {
	if m.countFn != nil {
		return m.countFn(ctx, licensePlate, gateID, eventType)
	}
	return 0, nil
}

func (m *mockAIRepo) Update(ctx context.Context, id int64, e *entity.AIEvent) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, id, e)
	}
	return nil
}

func (m *mockAIRepo) Delete(ctx context.Context, id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	return nil
}

func validCreateRequest() *dto.CreateAIEventRequest {
	return &dto.CreateAIEventRequest{
		LicensePlate:    "51A-12345",
		ConfidenceScore: 0.92,
		EventType:       "INBOUND",
		GateID:          "GATE-01",
	}
}

func TestAIEventCreateValid(t *testing.T) {
	var saved *entity.AIEvent
	repo := &mockAIRepo{
		createFn: func(ctx context.Context, e *entity.AIEvent) error {
			e.ID = 1
			saved = e
			return nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	e, warning, err := svc.Create(context.Background(), validCreateRequest())
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if e.ID != 1 {
		t.Errorf("Create() ID = %d, want 1", e.ID)
	}
	if saved.EventCode == "" {
		t.Error("Create() did not generate event_code")
	}
	if saved.Timestamp.IsZero() {
		t.Error("Create() did not set timestamp")
	}
	if warning != "" {
		t.Errorf("Create() warning = %q, want empty for high confidence", warning)
	}
}

func TestAIEventCreateWarnsOnLowConfidence(t *testing.T) {
	repo := &mockAIRepo{createFn: func(ctx context.Context, e *entity.AIEvent) error { e.ID = 1; return nil }}
	svc := service.NewServiceWithRepo(repo)

	req := validCreateRequest()
	req.ConfidenceScore = 0.55
	_, warning, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if warning == "" {
		t.Error("Create() warning empty, want low-confidence warning")
	}
}

func TestAIEventCreateRejectsInvalidConfidence(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockAIRepo{})
	req := validCreateRequest()
	req.ConfidenceScore = 1.5
	_, _, err := svc.Create(context.Background(), req)
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestAIEventCreateRequiresLicensePlate(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockAIRepo{})
	req := validCreateRequest()
	req.LicensePlate = ""
	_, _, err := svc.Create(context.Background(), req)
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Create() error = %v, want ErrBadRequest", err)
	}
}

func TestAIEventCreateRejectsDuplicateEventCode(t *testing.T) {
	repo := &mockAIRepo{
		getByCodeFn: func(ctx context.Context, code string) (*entity.AIEvent, error) {
			return &entity.AIEvent{ID: 1, EventCode: code}, nil
		},
	}
	svc := service.NewServiceWithRepo(repo)

	req := validCreateRequest()
	req.EventCode = "AIE-001"
	_, _, err := svc.Create(context.Background(), req)
	if !errors.Is(err, apierrors.ErrConflict) {
		t.Errorf("Create() error = %v, want ErrConflict", err)
	}
}

func TestAIEventGetReturnsErrNotFound(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockAIRepo{})
	_, err := svc.Get(context.Background(), 99)
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestAIEventListReturnsItemsAndCount(t *testing.T) {
	items := []entity.AIEvent{{ID: 1, LicensePlate: "51A-12345"}}
	repo := &mockAIRepo{
		listFn: func(ctx context.Context, lp, gate, et string, offset, limit int) ([]entity.AIEvent, error) {
			return items, nil
		},
		countFn: func(ctx context.Context, lp, gate, et string) (int, error) { return 1, nil },
	}
	svc := service.NewServiceWithRepo(repo)

	got, total, err := svc.List(context.Background(), "51A-12345", "", "", 0, 20)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if total != 1 || len(got) != 1 {
		t.Errorf("List() = (%d items, %d total), want (1, 1)", len(got), total)
	}
}

func TestAIEventUpdateAppliesPartialFields(t *testing.T) {
	existing := &entity.AIEvent{ID: 1, LicensePlate: "51A-12345", ConfidenceScore: 0.9, EventType: "INBOUND", GateID: "GATE-01"}
	repo := &mockAIRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.AIEvent, error) { return existing, nil }}
	svc := service.NewServiceWithRepo(repo)

	newGate := "GATE-02"
	newScore := 0.8
	e, err := svc.Update(context.Background(), 1, &dto.UpdateAIEventRequest{GateID: &newGate, ConfidenceScore: &newScore})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if e.GateID != newGate {
		t.Errorf("Update() GateID = %q, want %q", e.GateID, newGate)
	}
	if e.LicensePlate != "51A-12345" {
		t.Errorf("Update() LicensePlate = %q, want unchanged %q", e.LicensePlate, "51A-12345")
	}
}

func TestAIEventUpdateRejectsInvalidConfidence(t *testing.T) {
	repo := &mockAIRepo{getByIDFn: func(ctx context.Context, id int64) (*entity.AIEvent, error) {
		return &entity.AIEvent{ID: 1}, nil
	}}
	svc := service.NewServiceWithRepo(repo)
	bad := -0.1
	_, err := svc.Update(context.Background(), 1, &dto.UpdateAIEventRequest{ConfidenceScore: &bad})
	if !errors.Is(err, apierrors.ErrBadRequest) {
		t.Errorf("Update() error = %v, want ErrBadRequest", err)
	}
}

func TestAIEventUpdateReturnsErrNotFound(t *testing.T) {
	svc := service.NewServiceWithRepo(&mockAIRepo{})
	_, err := svc.Update(context.Background(), 99, &dto.UpdateAIEventRequest{})
	if !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestAIEventDeleteReturnsErrNotFound(t *testing.T) {
	repo := &mockAIRepo{deleteFn: func(ctx context.Context, id int64) error { return apierrors.ErrNotFound }}
	svc := service.NewServiceWithRepo(repo)
	if err := svc.Delete(context.Background(), 99); !errors.Is(err, apierrors.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}

func TestAIEventTimestampDefaultsToNow(t *testing.T) {
	var saved *entity.AIEvent
	repo := &mockAIRepo{createFn: func(ctx context.Context, e *entity.AIEvent) error { saved = e; return nil }}
	svc := service.NewServiceWithRepo(repo)

	_, _, err := svc.Create(context.Background(), validCreateRequest())
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if saved.Timestamp.Before(time.Now().Add(-time.Minute)) {
		t.Errorf("Create() timestamp = %v, want near now", saved.Timestamp)
	}
}
