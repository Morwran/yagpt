package grpc

import (
	"context"
	"io"

	"github.com/H-BF/corlib/pkg/conventions"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

//WithErrorWrapper ...
func WithErrorWrapper(c grpc.ClientConnInterface, serviceNamePrefix string) grpc.ClientConnInterface {
	if _, ok := c.(errWrapperInterface); ok {
		return c
	}
	base := &wrappedErrConn{
		serviceNamePrefix: serviceNamePrefix,
		wrapped:           c,
	}
	if closable, _ := c.(Closable); closable != nil {
		type resType = struct {
			errWrapperInterface
			Closable
		}
		return resType{
			errWrapperInterface: base,
			Closable:            closable,
		}
	}
	if closer, _ := c.(io.Closer); closer != nil {
		type resType = struct {
			errWrapperInterface
			io.Closer
		}
		return resType{
			errWrapperInterface: base,
			Closer:              closer,
		}
	}
	return base
}

type errWrapperInterface interface {
	grpc.ClientConnInterface
	isErrWrapper()
}

type wrappedErrConn struct {
	serviceNamePrefix string
	wrapped           grpc.ClientConnInterface
}

var _ errWrapperInterface = (*wrappedErrConn)(nil)

// Invoke performs a unary RPC and returns after the response is received into reply.
func (c *wrappedErrConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	e := c.wrapped.Invoke(ctx, method, args, reply, opts...)
	return c.wrapError(e, method)
}

// NewStream begins a streaming RPC.
func (c *wrappedErrConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	ret, e := c.wrapped.NewStream(ctx, desc, method, opts...)
	return ret, c.wrapError(e, method)
}

func (c *wrappedErrConn) isErrWrapper() {}

func (c *wrappedErrConn) wrapError(e error, meth string) error {
	if e != nil {
		var mi conventions.GrpcMethodInfo
		if mi.Init(meth) == nil {
			if len(c.serviceNamePrefix) > 0 {
				return errors.Wrapf(e, "%s/%s/%s", c.serviceNamePrefix, mi.Service, mi.Method)
			}
			return errors.Wrapf(e, "%s/%s", mi.Service, mi.Method)
		}
	}
	return e
}
