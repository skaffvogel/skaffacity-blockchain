package keeper

import (
    "context"
    
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    
    "skaffacity/x/marketplace/types"
)

type queryServer struct {
    Keeper
}

func NewQueryServer(k Keeper) types.QueryServer {
    return &queryServer{Keeper: k}
}

// Listing returns the listing information
func (k queryServer) Listing(c context.Context, req *types.QueryListingRequest) (*types.QueryListingResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    
    ctx := c
    resp, err := k.GetListing(ctx, req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// Listings returns all active listings
func (k queryServer) Listings(c context.Context, req *types.QueryListingsRequest) (*types.QueryListingsResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    
    ctx := c
    resp, err := k.GetAllListings(ctx, req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// ListingsByType returns all active listings of a specific NFT type
func (k queryServer) ListingsByType(c context.Context, req *types.QueryListingsByTypeRequest) (*types.QueryListingsByTypeResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    
    ctx := c
    resp, err := k.GetListingsByType(ctx, req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// ListingsByOwner returns all listings by an owner
func (k queryServer) ListingsByOwner(c context.Context, req *types.QueryListingsByOwnerRequest) (*types.QueryListingsByOwnerResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    
    ctx := c
    resp, err := k.GetListingsByOwner(ctx, req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}

// MarketStats returns marketplace statistics
func (k queryServer) MarketStats(c context.Context, req *types.QueryMarketStatsRequest) (*types.QueryMarketStatsResponse, error) {
    if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }
    
    ctx := c
    resp, err := k.GetMarketStats(ctx, req)
    if err != nil {
        return nil, err
    }
    return resp, nil
}
