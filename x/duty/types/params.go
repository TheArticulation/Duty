package types

import (
	"context"
	"fmt"

	"cosmossdk.io/core/address"
	"cosmossdk.io/depinject"
	paramtypes "cosmossdk.io/x/params/types"
)

const (
	DefaultQuorumNum = uint32(2)
	DefaultQuorumDen = uint32(3)
)

var (
	KeyQuorumNumerator   = []byte("QuorumNumerator")
	KeyQuorumDenominator = []byte("QuorumDenominator")
)

type Params struct {
	QuorumNumerator   uint32 `json:"quorum_num" yaml:"quorum_num"`
	QuorumDenominator uint32 `json:"quorum_den" yaml:"quorum_den"`
}

func (p Params) Validate() error {
	if p.QuorumNumerator == 0 || p.QuorumDenominator == 0 || p.QuorumNumerator > p.QuorumDenominator {
		return fmt.Errorf("invalid quorum %d/%d", p.QuorumNumerator, p.QuorumDenominator)
	}
	return nil
}

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyQuorumNumerator, &p.QuorumNumerator, validateQuorumNumerator),
		paramtypes.NewParamSetPair(KeyQuorumDenominator, &p.QuorumDenominator, validateQuorumDenominator),
	}
}

func validateQuorumNumerator(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("quorum numerator cannot be zero")
	}
	return nil
}

func validateQuorumDenominator(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("quorum denominator cannot be zero")
	}
	return nil
}

func DefaultParams() Params {
	return Params{QuorumNumerator: DefaultQuorumNum, QuorumDenominator: DefaultQuorumDen}
}

// ParamsService defines the interface for parameter operations
type ParamsService interface {
	Get(ctx context.Context, subspace, key string) ([]byte, error)
	Set(ctx context.Context, subspace, key string, value []byte) error
	Has(ctx context.Context, subspace, key string) (bool, error)
	GetParams(ctx context.Context) (Params, error)
	SetParams(ctx context.Context, params Params) error
}

// ParamsInputs defines the inputs for parameter operations
type ParamsInputs struct {
	depinject.In
	ParamsService ParamsService
	AddressCodec  address.Codec
}

// ParamsOutputs defines the outputs for parameter operations
type ParamsOutputs struct {
	depinject.Out
	Params Params
}

// ProvideParams provides the duty module parameters
func ProvideParams(in ParamsInputs) (ParamsOutputs, error) {
	// For now, return default params
	// In a full implementation, you would read from the params service
	return ParamsOutputs{
		Params: DefaultParams(),
	}, nil
}
