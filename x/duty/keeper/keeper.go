package keeper

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"yourapp/x/duty/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	paramSpace paramtypes.Subspace

	StakingKeeper interface {
		GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator
		GetValidator(ctx sdk.Context, addr sdk.ValAddress) (stakingtypes.Validator, bool)
	}
}

func NewKeeper(cdc codec.BinaryCodec, key sdk.StoreKey, ps paramtypes.Subspace, sk interface{}) Keeper {
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{cdc: cdc, storeKey: key, paramSpace: ps, StakingKeeper: sk}
}

// Set & Get params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}
func (k Keeper) SetParams(ctx sdk.Context, p types.Params) { k.paramSpace.SetParamSet(ctx, &p) }

// Duty metadata CRUD
func (k Keeper) SetDutyMetadata(ctx sdk.Context, valConsAddr sdk.ConsAddress, meta types.DutyMetadata) {
	store := ctx.KVStore(k.storeKey)
	bz, _ := json.Marshal(meta)
	store.Set(types.DutyMetaKey(valConsAddr.Bytes()), bz)
}
func (k Keeper) GetDutyMetadata(ctx sdk.Context, valConsAddr sdk.ConsAddress) (types.DutyMetadata, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DutyMetaKey(valConsAddr.Bytes()))
	if bz == nil {
		return types.DutyMetadata{}, false
	}
	var dm types.DutyMetadata
	_ = json.Unmarshal(bz, &dm)
	return dm, true
}

// DutySet view: expose current consensus validators with optional metadata
type DutyValidator struct {
	ValConsAddr string              `json:"val_cons_addr"`
	VotingPower string              `json:"voting_power"` // string to avoid precision issues in JSON
	Metadata    *types.DutyMetadata `json:"metadata,omitempty"`
}

func (k Keeper) GetDutySet(ctx sdk.Context) ([]DutyValidator, types.Params) {
	vals := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	out := make([]DutyValidator, 0, len(vals))
	for _, v := range vals {
		consAddr, _ := v.GetConsAddr()
		dv := DutyValidator{
			ValConsAddr: consAddr.String(),
			VotingPower: v.GetTokens().String(),
		}
		if meta, ok := k.GetDutyMetadata(ctx, consAddr); ok {
			dv.Metadata = &meta
		}
		out = append(out, dv)
	}
	return out, k.GetParams(ctx)
}
