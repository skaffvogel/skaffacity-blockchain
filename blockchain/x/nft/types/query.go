package types

type QueryServer interface{}

type QueryBadgesRequest struct{ Id string }
type QueryBadgesResponse struct{ Badges []Badge }
type QueryNFTRequest struct{ Id string }
type QueryNFTResponse struct{ Nft NFT }
type QueryNFTsRequest struct{ Owner string }
type QueryNFTsResponse struct{ Nfts []NFT }
type QueryLandRequest struct{}
type QueryLandResponse struct{ Land []Land }
type QueryItemsRequest struct{}
type QueryItemsResponse struct{ Items []Item }

type Badge struct{}
func (b *Badge) Marshal() ([]byte, error) { return nil, nil }
func (b *Badge) Unmarshal([]byte) error { return nil }
func (b *Badge) Size() int { return 0 }
func (b *Badge) MarshalTo([]byte) (int, error) { return 0, nil }
func (b *Badge) UnmarshalTo(interface{}) error { return nil }
func (b *Badge) MarshalToSizedBuffer([]byte) (int, error) { return 0, nil }

type Land struct{}
func (l *Land) Marshal() ([]byte, error) { return nil, nil }
func (l *Land) Unmarshal([]byte) error { return nil }
func (l *Land) Size() int { return 0 }
func (l *Land) MarshalTo([]byte) (int, error) { return 0, nil }
func (l *Land) UnmarshalTo(interface{}) error { return nil }
func (l *Land) MarshalToSizedBuffer([]byte) (int, error) { return 0, nil }

type Item struct{}
func (i *Item) Marshal() ([]byte, error) { return nil, nil }
func (i *Item) Unmarshal([]byte) error { return nil }
func (i *Item) Size() int { return 0 }
func (i *Item) MarshalTo([]byte) (int, error) { return 0, nil }
func (i *Item) UnmarshalTo(interface{}) error { return nil }
func (i *Item) MarshalToSizedBuffer([]byte) (int, error) { return 0, nil }
