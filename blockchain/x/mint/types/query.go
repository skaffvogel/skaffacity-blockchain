package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// QueryParamsRequest is the request type for the Query/Params RPC method.
type QueryParamsRequest struct{}

// QueryParamsResponse is the response type for the Query/Params RPC method.
type QueryParamsResponse struct {
	// params defines the parameters of the module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

// QueryInflationRequest is the request type for the Query/Inflation RPC method.
type QueryInflationRequest struct{}

// QueryInflationResponse is the response type for the Query/Inflation RPC
// method.
type QueryInflationResponse struct {
	// inflation is the current minting inflation value.
	Inflation sdk.Dec `protobuf:"bytes,1,opt,name=inflation,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"inflation"`
}

// QueryAnnualProvisionsRequest is the request type for the
// Query/AnnualProvisions RPC method.
type QueryAnnualProvisionsRequest struct{}

// QueryAnnualProvisionsResponse is the response type for the
// Query/AnnualProvisions RPC method.
type QueryAnnualProvisionsResponse struct {
	// annual_provisions is the current minting annual provisions value.
	AnnualProvisions sdk.Dec `protobuf:"bytes,1,opt,name=annual_provisions,json=annualProvisions,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"annual_provisions"`
}
