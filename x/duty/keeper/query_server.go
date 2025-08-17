package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/duty/types"
)

// Duty implements the Query/Duty gRPC method
func (k Keeper) Duty(c context.Context, req *types.QueryGetDutyRequest) (*types.QueryGetDutyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	duty, found := k.GetDuty(ctx, req.Id)
	if !found {
		return nil, status.Errorf(codes.NotFound, "duty not found")
	}

	return &types.QueryGetDutyResponse{Duty: &duty}, nil
}

// DutyAll implements the Query/DutyAll gRPC method
func (k Keeper) DutyAll(c context.Context, req *types.QueryAllDutyRequest) (*types.QueryAllDutyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	dutyStore := prefix.NewStore(store, types.KeyPrefix(types.DutyKey))

	var duties []types.Duty
	pageRes, err := query.Paginate(dutyStore, req.Pagination, func(key []byte, value []byte) error {
		var duty types.Duty
		if err := k.cdc.Unmarshal(value, &duty); err != nil {
			return err
		}

		duties = append(duties, duty)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDutyResponse{Duty: duties, Pagination: pageRes}, nil
}

// Params implements the Query/Params gRPC method
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

