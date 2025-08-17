package params

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/TheArticulation/Duty/x/duty/types"

	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
)

// Service implements the modern params service for the duty module
type Service struct {
	storeService store.KVStoreService
	subspace     string
}

// NewService creates a new params service for the duty module
func NewService(storeService store.KVStoreService) *Service {
	return &Service{
		storeService: storeService,
		subspace:     types.ModuleName,
	}
}

// Get retrieves a parameter value from the store
func (s *Service) Get(ctx context.Context, subspace, key string) ([]byte, error) {
	if subspace != s.subspace {
		return nil, fmt.Errorf("invalid subspace: %s", subspace)
	}

	store := s.storeService.OpenKVStore(ctx)
	value := store.Get([]byte(key))
	if value == nil {
		return nil, fmt.Errorf("parameter not found: %s", key)
	}

	return value, nil
}

// Set stores a parameter value in the store
func (s *Service) Set(ctx context.Context, subspace, key string, value []byte) error {
	if subspace != s.subspace {
		return fmt.Errorf("invalid subspace: %s", subspace)
	}

	store := s.storeService.OpenKVStore(ctx)
	store.Set([]byte(key), value)
	return nil
}

// Has checks if a parameter exists in the store
func (s *Service) Has(ctx context.Context, subspace, key string) (bool, error) {
	if subspace != s.subspace {
		return false, fmt.Errorf("invalid subspace: %s", subspace)
	}

	store := s.storeService.OpenKVStore(ctx)
	return store.Has([]byte(key)), nil
}

// GetParams retrieves all duty module parameters
func (s *Service) GetParams(ctx context.Context) (types.Params, error) {
	params := types.DefaultParams()

	// Try to get quorum numerator
	if numBytes, err := s.Get(ctx, s.subspace, "QuorumNumerator"); err == nil {
		var num uint32
		if err := json.Unmarshal(numBytes, &num); err == nil {
			params.QuorumNumerator = num
		}
	}

	// Try to get quorum denominator
	if denBytes, err := s.Get(ctx, s.subspace, "QuorumDenominator"); err == nil {
		var den uint32
		if err := json.Unmarshal(denBytes, &den); err == nil {
			params.QuorumDenominator = den
		}
	}

	return params, nil
}

// SetParams stores all duty module parameters
func (s *Service) SetParams(ctx context.Context, params types.Params) error {
	// Validate parameters
	if err := params.Validate(); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	// Store quorum numerator
	numBytes, err := json.Marshal(params.QuorumNumerator)
	if err != nil {
		return fmt.Errorf("failed to marshal quorum numerator: %w", err)
	}
	if err := s.Set(ctx, s.subspace, "QuorumNumerator", numBytes); err != nil {
		return fmt.Errorf("failed to set quorum numerator: %w", err)
	}

	// Store quorum denominator
	denBytes, err := json.Marshal(params.QuorumDenominator)
	if err != nil {
		return fmt.Errorf("failed to marshal quorum denominator: %w", err)
	}
	if err := s.Set(ctx, s.subspace, "QuorumDenominator", denBytes); err != nil {
		return fmt.Errorf("failed to set quorum denominator: %w", err)
	}

	return nil
}

// ServiceInputs defines the inputs for the params service
type ServiceInputs struct {
	depinject.In
	StoreService store.KVStoreService
}

// ServiceOutputs defines the outputs for the params service
type ServiceOutputs struct {
	depinject.Out
	ParamsService types.ParamsService
}

// ProvideParamsService provides the duty module params service
func ProvideParamsService(in ServiceInputs) ServiceOutputs {
	service := NewService(in.StoreService)
	return ServiceOutputs{
		ParamsService: service,
	}
}
