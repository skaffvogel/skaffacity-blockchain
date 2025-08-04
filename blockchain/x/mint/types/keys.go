package types

const (
	// ModuleName defines the module name
	ModuleName = "mint"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

var (
	// MinterKey is the key to use for the keeper store.
	MinterKey = []byte{0x00}
	
	// ParamsKey is the key to use for the params store.
	ParamsKey = []byte{0x01}
)
