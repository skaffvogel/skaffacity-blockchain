package web

import (
	"fmt"
	"github.com/cosmos/gogoproto/proto"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"skaffacity/x/web/keeper"
	"skaffacity/x/web/types"
)

// GenesisState represents the web module genesis state
type GenesisState struct {
	WebConfig types.WebConfig `json:"web_config"`
}

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		WebConfig: types.DefaultWebConfig(),
	}
}

// ValidateGenesis validates the genesis state
func (gs GenesisState) Validate() error {
	// Validate web config
	if gs.WebConfig.Port == 0 {
		return fmt.Errorf("web config port cannot be zero")
	}
	
	if gs.WebConfig.Host == "" {
		return fmt.Errorf("web config host cannot be empty")
	}
	
	return nil
}

// ProtoMessage implements proto.Message interface
func (gs *GenesisState) ProtoMessage() {}

// Reset implements proto.Message interface
func (gs *GenesisState) Reset() {
	*gs = GenesisState{}
}

// String implements proto.Message interface
func (gs *GenesisState) String() string {
	return fmt.Sprintf("GenesisState{WebConfig: %s}", gs.WebConfig.String())
}

// Marshal implements ProtoMarshaler interface
func (gs *GenesisState) Marshal() ([]byte, error) {
	return proto.Marshal(gs)
}

// Unmarshal implements ProtoMarshaler interface
func (gs *GenesisState) Unmarshal(data []byte) error {
	return proto.Unmarshal(data, gs)
}

// MarshalTo implements ProtoMarshaler interface
func (gs *GenesisState) MarshalTo(data []byte) (int, error) {
	marshaled, err := gs.Marshal()
	if err != nil {
		return 0, err
	}
	copy(data, marshaled)
	return len(marshaled), nil
}

// Size implements ProtoMarshaler interface
func (gs *GenesisState) Size() int {
	marshaled, _ := gs.Marshal()
	return len(marshaled)
}

// MarshalToSizedBuffer implements ProtoMarshaler interface
func (gs *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	marshaled, err := gs.Marshal()
	if err != nil {
		return 0, err
	}
	copy(dAtA[len(dAtA)-len(marshaled):], marshaled)
	return len(marshaled), nil
}

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState GenesisState) {
	// Set web config if provided
	k.SetWebConfig(ctx, genState.WebConfig)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *GenesisState {
	genesis := GenesisState{}
	
	// Get web config (using the new method that returns default if not found)
	webConfig := k.GetWebConfig(ctx)
	genesis.WebConfig = webConfig

	return &genesis
}
