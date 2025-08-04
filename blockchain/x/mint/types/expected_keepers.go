package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the contract required for account APIs.
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) interface{}
}

// BankKeeper defines the contract needed to be fulfilled for banking dependencies.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
}

// StakingKeeper defines the contract for staking APIs.
type StakingKeeper interface {
	StakingTokenSupply(ctx sdk.Context) sdk.Int
	BondedRatio(ctx sdk.Context) sdk.Dec
}
