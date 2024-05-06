package net

import (
	"context"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

//Listen listen net address
func Listen(endpoint *Endpoint) (net.Listener, error) {
	addr, err := endpoint.Address()
	if err != nil {
		return nil, err
	}
	switch endpoint.endpointAddress.(type) {
	case endpointAddressTCP:
		return net.Listen(endpoint.Network(), addr)
	case endpointAddressUnix:
		return ListenUnixDomain(addr)
	}
	return nil, errors.Errorf("Listen: unsupported network '%s'", endpoint.Network())
}

//ListenUnixDomain safe listen unix domain socket
func ListenUnixDomain(addr string) (net.Listener, error) {
	const (
		unix    = "unix"
		connect = "connect"
		check   = "check-socket"
		un      = "unlink"
	)

	retErr := &net.OpError{
		Net: unix,
		Addr: &net.UnixAddr{
			Net:  unix,
			Name: addr,
		},
	}
	if st, e := os.Stat(addr); e == nil {
		mode := st.Mode()
		if mode&os.ModeSocket != os.ModeSocket {
			retErr.Op = check
			retErr.Err = syscall.ENOTSOCK
			return nil, retErr
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		var c net.Conn
		c, e = (&net.Dialer{}).DialContext(ctx, unix, addr)
		cancel()
		if e == nil {
			_ = c.Close()
			retErr.Op = connect
			retErr.Err = syscall.EADDRINUSE
			return nil, retErr
		}
		if netErr := new(net.OpError); errors.As(e, &netErr) {
			if netErr.Timeout() {
				retErr.Op = connect
				retErr.Err = syscall.EADDRINUSE
				return nil, retErr
			}
		}
		if e = syscall.Unlink(addr); e != nil {
			retErr.Op = un
			retErr.Err = e
			return nil, retErr
		}
	}
	return net.Listen(unix, addr)
}
