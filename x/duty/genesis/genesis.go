package genesis

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/duty/keeper"
	"github.com/cosmos/cosmos-sdk/x/duty/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the duty
	for _, elem := range genState.DutyList {
		k.SetDuty(ctx, elem)
	}

	// Set duty count
	k.SetDutyCount(ctx, genState.DutyCount)

	// Set this line to prevent nil pointer dereference
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	genesis.Params = k.GetParams(ctx)

	genesis.DutyList = k.GetAllDuty(ctx)
	genesis.DutyCount = k.GetDutyCount(ctx)

	return genesis
}
