package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/duty/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the duty MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// CreateDuty handles the creation of a new duty
func (k msgServer) CreateDuty(goCtx context.Context, msg *types.MsgCreateDuty) (*types.MsgCreateDutyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Check if the duty already exists
	if k.HasDuty(ctx, msg.Id) {
		return nil, sdkerrors.Wrapf(types.ErrDutyAlreadyExists, "duty with id %s already exists", msg.Id)
	}

	// Create the duty
	duty := types.Duty{
		Id:          msg.Id,
		Title:       msg.Title,
		Description: msg.Description,
		Creator:     msg.Creator,
		Status:      types.DutyStatus_PENDING,
		CreatedAt:   ctx.BlockTime(),
	}

	// Store the duty
	k.SetDuty(ctx, duty)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDutyCreated,
			sdk.NewAttribute(types.AttributeKeyDutyId, duty.Id),
			sdk.NewAttribute(types.AttributeKeyCreator, duty.Creator),
		),
	)

	return &types.MsgCreateDutyResponse{}, nil
}

// UpdateDuty handles the update of an existing duty
func (k msgServer) UpdateDuty(goCtx context.Context, msg *types.MsgUpdateDuty) (*types.MsgUpdateDutyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Check if the duty exists
	duty, found := k.GetDuty(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrDutyNotFound, "duty with id %s not found", msg.Id)
	}

	// Check if the sender is the creator
	if duty.Creator != msg.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "only the creator can update the duty")
	}

	// Update the duty
	duty.Title = msg.Title
	duty.Description = msg.Description
	duty.UpdatedAt = ctx.BlockTime()

	// Store the updated duty
	k.SetDuty(ctx, duty)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDutyUpdated,
			sdk.NewAttribute(types.AttributeKeyDutyId, duty.Id),
			sdk.NewAttribute(types.AttributeKeyCreator, duty.Creator),
		),
	)

	return &types.MsgUpdateDutyResponse{}, nil
}

// DeleteDuty handles the deletion of a duty
func (k msgServer) DeleteDuty(goCtx context.Context, msg *types.MsgDeleteDuty) (*types.MsgDeleteDutyResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Check if the duty exists
	duty, found := k.GetDuty(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrDutyNotFound, "duty with id %s not found", msg.Id)
	}

	// Check if the sender is the creator
	if duty.Creator != msg.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "only the creator can delete the duty")
	}

	// Delete the duty
	k.RemoveDuty(ctx, msg.Id)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDutyDeleted,
			sdk.NewAttribute(types.AttributeKeyDutyId, duty.Id),
			sdk.NewAttribute(types.AttributeKeyCreator, duty.Creator),
		),
	)

	return &types.MsgDeleteDutyResponse{}, nil
}

