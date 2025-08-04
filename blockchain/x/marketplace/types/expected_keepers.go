package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	nfttypes "skaffacity/x/nft/types"
)

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

// NFTKeeper defines the expected NFT keeper interface
type NFTKeeper interface {
	GetNFT(ctx sdk.Context, nftID string) (nfttypes.NFT, error) // Use actual NFT type and error
	TransferNFT(ctx sdk.Context, nftID, from, to string) error
}

// NFT represents an NFT for marketplace operations - using actual NFT type
type NFT = nfttypes.NFT
