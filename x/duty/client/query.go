package client

import (
	"github.com/TheArticulation/Duty/x/duty/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the query commands for this module
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Query commands for the duty module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdDutySet(),
		GetCmdDutyMetadata(),
	)

	return cmd
}

// GetCmdDutySet returns the command to query the duty set
func GetCmdDutySet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "duty-set",
		Short: "Query the current duty set",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.DutySet(cmd.Context(), &types.QueryDutySetRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdDutyMetadata returns the command to query duty metadata
func GetCmdDutyMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "duty-metadata [consensus-address]",
		Short: "Query duty metadata for a validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			consensusAddress := args[0]
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.DutyMetadata(cmd.Context(), &types.QueryDutyMetadataRequest{
				ConsAddr: consensusAddress,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
