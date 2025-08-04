package types

type QueryServer interface{}

type QueryMarketStatsRequest struct{}
type QueryMarketStatsResponse struct{}
type QueryListingRequest struct{ Id string }
type QueryListingResponse struct{ Listing Listing }
type QueryListingsRequest struct{}
type QueryListingsResponse struct{ Listings []Listing }
type QueryListingsByTypeRequest struct{ Type string }
type QueryListingsByTypeResponse struct{ Listings []Listing }
type QueryListingsByOwnerRequest struct{ Owner string }
type QueryListingsByOwnerResponse struct{ Listings []Listing }
