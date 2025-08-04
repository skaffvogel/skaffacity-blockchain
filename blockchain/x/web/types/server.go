package types

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// QueryServer interface
type QueryServer interface {
	WebConfig(ctx context.Context, req *QueryGetWebConfigRequest) (*QueryGetWebConfigResponse, error)
	WebConfigAll(ctx context.Context, req *QueryAllWebConfigRequest) (*QueryAllWebConfigResponse, error)
}

// RegisterQueryHandlerClient registers the http handlers for service Query to "mux"
func RegisterQueryHandlerClient(ctx context.Context, mux *runtime.ServeMux, client QueryClient) error {
	// Simple implementation for now
	return nil
}

// RegisterMsgServer registers the msg server
func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	// Simple implementation for now
}

// RegisterQueryServer registers the query server  
func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	// Simple implementation for now
}
