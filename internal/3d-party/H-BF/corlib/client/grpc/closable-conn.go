package grpc

import (
	"io"

	"github.com/Morwran/yagpt/internal/3d-party/H-BF/corlib/client/grpc/internal"
	"google.golang.org/grpc"
)

// ErrConnClosed send when grpc is closed
var ErrConnClosed = internal.ErrConnClosed

// Closable is closable connect
type Closable interface {
	CloseConn() error
}

// ClosableClientConnInterface grpc client conn interface with close ability
type ClosableClientConnInterface interface {
	grpc.ClientConnInterface
	Closable
}

var _ io.Closer = (*grpc.ClientConn)(nil) //assert

// MakeCloseable ...
func MakeCloseable(c grpc.ClientConnInterface) ClosableClientConnInterface {
	if ret, ok := c.(ClosableClientConnInterface); ok {
		return ret
	}
	return internal.MakeCloseable(c)
}
