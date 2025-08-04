package types

import (
	"fmt"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(minter Minter, params Params) *GenesisState {
	return &GenesisState{
		Minter: minter,
		Params: params,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Minter: DefaultInitialMinter(),
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the mint genesis parameters
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	return ValidateMinter(data.Minter)
}

// ValidateMinter validates the minter object
func ValidateMinter(minter Minter) error {
	if minter.Inflation.IsNegative() {
		return fmt.Errorf("mint parameter inflation should be positive, is %s",
			minter.Inflation.String())
	}
	return nil
}

// GenesisState defines the mint module's genesis state.
type GenesisState struct {
	// minter is a space for holding current inflation information.
	Minter Minter `protobuf:"bytes,1,opt,name=minter,proto3" json:"minter"`
	// params defines all the paramaters of the module.
	Params Params `protobuf:"bytes,2,opt,name=params,proto3" json:"params"`
}
