package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

const (
	ModuleName = "staking"
	StoreKey   = ModuleName
)

var (
	ErrInvalidAmount = sdkerrors.Register(ModuleName, 101, "invalid amount")
	ErrInsufficientStake = sdkerrors.Register(ModuleName, 102, "insufficient stake")
	ErrNoDelegation = sdkerrors.Register(ModuleName, 103, "no delegation found")
)
