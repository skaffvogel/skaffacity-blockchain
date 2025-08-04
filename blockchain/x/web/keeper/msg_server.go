package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"skaffacity/x/web/types"
)

func (k msgServer) UpdateWebConfig(goCtx context.Context, msg *types.MsgUpdateWebConfig) (*types.MsgUpdateWebConfigResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	k.SetWebConfig(ctx, msg.Config)

	return &types.MsgUpdateWebConfigResponse{}, nil
}

// SetDeveloperAddress handles setting the developer address for fee distribution
func (k msgServer) SetDeveloperAddress(goCtx context.Context, msg *types.MsgSetDeveloperAddress) (*types.MsgSetDeveloperAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Add permission check (only admin/governance should be able to set this)
	
	err := k.Keeper.SetDeveloperAddress(ctx, msg.DeveloperAddress)
	if err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"developer_address_set",
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("developer_address", msg.DeveloperAddress),
		),
	)

	return &types.MsgSetDeveloperAddressResponse{}, nil
}

// EnableFeeDistribution handles enabling/disabling fee distribution
func (k msgServer) EnableFeeDistribution(goCtx context.Context, msg *types.MsgEnableFeeDistribution) (*types.MsgEnableFeeDistributionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Add permission check (only admin/governance should be able to set this)
	
	err := k.Keeper.EnableFeeDistribution(ctx, msg.Enabled)
	if err != nil {
		return nil, err
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"fee_distribution_enabled",
			sdk.NewAttribute("creator", msg.Creator),
			sdk.NewAttribute("enabled", fmt.Sprintf("%t", msg.Enabled)),
		),
	)

	return &types.MsgEnableFeeDistributionResponse{}, nil
}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
