package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StakingKeeper defines the expected staking keeper interface
type StakingKeeper interface {
	GetStakedAmount(ctx sdk.Context, address string) sdk.Int
}
