package keeper

import (
	context "context"
	"yourapp/x/duty/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type msgServer struct{ k Keeper }

func NewMsgServerImpl(k Keeper) types.MsgServer { return &msgServer{k: k} }

func (s *msgServer) SetDutyMetadata(goCtx context.Context, msg *types.MsgSetDutyMetadata) (*types.Empty, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	valAddr, err := sdk.ValAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	// Only allow the validator operator to set metadata for their own consensus key
	v, found := s.k.StakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "no validator")
	}
	consAddr, _ := v.GetConsAddr()

	s.k.SetDutyMetadata(ctx, consAddr, msg.Metadata)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent("duty_metadata_set",
			sdk.NewAttribute("cons_addr", consAddr.String()),
			sdk.NewAttribute("storage_uri", msg.Metadata.CheckpointStorageURI),
		),
	)
	return &types.Empty{}, nil
}
