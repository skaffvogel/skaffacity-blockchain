package types

// GenesisState defines the governance module's genesis state
type GenesisState struct {
	// Empty for now, basic structure
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}

func (gs GenesisState) Validate() error {
	return nil
}
