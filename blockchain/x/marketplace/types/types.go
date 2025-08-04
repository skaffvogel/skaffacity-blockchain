package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgCreateListing defines a message to create a new marketplace listing
type MsgCreateListing struct {
    Creator  string  `json:"creator"`
    ItemID   string  `json:"item_id"`
    Price    sdk.Coin `json:"price"`
}

// MsgBuyItem defines a message to buy an item from the marketplace
type MsgBuyItem struct {
    Buyer    string  `json:"buyer"`
    ListingID string `json:"listing_id"`
}
func (l *Listing) Size() int { return 0 }
func (l *Listing) MarshalTo([]byte) (int, error) { return 0, nil }
func (l *Listing) UnmarshalTo(interface{}) error { return nil }
func (l *Listing) MarshalToSizedBuffer([]byte) (int, error) { return 0, nil }
