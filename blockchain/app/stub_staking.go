package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
)

// StubStakingKeeper provides a minimal implementation for mint module
type StubStakingKeeper struct{}

func (k StubStakingKeeper) StakingTokenSupply(ctx sdk.Context) sdk.Int {
	return sdk.NewInt(1000000) // Return fixed supply for now
}

func (k StubStakingKeeper) BondedRatio(ctx sdk.Context) sdk.Dec {
	return sdk.NewDecWithPrec(67, 2) // Return fixed 67% bonded ratio
}

// Ensure interface compliance
var _ types.StakingKeeper = StubStakingKeeper{}
