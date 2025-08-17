package types

import (
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyMaxDutyDuration = []byte("MaxDutyDuration")
	KeyMinDutyDuration = []byte("MinDutyDuration")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(maxDutyDuration, minDutyDuration time.Duration) Params {
	return Params{
		MaxDutyDuration: maxDutyDuration,
		MinDutyDuration: minDutyDuration,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		24*time.Hour, // 24 hours max duty duration
		1*time.Hour,  // 1 hour min duty duration
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMaxDutyDuration, &p.MaxDutyDuration, validateMaxDutyDuration),
		paramtypes.NewParamSetPair(KeyMinDutyDuration, &p.MinDutyDuration, validateMinDutyDuration),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateMaxDutyDuration(p.MaxDutyDuration); err != nil {
		return err
	}
	if err := validateMinDutyDuration(p.MinDutyDuration); err != nil {
		return err
	}
	if p.MaxDutyDuration <= p.MinDutyDuration {
		return ErrInvalidDutyDuration
	}
	return nil
}

func validateMaxDutyDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return ErrInvalidDutyDuration
	}
	if v <= 0 {
		return ErrInvalidDutyDuration
	}
	return nil
}

func validateMinDutyDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return ErrInvalidDutyDuration
	}
	if v <= 0 {
		return ErrInvalidDutyDuration
	}
	return nil
}
