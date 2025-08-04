package keeper

import (
    "context"
    "github.com/cosmos/cosmos-sdk/codec"
    storetypes "github.com/cosmos/cosmos-sdk/store/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    "skaffacity/x/nft/types"
    typekeeper "skaffacity/x/nft/types/keeper"
)

type Keeper struct {
    storeKey    storetypes.StoreKey
    cdc         codec.BinaryCodec
    bankKeeper  typekeeper.BankKeeper
}

func NewKeeper(
    cdc codec.BinaryCodec,
    storeKey storetypes.StoreKey,
    bankKeeper typekeeper.BankKeeper,
) *Keeper {
    return &Keeper{
        storeKey:   storeKey,
        cdc:        cdc,
        bankKeeper: bankKeeper,
    }
}

// MintNFT creates a new NFT
func (k Keeper) MintNFT(ctx sdk.Context, nft types.NFT) error {
    store := ctx.KVStore(k.storeKey)
    key := append([]byte(nft.Type + "/"), []byte(nft.ID)...)
    
    if store.Has(key) {
        return sdkerrors.Wrap(types.ErrNFTExists, "NFT already exists")
    }
    
    bz := k.cdc.MustMarshal(&nft)
    store.Set(key, bz)
    return nil
}

// TransferNFT transfers ownership of an NFT
func (k Keeper) TransferNFT(ctx sdk.Context, nftID string, from, to string) error {
    nft, err := k.GetNFT(ctx, nftID)
    if err != nil {
        return err
    }
    
    if nft.Owner != from {
        return sdkerrors.Wrap(types.ErrUnauthorized, "sender is not the owner")
    }
    
    if !nft.Transferable {
        return sdkerrors.Wrap(types.ErrNonTransferable, "NFT cannot be transferred")
    }
    
    nft.Owner = to
    return k.UpdateNFT(ctx, nft)
}

// GetNFT retrieves an NFT by ID
func (k Keeper) GetNFT(ctx sdk.Context, id string) (types.NFT, error) {
    store := ctx.KVStore(k.storeKey)
    bz := store.Get([]byte(id))
    if bz == nil {
        return types.NFT{}, sdkerrors.Wrap(types.ErrNFTNotFound, id)
    }
    
    var nft types.NFT
    k.cdc.MustUnmarshal(bz, &nft)
    return nft, nil
}

// UpdateNFT updates an existing NFT
func (k Keeper) UpdateNFT(ctx sdk.Context, nft types.NFT) error {
    store := ctx.KVStore(k.storeKey)
    key := append([]byte(nft.Type + "/"), []byte(nft.ID)...)
    
    if !store.Has(key) {
        return sdkerrors.Wrap(types.ErrNFTNotFound, nft.ID)
    }
    
    bz := k.cdc.MustMarshal(&nft)
    store.Set(key, bz)
    return nil
}

// Stub missing keeper methods for nft
func (k Keeper) GetNFTsByOwner(ctx context.Context, req *types.QueryNFTsRequest) (*types.QueryNFTsResponse, error) { return nil, nil }
func (k Keeper) GetAllLand(ctx context.Context, req *types.QueryLandRequest) (*types.QueryLandResponse, error) { return nil, nil }
func (k Keeper) GetAllItems(ctx context.Context, req *types.QueryItemsRequest) (*types.QueryItemsResponse, error) { return nil, nil }
func (k Keeper) GetAllBadges(ctx context.Context, req *types.QueryBadgesRequest) (*types.QueryBadgesResponse, error) { return nil, nil }
