package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
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
		paramtypes.NewParamSetPair(KeyQuorumNumerator, &p.QuorumNumerator, func(i interface{}) error { return nil }),
		paramtypes.NewParamSetPair(KeyQuorumDenominator, &p.QuorumDenominator, func(i interface{}) error { return nil }),
	}
}

func DefaultParams() Params {
	return Params{QuorumNumerator: DefaultQuorumNum, QuorumDenominator: DefaultQuorumDen}
}
