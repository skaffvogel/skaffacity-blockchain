package types

import (
	"time"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Listing represents a marketplace listing
type Listing struct {
	ID        uint64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Creator   string    `protobuf:"bytes,2,opt,name=creator,proto3" json:"creator,omitempty"`
	NFTID     string    `protobuf:"bytes,3,opt,name=nft_id,json=nftId,proto3" json:"nft_id,omitempty"`
	Price     sdk.Coin  `protobuf:"bytes,4,opt,name=price,proto3" json:"price"`
	Active    bool      `protobuf:"varint,5,opt,name=active,proto3" json:"active,omitempty"`
	CreatedAt time.Time `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3,stdtime" json:"created_at"`
	SoldTo    string    `protobuf:"bytes,7,opt,name=sold_to,json=soldTo,proto3" json:"sold_to,omitempty"`
	SoldAt    time.Time `protobuf:"bytes,8,opt,name=sold_at,json=soldAt,proto3,stdtime" json:"sold_at"`
}

// MarketStats represents marketplace statistics
type MarketStats struct {
	TotalListings  uint64   `protobuf:"varint,1,opt,name=total_listings,json=totalListings,proto3" json:"total_listings,omitempty"`
	ActiveListings uint64   `protobuf:"varint,2,opt,name=active_listings,json=activeListings,proto3" json:"active_listings,omitempty"`
	SoldItems      uint64   `protobuf:"varint,3,opt,name=sold_items,json=soldItems,proto3" json:"sold_items,omitempty"`
	TotalVolume    sdk.Coin `protobuf:"bytes,4,opt,name=total_volume,json=totalVolume,proto3" json:"total_volume"`
}

// ProtoMessage implements the proto.Message interface for Listing.
func (l *Listing) ProtoMessage() {}

// Reset implements the proto.Message interface for Listing.
func (l *Listing) Reset() {}

// String implements the fmt.Stringer interface for Listing.
func (l *Listing) String() string { return "Listing stub" }