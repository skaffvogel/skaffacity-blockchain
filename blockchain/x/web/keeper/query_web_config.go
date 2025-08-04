package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"skaffacity/x/web/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) WebConfig(goCtx context.Context, req *types.QueryGetWebConfigRequest) (*types.QueryGetWebConfigResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val := k.GetWebConfig(ctx)

	return &types.QueryGetWebConfigResponse{WebConfig: val}, nil
}

func (k Keeper) WebConfigAll(goCtx context.Context, req *types.QueryAllWebConfigRequest) (*types.QueryAllWebConfigResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var webConfigs []types.WebConfig
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	webConfigStore := prefix.NewStore(store, types.KeyPrefix(types.WebConfigKey))

	pageRes, err := query.Paginate(webConfigStore, req.Pagination, func(key []byte, value []byte) error {
		var webConfig types.WebConfig
		if err := k.cdc.Unmarshal(value, &webConfig); err != nil {
			return err
		}

		webConfigs = append(webConfigs, webConfig)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllWebConfigResponse{WebConfig: webConfigs, Pagination: pageRes}, nil
}
