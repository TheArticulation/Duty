package genesis

import (
	"yourapp/x/duty/keeper"
	"yourapp/x/duty/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Params types.Params `json:"params"`
}

func DefaultGenesis() *GenesisState { return &GenesisState{Params: types.DefaultParams()} }

func InitGenesis(ctx sdk.Context, k keeper.Keeper, data *GenesisState) {
	if err := data.Params.Validate(); err == nil {
		k.SetParams(ctx, data.Params)
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *GenesisState {
	return &GenesisState{Params: k.GetParams(ctx)}
}
