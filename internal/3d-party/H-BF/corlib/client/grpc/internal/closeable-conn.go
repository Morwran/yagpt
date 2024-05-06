package internal

import (
	"context"
	"io"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

//ErrConnClosed send when grpc is closed
var ErrConnClosed = errors.New("GRPC conn is closed")

//MakeCloseable ...
func MakeCloseable(c grpc.ClientConnInterface) *closableConn { //nolint:revive
	type connType = struct {
		grpc.ClientConnInterface
	}
	if false {
		var _ grpc.ClientConnInterface = *(*connType)(nil)
	}
	ret := &closableConn{c: new(atomic.Value)}
	ret.c.Store(connType{c})
	var mx sync.Mutex
	var closed bool
	v := ret.c
	ret.close = func() error {
		mx.Lock()
		defer mx.Unlock()
		var e error
		if !closed {
			if closer, _ := c.(io.Closer); closer != nil {
				e = closer.Close()
			}
			if e == nil {
				v.Store(connType{&InvalidConn{Err: ErrConnClosed}})
				closed = true
			}
		}
		return e
	}
	runtime.SetFinalizer(ret, func(o *closableConn) {
		_ = o.close()
	})
	return ret
}

type closableConn struct {
	c     *atomic.Value
	close func() error
}

//Invoke call unary RPC
func (c *closableConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	conn := c.conn()
	return conn.Invoke(ctx, method, args, reply, opts...)
}

//NewStream begins a streaming RPC.
func (c *closableConn) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	conn := c.conn()
	return conn.NewStream(ctx, desc, method, opts...)
}

//CloseConn it closes connection
func (c *closableConn) CloseConn() error {
	return c.close()
}

func (c *closableConn) conn() grpc.ClientConnInterface {
	return c.c.Load().(grpc.ClientConnInterface)
}
