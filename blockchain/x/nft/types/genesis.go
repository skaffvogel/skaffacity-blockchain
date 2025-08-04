package types

// GenesisState defines the NFT module's genesis state
type GenesisState struct {
    NFTs []NFT `json:"nfts"`
}

// DefaultGenesisState returns the NFT module's default genesis state
func DefaultGenesisState() *GenesisState {
    return &GenesisState{
        NFTs: []NFT{},
    }
}

// Validate performs basic genesis state validation
func (gs GenesisState) Validate() error {
    return nil // Simplified for now
}
