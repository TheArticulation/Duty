package duty

import (
	"context"
	"encoding/json"

	"cosmossdk.io/depinject"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/TheArticulation/Duty/x/duty/client"
	"github.com/TheArticulation/Duty/x/duty/genesis"
	"github.com/TheArticulation/Duty/x/duty/keeper"
	"github.com/TheArticulation/Duty/x/duty/modulev1"
	"github.com/TheArticulation/Duty/x/duty/types"
)

// ModuleInputs defines the inputs for the duty module
type ModuleInputs struct {
	depinject.In

	Codec         codec.Codec
	StoreService  store.KVStoreService
	ParamSpace    paramtypes.Subspace
	StakingKeeper *stakingkeeper.Keeper
	Logger        log.Logger
	Config        *modulev1.Module
	ParamsService types.ParamsService
}

// ModuleOutputs defines the outputs for the duty module
type ModuleOutputs struct {
	depinject.Out

	Keeper    keeper.Keeper
	AppModule module.AppModule
}

// HooksInputs defines the inputs for registering hooks
type HooksInputs struct {
	depinject.In

	Keeper        keeper.Keeper
	StakingKeeper *stakingkeeper.Keeper
}

// ProvideModule provides the duty module with dependency injection
func ProvideModule(in ModuleInputs) (ModuleOutputs, error) {
	// Set default parameters from config if provided
	if in.Config != nil {
		// Set the default parameters in the param space
		in.ParamSpace.Set(context.Background(), types.KeyQuorumNumerator, in.Config.QuorumNum)
		in.ParamSpace.Set(context.Background(), types.KeyQuorumDenominator, in.Config.QuorumDen)
	}

	k := keeper.NewKeeper(
		in.Codec,
		in.StoreService,
		in.ParamSpace,
		in.StakingKeeper,
		in.Logger,
		in.ParamsService,
	)

	appModule := NewAppModule(k)

	return ModuleOutputs{
		Keeper:    k,
		AppModule: appModule,
	}, nil
}

// RegisterHooks registers the duty module hooks with the staking module
func RegisterHooks(in HooksInputs) error {
	hooks := keeper.NewDutyHooks(in.Keeper)
	in.StakingKeeper.SetHooks(hooks)
	return nil
}

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string { return types.ModuleName }
func (AppModuleBasic) RegisterInterfaces(reg codec.Types) {
	types.RegisterInterfaces(reg)
}
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(genesis.DefaultGenesis())
}
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ module.ClientTxEncodingConfig, bz json.RawMessage) error {
	return nil
}

func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return client.GetTxCmd()
}

func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return client.GetQueryCmd()
}

type AppModule struct {
	AppModuleBasic
	Keeper keeper.Keeper
}

func NewAppModule(k keeper.Keeper) AppModule { return AppModule{Keeper: k} }

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.Keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQueryServer(am.Keeper))
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var gs genesis.GenesisState
	cdc.MustUnmarshalJSON(data, &gs)
	genesis.InitGenesis(ctx, am.Keeper, &gs)
	return nil
}
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(genesis.ExportGenesis(ctx, am.Keeper))
}
