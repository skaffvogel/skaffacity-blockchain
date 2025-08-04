package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(ctx sdk.Context, key []byte, ptr interface{})
	Set(ctx sdk.Context, key []byte, param interface{})
	GetParamSet(ctx sdk.Context, ps ParamSet)
	SetParamSet(ctx sdk.Context, ps ParamSet)
}

// ParamSet defines an interface for structs containing parameters for a module
type ParamSet interface {
	ParamSetPairs() ParamSetPairs
	Validate() error
}

// ParamSetPairs - Implements params ParamSet
type ParamSetPairs []ParamSetPair

// ParamSetPair - a key/value struct pair
type ParamSetPair struct {
	Key         []byte
	Value       interface{}
	ValidatorFn func(value interface{}) error
}

// NewParamSetPair creates a new ParamSetPair instance
func NewParamSetPair(key []byte, value interface{}, vfn func(value interface{}) error) ParamSetPair {
	return ParamSetPair{key, value, vfn}
}
