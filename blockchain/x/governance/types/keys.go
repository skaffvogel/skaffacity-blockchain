package types

const (
	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// Keys for governance store
var (
	ProposalKey       = []byte{0x01}
	NextProposalIDKey = []byte{0x02}
	VoteKey           = []byte{0x03}
)
