package net

import (
	"net"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

//Endpoint endpoint to connect to
type Endpoint struct {
	endpointAddress
}

//ParseEndpoint parse endpoint address
func ParseEndpoint(src string) (*Endpoint, error) {
	const (
		api  = "ParseEndpoint"
		unix = "unix"
		tcp  = "tcp"
		tcp4 = "tcp4"
		tcp6 = "tcp6"
	)
	parts := reSchemaAndAddress.FindStringSubmatch(src)
	if len(parts) < 3 {
		return nil, errors.Errorf("%s: address '%s' seems invalid", api, src)
	}
	ep := new(Endpoint)
	schema := strings.ToLower(parts[1])
	addr := parts[2]
	switch schema {
	case unix:
		ep.endpointAddress = endpointAddressUnix{socketPath: addr}
	case tcp, tcp4, tcp6, "":
		h, p, e := net.SplitHostPort(addr)
		if e != nil {
			return nil, errors.Errorf("%s: the addr '%s' has invalid host:port", api, src)
		}
		h, p = strings.TrimSpace(h), strings.TrimSpace(p)
		if len(p) == 0 {
			return nil, errors.Errorf("%s: in addr('%s') port is empty", api, src)
		}
		for _, s := range []*string{&h, &p} {
			if strings.ContainsAny(*s, "/?\\=% #") {
				return nil, errors.Errorf("%s: the addr('%s') is invalid", api, src)
			}
		}
		ep.endpointAddress = endpointAddressTCP{host: h, port: p}
	default:
		return nil, errors.Errorf("%s: the addr '%s' has unsupported schema '%s'", api, src, schema)
	}
	return ep, nil
}

func (ep *Endpoint) String() string {
	s, _ := ep.Address()
	return s
}

//Address makes net address
func (ep *Endpoint) Address() (string, error) {
	const api = "Endpoint.Address"
	switch t := ep.endpointAddress.(type) {
	case endpointAddressTCP:
		return net.JoinHostPort(t.host, t.port), nil
	case endpointAddressUnix:
		return t.socketPath, nil
	}
	return "", errors.Errorf("%s: endpoint is not initialized", api)
}

//HostPort gives host - port if TCP case is
func (ep *Endpoint) HostPort() (host, port string, err error) {
	const api = "Endpoint.HostPort"
	switch t := ep.endpointAddress.(type) {
	case endpointAddressTCP:
		host, port = t.host, t.port
	default:
		err = errors.Errorf("%s: endpoint is not a TCP", api)
	}
	return
}

//IsUnixDomain returns true when endpoint is unix domain socket
func (ep *Endpoint) IsUnixDomain() bool {
	_, ret := ep.endpointAddress.(endpointAddressUnix)
	return ret
}

//FQN full qualified name
func (ep *Endpoint) FQN() string {
	if a, _ := ep.Address(); len(a) > 0 {
		return ep.Network() + "://" + a
	}
	return ""
}

var (
	reSchemaAndAddress          = regexp.MustCompile(`^\s*(?:(\w+):(?://)?)?([^\\\s#?=]+)\s*$`)
	_                           = ParseEndpoint
	_                  net.Addr = (*Endpoint)(nil)
)

type (
	endpointAddress interface {
		isEndpointAddress()
		Network() string
	}

	endpointAddressTCP struct {
		endpointAddress
		host string
		port string
	}

	endpointAddressUnix struct {
		endpointAddress
		socketPath string
	}
)

func (endpointAddressTCP) Network() string {
	return "tcp"
}

func (endpointAddressUnix) Network() string {
	return "unix"
}
