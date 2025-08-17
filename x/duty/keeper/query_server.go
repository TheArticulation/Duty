package keeper

import (
	context "context"
	"yourapp/x/duty/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type queryServer struct{ k Keeper }

func NewQueryServer(k Keeper) types.QueryServer { return &queryServer{k: k} }

func (q *queryServer) DutySet(goCtx context.Context, _ *types.QueryDutySetRequest) (*types.QueryDutySetResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	set, params := q.k.GetDutySet(ctx)
	return &types.QueryDutySetResponse{Validators: set, QuorumNum: params.QuorumNumerator, QuorumDen: params.QuorumDenominator}, nil
}
func (q *queryServer) DutyMetadata(goCtx context.Context, req *types.QueryDutyMetadataRequest) (*types.QueryDutyMetadataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	consAddr, err := sdk.ConsAddressFromBech32(req.ConsAddr)
	if err != nil {
		return nil, err
	}
	meta, ok := q.k.GetDutyMetadata(ctx, consAddr)
	if !ok {
		return &types.QueryDutyMetadataResponse{}, nil
	}
	return &types.QueryDutyMetadataResponse{Metadata: &meta}, nil
}
