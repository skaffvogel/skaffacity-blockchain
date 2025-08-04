package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

const ModuleName = "nft"

var (
	ErrNFTExists           = sdkerrors.Register(ModuleName, 1, "NFT already exists")
	ErrNFTNotFound         = sdkerrors.Register(ModuleName, 2, "NFT not found")
	ErrUnauthorized        = sdkerrors.Register(ModuleName, 3, "unauthorized")
	ErrNonTransferable     = sdkerrors.Register(ModuleName, 4, "NFT cannot be transferred")
	ErrNonTransferableNFT  = sdkerrors.Register(ModuleName, 5, "NFT cannot be transferred")
	ErrInvalidNFTType      = sdkerrors.Register(ModuleName, 6, "invalid NFT type")
)
