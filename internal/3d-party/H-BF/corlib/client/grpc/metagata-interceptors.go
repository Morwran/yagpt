package grpc

import (
	"context"

	"github.com/H-BF/corlib/pkg/conventions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AddMetaToOutgoingUnaryInterceptor - specifies fixed headers for every rpc as kv array
func AddMetaToOutgoingUnaryInterceptor(kv ...string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if len(kv) != 0 {
			ctx = metadata.AppendToOutgoingContext(ctx, kv...)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// AddMetaToOutgoingStreamInterceptor - specifies fixed headers for every rpc as kv array
func AddMetaToOutgoingStreamInterceptor(kv ...string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, kv...)
		return streamer(ctx, desc, cc, method, opts...)
	}
}

// HostNamePropagator add host name into outgoing metadata
type HostNamePropagator string

// ClientUnary unary client interceptor
func (h HostNamePropagator) ClientUnary() grpc.UnaryClientInterceptor {
	return AddMetaToOutgoingUnaryInterceptor(conventions.HostNameHeader, string(h))
}

// ClientStream stream client interceptor
func (h HostNamePropagator) ClientStream() grpc.StreamClientInterceptor {
	return AddMetaToOutgoingStreamInterceptor(conventions.HostNameHeader, string(h))
}
