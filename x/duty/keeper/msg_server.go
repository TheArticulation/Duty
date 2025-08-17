package keeper

import (
	context "context"
	"fmt"
	"yourapp/x/duty/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type msgServer struct{ k Keeper }

func NewMsgServerImpl(k Keeper) types.MsgServer { return &msgServer{k: k} }

func (s *msgServer) SetDutyMetadata(goCtx context.Context, msg *types.MsgSetDutyMetadata) (*emptypb.Empty, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	valAddr, err := sdk.ValAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid valoper")
	}
	// Only allow the validator operator to set metadata for their own consensus key
	v, found := s.k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "no validator")
	}
	consAddr, _ := v.GetConsAddr()

	// Convert proto DutyMetadata to internal DutyMetadata
	metadata := types.DutyMetadata{
		CheckpointPubKey:     msg.Metadata.CheckpointPubKey,
		CheckpointStorageURI: msg.Metadata.CheckpointStorageUri,
	}

	s.k.SetDutyMetadata(ctx, consAddr, metadata)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent("duty_metadata_set",
			sdk.NewAttribute("cons_addr", consAddr.String()),
			sdk.NewAttribute("val_addr", valAddr.String()),
			sdk.NewAttribute("checkpoint_pub_key", metadata.CheckpointPubKey),
			sdk.NewAttribute("storage_uri", metadata.CheckpointStorageURI),
			sdk.NewAttribute("block_height", fmt.Sprintf("%d", ctx.BlockHeight())),
		),
	)
	return &emptypb.Empty{}, nil
}

func (s *msgServer) RotateCheckpointKey(goCtx context.Context, msg *types.MsgRotateCheckpointKey) (*emptypb.Empty, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	valAddr, err := sdk.ValAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid valoper")
	}

	// Only allow the validator operator to rotate their own checkpoint key
	v, found := s.k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "no validator")
	}

	consAddr, _ := v.GetConsAddr()

	// Get existing metadata
	existingMeta, found := s.k.GetDutyMetadata(ctx, consAddr)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "no existing duty metadata")
	}

	// TODO: Verify attestation signature
	// This would verify that the new key is owned by the validator
	// For now, we'll just update the key

	// Update metadata with new checkpoint key
	updatedMeta := types.DutyMetadata{
		CheckpointPubKey:     msg.NewCheckpointPubKey,
		CheckpointStorageURI: existingMeta.CheckpointStorageURI,
	}

	s.k.SetDutyMetadata(ctx, consAddr, updatedMeta)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("duty_checkpoint_key_rotated",
			sdk.NewAttribute("cons_addr", consAddr.String()),
			sdk.NewAttribute("val_addr", valAddr.String()),
			sdk.NewAttribute("old_checkpoint_pub_key", existingMeta.CheckpointPubKey),
			sdk.NewAttribute("new_checkpoint_pub_key", msg.NewCheckpointPubKey),
			sdk.NewAttribute("block_height", fmt.Sprintf("%d", ctx.BlockHeight())),
		),
	)

	return &emptypb.Empty{}, nil
}

func (s *msgServer) BindCheckpointKey(goCtx context.Context, msg *types.MsgBindCheckpointKey) (*emptypb.Empty, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	valAddr, err := sdk.ValAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid valoper")
	}

	// Parse consensus address
	consAddr, err := sdk.ConsAddressFromBech32(msg.ConsensusAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid consensus address")
	}

	// Verify the validator exists and the signer is the operator
	v, found := s.k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "no validator")
	}

	// Verify the consensus address matches the validator's consensus address
	validatorConsAddr, _ := v.GetConsAddr()
	if !consAddr.Equals(validatorConsAddr) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "consensus address mismatch")
	}

	// TODO: Verify binding signature
	// This would verify that the checkpoint key is bound to the consensus address
	// For now, we'll just create the binding

	// Create or update metadata with the bound checkpoint key
	metadata := types.DutyMetadata{
		CheckpointPubKey:     msg.CheckpointPubKey,
		CheckpointStorageURI: "", // Will be set separately via SetDutyMetadata
	}

	s.k.SetDutyMetadata(ctx, consAddr, metadata)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("duty_checkpoint_key_bound",
			sdk.NewAttribute("cons_addr", consAddr.String()),
			sdk.NewAttribute("val_addr", valAddr.String()),
			sdk.NewAttribute("checkpoint_pub_key", msg.CheckpointPubKey),
			sdk.NewAttribute("binding_signature", msg.BindingSignature),
			sdk.NewAttribute("block_height", fmt.Sprintf("%d", ctx.BlockHeight())),
		),
	)

	return &emptypb.Empty{}, nil
}
