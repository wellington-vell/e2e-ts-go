package gorpc

import "net/http"

type Context struct {
	Req    *http.Request
	Res    http.ResponseWriter
	Params map[string]string
}
