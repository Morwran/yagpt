package internal

import (
	"context"

	"google.golang.org/grpc"
)

//InvalidConn ...
type InvalidConn struct {
	Err error
}

//Invoke call unary RPC
func (c *InvalidConn) Invoke(_ context.Context, _ string, _ interface{}, _ interface{}, _ ...grpc.CallOption) error {
	return c.Err
}

//NewStream begins a streaming RPC.
func (c *InvalidConn) NewStream(_ context.Context, _ *grpc.StreamDesc, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, c.Err
}
