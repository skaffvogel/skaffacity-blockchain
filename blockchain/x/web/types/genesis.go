package types

import (
	"fmt"
	"github.com/cosmos/gogoproto/proto"
)

// GenesisState defines the web module's genesis state
type GenesisState struct {
	// this line is used by starport scaffolding # genesis/types/default
	WebConfig WebConfig `json:"web_config"`
}

// DefaultGenesisState returns the default genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		WebConfig: DefaultWebConfig(),
	}
}

// Validate performs basic genesis state validation returning an error upon any failure
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate
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
