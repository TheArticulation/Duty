package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/TheArticulation/Duty/x/duty/types"
)

type Keeper struct {
	cdc           codec.Codec
	storeService  store.KVStoreService
	paramSpace    paramtypes.Subspace
	stakingKeeper *stakingkeeper.Keeper
	logger        log.Logger
	paramsService types.ParamsService
}

func NewKeeper(
	cdc codec.Codec,
	storeService store.KVStoreService,
	ps paramtypes.Subspace,
	stakingKeeper *stakingkeeper.Keeper,
	logger log.Logger,
	paramsService types.ParamsService,
) Keeper {
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		cdc:           cdc,
		storeService:  storeService,
		paramSpace:    ps,
		stakingKeeper: stakingKeeper,
		logger:        logger,
		paramsService: paramsService,
	}
}

// Set & Get params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	// Try to get params from the modern params service first
	if k.paramsService != nil {
		if params, err := k.getParamsFromService(ctx); err == nil {
			return params
		}
	}

	// Fallback to legacy param space
	var p types.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

func (k Keeper) SetParams(ctx sdk.Context, p types.Params) {
	// Try to set params using the modern params service first
	if k.paramsService != nil {
		if err := k.setParamsToService(ctx, p); err == nil {
			return
		}
	}

	// Fallback to legacy param space
	k.paramSpace.SetParamSet(ctx, &p)
}

// getParamsFromService retrieves parameters from the modern params service
func (k Keeper) getParamsFromService(ctx sdk.Context) (types.Params, error) {
	if k.paramsService == nil {
		return types.DefaultParams(), fmt.Errorf("params service not available")
	}

	// Convert sdk.Context to context.Context for the service
	serviceCtx := context.Background()
	return k.paramsService.GetParams(serviceCtx)
}

// setParamsToService sets parameters using the modern params service
func (k Keeper) setParamsToService(ctx sdk.Context, params types.Params) error {
	if k.paramsService == nil {
		return fmt.Errorf("params service not available")
	}

	// Convert sdk.Context to context.Context for the service
	serviceCtx := context.Background()
	return k.paramsService.SetParams(serviceCtx, params)
}

// Duty metadata CRUD
func (k Keeper) SetDutyMetadata(ctx sdk.Context, valConsAddr sdk.ConsAddress, meta types.DutyMetadata) {
	store := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(meta)
	store.Set(types.DutyMetaKey(valConsAddr.Bytes()), bz)
}
func (k Keeper) GetDutyMetadata(ctx sdk.Context, valConsAddr sdk.ConsAddress) (types.DutyMetadata, bool) {
	store := k.storeService.OpenKVStore(ctx)
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
	vals := k.stakingKeeper.GetBondedValidatorsByPower(ctx)
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

// NewDutyHooks creates a new DutyHooks instance
func NewDutyHooks(k Keeper) DutyHooks {
	return DutyHooks{k: k}
}
