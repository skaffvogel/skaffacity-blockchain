package types

const (
    StoreKey     = "marketplace"
    RouterKey    = "marketplace"
    QuerierRoute = "marketplace"
)

// Keys for marketplace store
var (
	ListingKey       = []byte{0x01}
	NextListingIDKey = []byte{0x02}
	NFTListingKey    = []byte{0x03}
)
