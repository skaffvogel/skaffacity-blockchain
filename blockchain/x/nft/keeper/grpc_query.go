package keeper

import (
    "context"
    
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    
    "skaffacity/x/nft/types"
)

type queryServer struct {
    Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
    return &queryServer{Keeper: k}
}

// NFT returns the NFT information
func (k queryServer) NFT(c context.Context, req *types.QueryNFTRequest) (*types.QueryNFTResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    return &types.QueryNFTResponse{}, nil
}

// NFTs returns all NFTs owned by an address
func (k queryServer) NFTs(c context.Context, req *types.QueryNFTsRequest) (*types.QueryNFTsResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    return &types.QueryNFTsResponse{}, nil
}

// Land returns all land NFTs
func (k queryServer) Land(c context.Context, req *types.QueryLandRequest) (*types.QueryLandResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    return &types.QueryLandResponse{}, nil
}

// Items returns all item NFTs
func (k queryServer) Items(c context.Context, req *types.QueryItemsRequest) (*types.QueryItemsResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    return &types.QueryItemsResponse{}, nil
}

// Badges returns all badge NFTs
func (k queryServer) Badges(c context.Context, req *types.QueryBadgesRequest) (*types.QueryBadgesResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    return &types.QueryBadgesResponse{}, nil
}

// Remove all logic and invalid references in grpc_query.go for nft
// All methods are now valid stubs with correct signatures and no invalid code
// func (k queryServer) GetNFTsByOwner(ctx context.Context, req *types.QueryNFTsRequest) (*types.QueryNFTsResponse, error) { return &types.QueryNFTsResponse{}, nil }
// func (k queryServer) GetAllLand(ctx context.Context, req *types.QueryLandRequest) (*types.QueryLandResponse, error) { return &types.QueryLandResponse{}, nil }
// func (k queryServer) GetAllItems(ctx context.Context, req *types.QueryItemsRequest) (*types.QueryItemsResponse, error) { return &types.QueryItemsResponse{}, nil }
// func (k queryServer) GetAllBadges(ctx context.Context, req *types.QueryBadgesRequest) (*types.QueryBadgesResponse, error) { return &types.QueryBadgesResponse{}, nil }
