package gorpc

type MiddlewareFunc func(ctx *Context, input interface{}, next func(ctx *Context, input interface{}) (interface{}, error)) (interface{}, error)

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

// Execute executes middleware in the order it was added, with the handler called last
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
