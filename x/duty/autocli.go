package duty

import (
	"github.com/TheArticulation/Duty/x/duty/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/autocli"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func (am AppModuleBasic) GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Transaction commands for the duty module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		autocli.GetTxMsg[*types.MsgSetDutyMetadata](),
		autocli.GetTxMsg[*types.MsgRotateCheckpointKey](),
		autocli.GetTxMsg[*types.MsgBindCheckpointKey](),
	)

	return cmd
}

// GetQueryCmd returns the query commands for this module
func (am AppModuleBasic) GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Query commands for the duty module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		autocli.GetQuery[*types.QueryDutySetRequest](),
		autocli.GetQuery[*types.QueryDutyMetadataRequest](),
	)

	return cmd
}

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModuleBasic) AutoCLIOptions() *autocli.AutoCLIOptions {
	return &autocli.AutoCLIOptions{
		Module:        types.ModuleName,
		AddressPrefix: "cosmos",
	}
}
