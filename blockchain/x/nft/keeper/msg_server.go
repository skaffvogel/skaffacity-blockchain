package keeper

import (
    "context"
    "skaffacity/x/nft/types"
)

type msgServer struct {
    Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
    return &msgServer{Keeper: keeper}
}

func (k msgServer) MintNFT(goCtx context.Context, msg *types.MsgMintNFT) (*types.MsgMintNFTResponse, error) {
    return &types.MsgMintNFTResponse{}, nil
}

func (k msgServer) TransferNFT(goCtx context.Context, msg *types.MsgTransferNFT) (*types.MsgTransferNFTResponse, error) {
    return &types.MsgTransferNFTResponse{}, nil
}

func (k msgServer) AttachToItem(goCtx context.Context, msg *types.MsgAttachToItem) (*types.MsgAttachToItemResponse, error) {
    return &types.MsgAttachToItemResponse{}, nil
}
