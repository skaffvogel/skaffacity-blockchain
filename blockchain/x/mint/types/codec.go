package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
)

// RegisterLegacyAminoCodec registers the necessary x/mint interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// No tx types in mint module
}

// RegisterInterfaces register the interfaces
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	// No interfaces to register in mint module
}

var (
	// Amino is a codec for amino serialization
	Amino     = codec.NewLegacyAmino()
	// ModuleCdc is the codec used by this module
	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptocodec.RegisterCrypto(Amino)
	Amino.Seal()
}
