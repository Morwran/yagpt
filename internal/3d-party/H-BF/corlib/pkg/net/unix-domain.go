package net

import (
	"context"
	"errors"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

//UDS unix domain socket utils
var UDS unixDomainSocketUtils

//findUnixSocketFromURI обнаруживаем unix-socket из URI
func findUnixSocketFromURI(uri string) string {
	var ret string
	sm := reDetermineUnixSocket.FindAllStringSubmatch(uri, -1)
	for i := range sm {
		if found := sm[i]; len(found) > 1 {
			ret = path.Join(ret, found[1])
			if stat, _ := os.Stat(ret); stat != nil && stat.Mode()&fs.ModeSocket == fs.ModeSocket {
				return ret
			}
		}
	}
	return ""
}

var (
	reDetermineUnixSocket = regexp.MustCompile(`(?:[^:]+:[\\/]{2})?(/?\s*[^/\\]+)`)
)

type (
	unixDomainSocketUtils struct{}

	unixDomainRoundTripper struct {
		http.RoundTripper
		overrider func(req *http.Request) (*http.Response, error)
	}
)

//AddressFromURI address from URI
func (unixDomainSocketUtils) AddressFromURI(uri string) net.Addr {
	if a := findUnixSocketFromURI(uri); len(a) > 0 {
		return &net.UnixAddr{Net: "unix", Name: a}
	}
	return nil
}

//EnrichClient adds RoundTripper for http+unix:// scheme
func (unixDomainSocketUtils) EnrichClient(c *http.Client) *http.Client {
	ret := new(http.Client)
	if c == nil {
		*ret = *http.DefaultClient
	} else {
		*ret = *c
	}
	ret.Transport = enrichWithUnixDomainRoundTripper(ret.Transport)
	return ret
}

func enrichWithUnixDomainRoundTripper(transport http.RoundTripper) http.RoundTripper {
	if transport == nil {
		transport = http.DefaultTransport
	}
	ret := &unixDomainRoundTripper{
		RoundTripper: transport,
	}
	var dialer net.Dialer
	dial := func(ctx context.Context, _, addr string) (net.Conn, error) {
		h, _, _ := net.SplitHostPort(addr)
		return dialer.DialContext(ctx, "unix", h)
	}
	t := http.Transport{
		DialContext: dial,
	}
	ret.overrider = func(req *http.Request) (*http.Response, error) {
		return t.RoundTrip(req)
	}
	return ret
}

//RoundTrip impl http.RoundTripper for http+unix:// scheme
func (rt *unixDomainRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if SchemeUnixHTTP.Is(req.URL.Scheme) {
		path1 := path.Join(req.URL.Host, req.URL.Path)
		sock := findUnixSocketFromURI(path1)
		if len(sock) == 0 {
			return nil, errors.New("address socket is unresolved")
		}
		oldHost, oldURL := req.Host, req.URL
		defer func() {
			req.Host, req.URL = oldHost, oldURL
		}()
		url2 := *oldURL
		url2.Scheme = string(SchemeHTTP)
		url2.Host = sock
		url2.Path = path1[strings.Index(path1, sock)+len(sock):]
		req.Host, req.URL = "", &url2
		return rt.overrider(req)
	}
	return rt.RoundTripper.RoundTrip(req)
}
