package gorpc

import (
	"encoding/json"
	"reflect"
)

// ProcedureBuilder is a fluent builder for creating type-safe procedures.
// It uses generics to ensure compile-time type safety while providing a clean
// API for configuring procedure behavior. The builder captures the input and
// output types from the generic parameters, eliminating the need for manual
// type registration.
type ProcedureBuilder[TInput, TOutput any] struct {
	inputType   reflect.Type
	outputType  reflect.Type
	handler     HandlerFunc[TInput, TOutput]
	middleware  *MiddlewareChain
	contextType reflect.Type
	meta        Meta
	route       *Route
	tags        []string
	errorCodes  []int
}

// OS (Operation Specification) creates a new ProcedureBuilder for defining
// a typed procedure. The generic parameters TInput and TOutput define the
// request and response types respectively. This is the starting point for
// creating a procedure using the fluent builder pattern:
//
//	procedure := gorpc.OS[CreateTodoInput, Todo]().
//	    Handler(func(ctx *gorpc.Context, input CreateTodoInput) (Todo, error) { ... }).
//	    Route(gorpc.Route{Method: "POST", Path: "/todos"}).
//	    Build()
func OS[TInput, TOutput any]() *ProcedureBuilder[TInput, TOutput] {
	return &ProcedureBuilder[TInput, TOutput]{
		middleware: NewMiddlewareChain(),
	}
}

func (pb *ProcedureBuilder[TInput, TOutput]) Handler(handler HandlerFunc[TInput, TOutput]) *ProcedureBuilder[TInput, TOutput] {
	pb.handler = handler
	return pb
}

func (pb *ProcedureBuilder[TInput, TOutput]) Use(middleware MiddlewareFunc) *ProcedureBuilder[TInput, TOutput] {
	if middleware != nil {
		pb.middleware.Add(middleware)
	}
	return pb
}

func (pb *ProcedureBuilder[TInput, TOutput]) Context(contextType interface{}) *ProcedureBuilder[TInput, TOutput] {
	if contextType != nil {
		pb.contextType = reflect.TypeOf(contextType)
	}
	return pb
}

func (pb *ProcedureBuilder[TInput, TOutput]) Meta(meta Meta) *ProcedureBuilder[TInput, TOutput] {
	pb.meta = meta
	return pb
}

func (pb *ProcedureBuilder[TInput, TOutput]) Route(route Route) *ProcedureBuilder[TInput, TOutput] {
	pb.route = &route
	return pb
}

func (pb *ProcedureBuilder[TInput, TOutput]) Tag(tags ...string) *ProcedureBuilder[TInput, TOutput] {
	pb.tags = append(pb.tags, tags...)
	return pb
}

func (pb *ProcedureBuilder[TInput, TOutput]) Errors(errorCodes ...int) *ProcedureBuilder[TInput, TOutput] {
	pb.errorCodes = errorCodes
	return pb
}

// Build constructs a Procedure from the builder configuration. It extracts
// type information from the generic parameters by creating zero values and
// using reflection. The handler is wrapped to provide JSON serialization and
// deserialization, allowing handlers to work with strongly-typed inputs and
// outputs while the framework handles the HTTP request/response mapping.
//
// Panics if handler or route is not configured, as these are required for
// a valid procedure. This fail-fast approach catches configuration errors
// at startup rather than on first request.
func (pb *ProcedureBuilder[TInput, TOutput]) Build() *Procedure[TInput, TOutput] {
	if pb.handler == nil {
		panic("procedure handler is required")
	}

	if pb.route == nil {
		panic("procedure route is required")
	}

	var zeroInput TInput
	inputType := reflect.TypeOf(zeroInput)
	if inputType != nil {
		pb.inputType = inputType
		if pb.inputType.Kind() == reflect.Ptr {
			pb.inputType = pb.inputType.Elem()
		}
	}

	var zeroOutput TOutput
	outputType := reflect.TypeOf(zeroOutput)
	if outputType != nil {
		pb.outputType = outputType
		if pb.outputType.Kind() == reflect.Ptr {
			pb.outputType = pb.outputType.Elem()
		}
	}

	handlerAny := func(ctx *Context, input any) (any, error) {
		var typedInput TInput

		if input == nil {
		} else {
			inputBytes, err := json.Marshal(input)
			if err != nil {
				return nil, NewHTTPError(400, "Invalid input format: "+err.Error())
			}

			if err := json.Unmarshal(inputBytes, &typedInput); err != nil {
				return nil, NewHTTPError(400, "Invalid input structure: "+err.Error())
			}
		}

		result, err := pb.handler(ctx, typedInput)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	return &Procedure[TInput, TOutput]{
		Handler:     pb.handler,
		handlerAny:  handlerAny,
		InputType:   pb.inputType,
		OutputType:  pb.outputType,
		Middleware:  pb.middleware,
		ContextType: pb.contextType,
		Meta:        pb.meta,
		Route:       pb.route,
		Tags:        pb.tags,
		ErrorCodes:  pb.errorCodes,
	}
}
