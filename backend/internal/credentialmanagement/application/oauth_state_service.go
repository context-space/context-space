package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	observability "github.com/context-space/cloud-observability"
	"github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
)

const (
	// DefaultStateExpiration is the default expiration time for OAuth states
	DefaultStateExpiration = 15 * time.Minute
)

var (
	ErrOAuthStateDataNotFound = errors.New("oauth state data not found")
)

// OAuthStateServiceImpl implements the OAuthStateService interface
type OAuthStateServiceImpl struct {
	stateRepo domain.OAuthStateRepository
	obs       *observability.ObservabilityProvider
}

// NewOAuthStateService creates a new OAuthStateService
func NewOAuthStateService(stateRepo domain.OAuthStateRepository, obs *observability.ObservabilityProvider) domain.OAuthStateService {
	return &OAuthStateServiceImpl{
		stateRepo: stateRepo,
		obs:       obs,
	}
}

// StoreStateData stores a state with associated data
func (s *OAuthStateServiceImpl) StoreStateData(ctx context.Context, data *domain.OAuthStateData) error {
	ctx, span := s.obs.Tracer.Start(ctx, "OAuthStateService.StoreStateData")
	defer span.End()

	err := s.stateRepo.StoreStateData(ctx, data, DefaultStateExpiration)
	if err != nil {
		return fmt.Errorf("failed to store state: %w", err)
	}

	return nil
}

// GetByState gets by a state and returns associated data
func (s *OAuthStateServiceImpl) GetByState(ctx context.Context, state string) (*domain.OAuthStateData, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "OAuthStateService.GetByState")
	defer span.End()

	data, err := s.stateRepo.GetStateDataByState(ctx, state)
	if err != nil {
		return nil, fmt.Errorf("failed to get state by state: %w", err)
	}

	if data == nil {
		return nil, ErrOAuthStateDataNotFound
	}

	return data, nil
}

// GetByID gets by an id and returns associated data
func (s *OAuthStateServiceImpl) GetByID(ctx context.Context, id string) (*domain.OAuthStateData, error) {
	ctx, span := s.obs.Tracer.Start(ctx, "OAuthStateService.GetByID")
	defer span.End()

	data, err := s.stateRepo.GetStateDataByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get state by id: %w", err)
	}

	if data == nil {
		return nil, ErrOAuthStateDataNotFound
	}

	return data, nil
}

// UpdateStateData updates the data of a state
func (s *OAuthStateServiceImpl) UpdateStateData(ctx context.Context, data *domain.OAuthStateData) error {
	ctx, span := s.obs.Tracer.Start(ctx, "OAuthStateService.UpdateStateData")
	defer span.End()

	return s.stateRepo.UpdateStateData(ctx, data)
}
