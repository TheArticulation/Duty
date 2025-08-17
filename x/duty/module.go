package duty

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	abci "github.com/tendermint/tendermint/abci/types"

	"yourapp/x/duty/genesis"
	"yourapp/x/duty/keeper"
	"yourapp/x/duty/types"
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string                       { return types.ModuleName }
func (AppModuleBasic) RegisterInterfaces(reg codec.Types) {}
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(genesis.DefaultGenesis())
}
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ module.ClientTxEncodingConfig, bz json.RawMessage) error {
	return nil
}

type AppModule struct {
	AppModuleBasic
	k keeper.Keeper
}

func NewAppModule(k keeper.Keeper) AppModule { return AppModule{k: k} }

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.k))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.k))
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var gs genesis.GenesisState
	cdc.MustUnmarshalJSON(data, &gs)
	genesis.InitGenesis(ctx, am.k, &gs)
	return nil
}
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(genesis.ExportGenesis(ctx, am.k))
}
