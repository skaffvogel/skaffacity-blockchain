package keeper

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines the expected bank keeper
type BankKeeper interface {
    SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
    SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
