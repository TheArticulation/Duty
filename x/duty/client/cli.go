package client

import (
	"github.com/TheArticulation/Duty/x/duty/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Transaction commands for the duty module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdSetDutyMetadata(),
		GetCmdRotateCheckpointKey(),
		GetCmdBindCheckpointKey(),
	)

	return cmd
}

// GetCmdSetDutyMetadata returns the command to set duty metadata
func GetCmdSetDutyMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-duty-metadata [signer] [checkpoint-pub-key] [checkpoint-storage-uri]",
		Short: "Set duty metadata for a validator",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := args[0]
			checkpointPubKey := args[1]
			checkpointStorageURI := args[2]

			msg := &types.MsgSetDutyMetadata{
				Signer: signer,
				Metadata: &types.DutyMetadata{
					CheckpointPubKey:     checkpointPubKey,
					CheckpointStorageUri: checkpointStorageURI,
				},
			}

			return clientCtx.PrintProto(msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdRotateCheckpointKey returns the command to rotate checkpoint key
func GetCmdRotateCheckpointKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rotate-checkpoint-key [signer] [new-checkpoint-pub-key] [attestation-signature]",
		Short: "Rotate checkpoint signing key for a validator",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := args[0]
			newCheckpointPubKey := args[1]
			attestationSignature := args[2]

			msg := &types.MsgRotateCheckpointKey{
				Signer:               signer,
				NewCheckpointPubKey:  newCheckpointPubKey,
				AttestationSignature: attestationSignature,
			}

			return clientCtx.PrintProto(msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// GetCmdBindCheckpointKey returns the command to bind checkpoint key
func GetCmdBindCheckpointKey() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind-checkpoint-key [signer] [checkpoint-pub-key] [binding-signature] [consensus-address]",
		Short: "Bind checkpoint key to consensus validator",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			signer := args[0]
			checkpointPubKey := args[1]
			bindingSignature := args[2]
			consensusAddress := args[3]

			msg := &types.MsgBindCheckpointKey{
				Signer:           signer,
				CheckpointPubKey: checkpointPubKey,
				BindingSignature: bindingSignature,
				ConsensusAddress: consensusAddress,
			}

			return clientCtx.PrintProto(msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
