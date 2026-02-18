package gorpc

// MiddlewareFunc is a function that wraps procedure execution, allowing cross-cutting
// concerns like logging, authentication, and validation to be applied uniformly.
// The next parameter represents the next middleware in the chain or the final handler.
// Middlewares can inspect/modify input, short-circuit by returning early, or handle
// errors before they reach the handler.
type MiddlewareFunc func(ctx *Context, input interface{}, next func(ctx *Context, input interface{}) (interface{}, error)) (interface{}, error)

// MiddlewareChain manages a collection of middleware functions and provides
// execution in the order they were added. It implements the chain of responsibility
// pattern, where each middleware can either delegate to the next or short-circuit
// the request processing.
type MiddlewareChain struct {
	middlewares []MiddlewareFunc
}

func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (mc *MiddlewareChain) Add(middleware MiddlewareFunc) {
	if middleware != nil {
		mc.middlewares = append(mc.middlewares, middleware)
	}
}

// Execute runs the middleware chain in forward order (first added = first executed),
// with the handler called after all middleware. It builds the chain dynamically by
// wrapping each middleware around the next one, starting from the handler.
// This approach avoids recursion and provides predictable execution order while
// allowing middleware to modify both input and output.
func (mc *MiddlewareChain) Execute(ctx *Context, input interface{}, handler handlerFuncAny) (interface{}, error) {
	if len(mc.middlewares) == 0 {
		return handler(ctx, input)
	}

	// Reverse iterate to build the chain from the end (handler) to the beginning
	var next func(ctx *Context, input interface{}) (interface{}, error)
	next = func(ctx *Context, input interface{}) (interface{}, error) {
		return handler(ctx, input)
	}

	for i := len(mc.middlewares) - 1; i >= 0; i-- {
		middleware := mc.middlewares[i]
		currentNext := next
		next = func(ctx *Context, input interface{}) (interface{}, error) {
			return middleware(ctx, input, currentNext)
		}
	}

	return next(ctx, input)
}
