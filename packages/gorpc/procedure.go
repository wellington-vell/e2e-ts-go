package gorpc

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type HandlerFunc func(ctx *Context, input any) (any, error)

type Procedure struct {
	Handler     HandlerFunc
	InputType   reflect.Type
	OutputType  reflect.Type
	Middleware  *MiddlewareChain
	ContextType reflect.Type
	Meta        Meta
	Route       *Route
	Tags        []string
}

func (p *Procedure) HandleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Req:    r,
		Res:    w,
		Params: make(map[string]string),
	}
	p.HandleRequestWithContext(ctx)
}

func (p *Procedure) HandleRequestWithContext(ctx *Context) {
	var input any

	if ctx.Req.Body != nil && ctx.Req.ContentLength > 0 &&
		(ctx.Req.Method == http.MethodPost || ctx.Req.Method == http.MethodPut || ctx.Req.Method == http.MethodPatch) {
		decoder := json.NewDecoder(ctx.Req.Body)
		decoder.Decode(&input)
	}

	if ctx.Params == nil {
		ctx.Params = make(map[string]string)
	}

	var result any
	var err error
	if p.Middleware != nil && len(p.Middleware.middlewares) > 0 {
		result, err = p.Middleware.Execute(ctx, input, p.Handler)
	} else {
		result, err = p.Handler(ctx, input)
	}

	if err != nil {
		writeError(ctx.Res, err)
		return
	}

	ctx.Res.Header().Set("Content-Type", "application/json")

	if ctx.Req.Method == http.MethodPost {
		ctx.Res.WriteHeader(http.StatusCreated)
	} else {
		ctx.Res.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(ctx.Res).Encode(result); err != nil {
		writeError(ctx.Res, err)
		return
	}
}
