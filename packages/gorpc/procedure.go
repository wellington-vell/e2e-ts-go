package gorpc

import (
	"encoding/json"
	"net/http"
	"reflect"
)

// HandlerFunc is a generic type-safe handler function
type HandlerFunc[TInput, TOutput any] func(ctx *Context, input TInput) (TOutput, error)

// handlerFuncAny is a type-erased handler function for internal storage
type handlerFuncAny func(ctx *Context, input any) (any, error)

// ProcedureAny is an interface for type-erased procedure storage
type ProcedureAny interface {
	HandleRequest(w http.ResponseWriter, r *http.Request)
	HandleRequestWithContext(ctx *Context)
}

// Procedure is a generic type-safe procedure
type Procedure[TInput, TOutput any] struct {
	Handler     HandlerFunc[TInput, TOutput]
	handlerAny  handlerFuncAny // internal type-erased handler
	InputType   reflect.Type
	OutputType  reflect.Type
	Middleware  *MiddlewareChain
	ContextType reflect.Type
	Meta        Meta
	Route       *Route
	Tags        []string
	ErrorCodes  []int
	PathParams  []string // explicit path parameters like ["id"] for /todos/:id
}

func (p *Procedure[TInput, TOutput]) HandleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Req:    r,
		Res:    w,
		Params: make(map[string]string),
	}
	p.HandleRequestWithContext(ctx)
}

func (p *Procedure[TInput, TOutput]) HandleRequestWithContext(ctx *Context) {
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
		result, err = p.Middleware.Execute(ctx, input, p.handlerAny)
	} else {
		result, err = p.handlerAny(ctx, input)
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
