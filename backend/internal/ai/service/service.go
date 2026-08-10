package service

import (
	"context"
	"fmt"
	"time"

	"my-web-app.com/smart-logistic-hub/internal/ai/dto"
	"my-web-app.com/smart-logistic-hub/internal/ai/entity"
	apierrors "my-web-app.com/smart-logistic-hub/internal/common/errors"
)

const lowConfidenceThreshold = 0.7

type AIRepository interface {
	Create(ctx context.Context, e *entity.AIEvent) error
	GetByID(ctx context.Context, id int64) (*entity.AIEvent, error)
	GetByCode(ctx context.Context, code string) (*entity.AIEvent, error)
	List(ctx context.Context, licensePlate, gateID, eventType string, offset, limit int) ([]entity.AIEvent, error)
	Count(ctx context.Context, licensePlate, gateID, eventType string) (int, error)
	Update(ctx context.Context, id int64, e *entity.AIEvent) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	repo AIRepository
}

func NewService(repo AIRepository) *Service {
	return &Service{repo: repo}
}

func NewServiceWithRepo(repo AIRepository) *Service {
	return NewService(repo)
}

func (s *Service) Create(ctx context.Context, req *dto.CreateAIEventRequest) (*entity.AIEvent, string, error) {
	if req.LicensePlate == "" {
		return nil, "", fmt.Errorf("%w: license_plate is required", apierrors.ErrBadRequest)
	}
	if req.EventType == "" {
		return nil, "", fmt.Errorf("%w: event_type is required", apierrors.ErrBadRequest)
	}
	if req.GateID == "" {
		return nil, "", fmt.Errorf("%w: gate_id is required", apierrors.ErrBadRequest)
	}
	if req.ConfidenceScore < 0 || req.ConfidenceScore > 1 {
		return nil, "", fmt.Errorf("%w: confidence_score must be between 0 and 1", apierrors.ErrBadRequest)
	}

	eventCode := req.EventCode
	if eventCode == "" {
		eventCode = fmt.Sprintf("AIE-%d", time.Now().UnixNano())
	}
	existing, err := s.repo.GetByCode(ctx, eventCode)
	if err != nil && err != apierrors.ErrNotFound {
		return nil, "", err
	}
	if existing != nil {
		return nil, "", fmt.Errorf("%w: event_code %s already exists", apierrors.ErrConflict, eventCode)
	}

	ts := req.Timestamp
	if ts.IsZero() {
		ts = time.Now().UTC()
	}
	e := &entity.AIEvent{
		EventCode:       eventCode,
		LicensePlate:    req.LicensePlate,
		ConfidenceScore: req.ConfidenceScore,
		EventType:       req.EventType,
		GateID:          req.GateID,
		Timestamp:       ts,
		CreatedAt:       time.Now().UTC(),
	}
	if err := s.repo.Create(ctx, e); err != nil {
		return nil, "", err
	}

	warning := ""
	if e.ConfidenceScore < lowConfidenceThreshold {
		warning = fmt.Sprintf("low confidence: %.2f (below %.2f threshold)", e.ConfidenceScore, lowConfidenceThreshold)
	}
	return e, warning, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*entity.AIEvent, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByCode(ctx context.Context, eventCode string) (*entity.AIEvent, error) {
	return s.repo.GetByCode(ctx, eventCode)
}

func (s *Service) List(ctx context.Context, licensePlate, gateID, eventType string, offset, limit int) ([]entity.AIEvent, int, error) {
	events, err := s.repo.List(ctx, licensePlate, gateID, eventType, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.repo.Count(ctx, licensePlate, gateID, eventType)
	if err != nil {
		return nil, 0, err
	}
	return events, count, nil
}

func (s *Service) Update(ctx context.Context, id int64, req *dto.UpdateAIEventRequest) (*entity.AIEvent, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.LicensePlate != nil {
		existing.LicensePlate = *req.LicensePlate
	}
	if req.ConfidenceScore != nil {
		if *req.ConfidenceScore < 0 || *req.ConfidenceScore > 1 {
			return nil, fmt.Errorf("%w: confidence_score must be between 0 and 1", apierrors.ErrBadRequest)
		}
		existing.ConfidenceScore = *req.ConfidenceScore
	}
	if req.EventType != nil {
		existing.EventType = *req.EventType
	}
	if req.GateID != nil {
		existing.GateID = *req.GateID
	}
	if req.Timestamp != nil {
		existing.Timestamp = *req.Timestamp
	}

	if err := s.repo.Update(ctx, id, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
