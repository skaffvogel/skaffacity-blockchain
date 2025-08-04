package app

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mint"
	
	// "skaffacity/x/governance"  // Commented out until AppModuleBasic implemented
	// "skaffacity/x/marketplace" // Commented out until AppModuleBasic implemented
	// "skaffacity/x/nft"         // Commented out until AppModuleBasic implemented
	// "skaffacity/x/staking"     // Commented out until AppModuleBasic implemented
	"skaffacity/x/web"
)

// ModuleBasics defines the module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	authvesting.AppModuleBasic{},
	bank.AppModuleBasic{},
	mint.AppModuleBasic{},
	// nft.AppModuleBasic{},      // TODO: implement AppModuleBasic
	// marketplace.AppModuleBasic{}, // TODO: implement AppModuleBasic  
	// governance.AppModuleBasic{}, // TODO: implement AppModuleBasic
	// staking.AppModuleBasic{},    // TODO: implement AppModuleBasic
	web.AppModuleBasic{},
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() EncodingConfig {
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txConfig := tx.NewTxConfig(marshaler, tx.DefaultSignModes)
	
	encodingConfig := EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txConfig,
		Amino:             codec.NewLegacyAmino(),
	}
	
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

// EncodingConfig specifies the concrete encoding types to use for a given app.
type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Marshaler         codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

// MakeCodecs constructs the *std.Codec and *codec.LegacyAmino instances used by
// simapp. It is useful for tests and clients who do not want to construct the
// full simapp
func MakeCodecs() EncodingConfig {
	return MakeEncodingConfig()
}