package conventions

import (
	"context"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
)

const (
	//SysHeaderPrefix common of system GRPC metadata and HTTP headers
	SysHeaderPrefix = "x-sys-"
)

const (
	// UserAgentHeader web requests
	UserAgentHeader = "user-agent"

	// HostNameHeader remote host name going from grpc-client
	HostNameHeader = "host-name"
)

const (
	//LoggerLevelHeader notes to change log level in current context of operation
	LoggerLevelHeader = SysHeaderPrefix + "log-lvl"

	//AppNameHeader holds application name for incoming outgoing requests
	AppNameHeader = SysHeaderPrefix + "app-name"

	//AppVersionHeader holds application version for incoming outgoing requests
	AppVersionHeader = SysHeaderPrefix + "app-ver"
)

// ClientName user agent extractor
var ClientName clientNameExtractor

// Incoming extracts user agent from incoming context
func (a clientNameExtractor) Incoming(ctx context.Context, defVal string) string {
	if ret, ok := a.extractClientName(ctx, a.mdIncoming); ok {
		return ret
	}
	return defVal
}

// Outgoing extracts user agent from outgoing context
func (a clientNameExtractor) Outgoing(ctx context.Context, defVal string) string {
	if ret, ok := a.extractClientName(ctx, a.mdOutgoing); ok {
		return ret
	}
	return defVal
}

type clientNameExtractor struct{}

func (clientNameExtractor) mdIncoming(ctx context.Context) metadata.MD {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return md
	}
	return nil
}

func (clientNameExtractor) mdOutgoing(ctx context.Context) metadata.MD {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		return md
	}
	return nil
}

func (clientNameExtractor) extractClientName(ctx context.Context, mdExtractor func(ctx context.Context) metadata.MD) (string, bool) {
	if md := mdExtractor(ctx); md != nil {
		ff := [...]func(metadata.MD, string) string{
			GetAppName, GetUserAgent,
		}
		for _, f := range ff {
			if v := f(md, ""); v != "" {
				return v, true
			}
		}
	}
	return "", false
}

// GetAppName -
func GetAppName(md metadata.MD, defVal string) string {
	return getAnyFromMD(md, defVal, AppNameHeader)
}

// GetUserAgent -
func GetUserAgent(md metadata.MD, defVal string) string {
	return getAnyFromMD(md, defVal, runtime.MetadataPrefix+UserAgentHeader, UserAgentHeader)
}

// GetHostName -
func GetHostName(md metadata.MD, defVal string) string {
	return getAnyFromMD(md, defVal, HostNameHeader)
}

func getAnyFromMD(md metadata.MD, defVal string, keys ...string) string {
	var s string
	var k string
	for _, k = range keys {
		if data := md.Get(k); len(data) > 0 {
			s = data[0]
			break
		}
	}
	if strings.HasSuffix(k, UserAgentHeader) {
		s = removeGrpcGoPrefixAndVersion(s)
	}
	if s == "" {
		return defVal
	}
	return s
}

func removeGrpcGoPrefixAndVersion(s string) string {
	const suffix = "grpc-go/" //мерзотный суффикс
	if n := strings.Index(s, suffix); n >= 0 {
		return strings.TrimRight(s[:n], " ")
	}
	return s
}
