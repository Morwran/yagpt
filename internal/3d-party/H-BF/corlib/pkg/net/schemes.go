package net

import (
	"strings"
)

//Schema ...
type Schema string

//Is check is s == another
func (s Schema) Is(another interface{}) bool {
	var s1 string
	switch v := another.(type) {
	case string:
		s1 = v
	case []byte:
		s1 = string(v)
	case Schema:
		s1 = string(v)
	default:
		return false
	}
	return strings.EqualFold(string(s), s1)
}

const (
	//SchemeUnixHTTP обозначим схему как http-unix://[/]socket_path/path[?query]
	SchemeUnixHTTP Schema = "http+unix"
	//SchemeHTTPS like https://
	SchemeHTTPS Schema = "https"
	//SchemeHTTP like http://
	SchemeHTTP Schema = "http"
)
