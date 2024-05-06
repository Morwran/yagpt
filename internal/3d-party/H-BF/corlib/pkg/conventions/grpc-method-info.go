package conventions

import (
	"bytes"
	"context"
	"regexp"
	"sync"

	"github.com/pkg/errors"
)

// GrpcMethodInfo ...
type GrpcMethodInfo struct {
	ServiceFQN string
	Method     string
	Service    string
	Package    string
}

func (m GrpcMethodInfo) String() string {
	if len(m.ServiceFQN) == 0 || len(m.Method) == 0 {
		return ""
	}
	b := bytes.NewBuffer(nil)
	_ = b.WriteByte('/')
	_, _ = b.WriteString(m.ServiceFQN)
	_ = b.WriteByte('/')
	_, _ = b.WriteString(m.Method)
	return b.String()
}

type methodInfoCtxKey struct{}

// WrapContext ...
func (m *GrpcMethodInfo) WrapContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, methodInfoCtxKey{}, m)
}

// FromContext ...
func (m *GrpcMethodInfo) FromContext(ctx context.Context) bool {
	switch v := ctx.Value(methodInfoCtxKey{}).(type) {
	case *GrpcMethodInfo:
		*m = *v
		return true
	}
	return false
}

// Init ...
func (m *GrpcMethodInfo) Init(source string) error {
	const api = "MethodInfo.Init"

	grpcMethodInfoMx.RLock()
	cached := cachedGrpcMethodInfo[source]
	grpcMethodInfoMx.RUnlock()

	switch t := cached.(type) {
	case *GrpcMethodInfo:
		*m = *t
		return nil
	case error:
		return t
	}
	grpcMethodInfoMx.Lock()
	defer grpcMethodInfoMx.Unlock()
	cached = cachedGrpcMethodInfo[source]
	switch t := cached.(type) {
	case *GrpcMethodInfo:
		*m = *t
		return nil
	case error:
		return t
	}
	r := reMethodName.FindAllStringSubmatchIndex(source, -1)
	if len(r) > 0 {
		parts := r[0]
		if len(parts) == 6 {
			res := GrpcMethodInfo{
				Method:     source[parts[4]:parts[5]],
				ServiceFQN: source[parts[2]:parts[3]],
			}
			r1 := reServiceName.FindAllSubmatchIndex([]byte(res.ServiceFQN), -1)
			if len(r1) > 0 {
				parts1 := r1[0]
				if len(parts1) == 6 {
					if parts1[3]-parts1[2] > 0 {
						res.Service = res.ServiceFQN[parts1[2]:parts1[3]]
					} else if parts1[5]-parts1[4] > 0 {
						res.Service = res.ServiceFQN[parts1[4]:parts1[5]]
					}
					if n := len(res.ServiceFQN) - len(res.Service); n > 0 {
						res.Package = res.ServiceFQN[:n-1]
					}
					cachedGrpcMethodInfo[source] = &res
					*m = res
					return nil
				}
			}
		}
	}
	e := errors.Errorf("%s: invnalid source '%s'", api, source)
	cachedGrpcMethodInfo[source] = e
	return e
}

var (
	cachedGrpcMethodInfo = make(map[string]interface{})
	grpcMethodInfoMx     sync.RWMutex
	//reMethodName         = regexp.MustCompile(`^/(\w[^/]+)/(\w+)$`)
	reMethodName  = regexp.MustCompile(`^[\S]*/(\w[^/]+)/(\w+)$`)
	reServiceName = regexp.MustCompile(`(?:\w\.(\w+)$)|(^\w+$)`)
)
