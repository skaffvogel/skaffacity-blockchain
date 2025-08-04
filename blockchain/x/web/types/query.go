package types

import (
	"context"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc"
)

// Query client interface
type QueryClient interface {
	WebConfig(ctx context.Context, req *QueryGetWebConfigRequest, opts ...grpc.CallOption) (*QueryGetWebConfigResponse, error)
	WebConfigAll(ctx context.Context, req *QueryAllWebConfigRequest, opts ...grpc.CallOption) (*QueryAllWebConfigResponse, error)
}

// NewQueryClient creates a new query client
func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func (c *queryClient) WebConfig(ctx context.Context, req *QueryGetWebConfigRequest, opts ...grpc.CallOption) (*QueryGetWebConfigResponse, error) {
	// Simple implementation for now
	return &QueryGetWebConfigResponse{
		WebConfig: DefaultWebConfig(),
	}, nil
}

func (c *queryClient) WebConfigAll(ctx context.Context, req *QueryAllWebConfigRequest, opts ...grpc.CallOption) (*QueryAllWebConfigResponse, error) {
	// Simple implementation for now
	return &QueryAllWebConfigResponse{
		WebConfig:  []WebConfig{DefaultWebConfig()},
		Pagination: &query.PageResponse{},
	}, nil
}

// Query request/response types
type QueryGetWebConfigRequest struct{}

func (q *QueryGetWebConfigRequest) ProtoMessage() {}
func (q *QueryGetWebConfigRequest) Reset()        { *q = QueryGetWebConfigRequest{} }
func (q *QueryGetWebConfigRequest) String() string { return "QueryGetWebConfigRequest{}" }

type QueryGetWebConfigResponse struct {
	WebConfig WebConfig `json:"web_config"`
}

func (q *QueryGetWebConfigResponse) ProtoMessage() {}
func (q *QueryGetWebConfigResponse) Reset()        { *q = QueryGetWebConfigResponse{} }
func (q *QueryGetWebConfigResponse) String() string { return "QueryGetWebConfigResponse{}" }

type QueryAllWebConfigRequest struct{
	Pagination *query.PageRequest `json:"pagination,omitempty"`
}

func (q *QueryAllWebConfigRequest) ProtoMessage() {}
func (q *QueryAllWebConfigRequest) Reset()         { *q = QueryAllWebConfigRequest{} }
func (q *QueryAllWebConfigRequest) String() string { return "QueryAllWebConfigRequest{}" }

type QueryAllWebConfigResponse struct {
	WebConfig  []WebConfig         `json:"web_config"`
	Pagination *query.PageResponse `json:"pagination,omitempty"`
}

func (q *QueryAllWebConfigResponse) ProtoMessage() {}
func (q *QueryAllWebConfigResponse) Reset()        { *q = QueryAllWebConfigResponse{} }
func (q *QueryAllWebConfigResponse) String() string { return "QueryAllWebConfigResponse{}" }
