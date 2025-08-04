package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

const ModuleName = "marketplace"

var (
	ErrNFTNotFound        = sdkerrors.Register(ModuleName, 1, "NFT not found")
	ErrUnauthorized       = sdkerrors.Register(ModuleName, 2, "unauthorized")
	ErrNFTAlreadyListed   = sdkerrors.Register(ModuleName, 3, "NFT already listed")
	ErrInvalidPrice       = sdkerrors.Register(ModuleName, 4, "invalid price")
	ErrListingNotFound    = sdkerrors.Register(ModuleName, 5, "listing not found")
	ErrListingNotActive   = sdkerrors.Register(ModuleName, 6, "listing not active")
	ErrSelfPurchase       = sdkerrors.Register(ModuleName, 7, "cannot buy your own NFT")
	ErrNonTransferable    = sdkerrors.Register(ModuleName, 8, "NFT cannot be transferred")
)
