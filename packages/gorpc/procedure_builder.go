package gorpc

import (
	"reflect"
)

type ProcedureBuilder struct {
	inputType   reflect.Type
	outputType  reflect.Type
	handler     HandlerFunc
	middleware  *MiddlewareChain
	contextType reflect.Type
	meta        Meta
	route       *Route
	tags        []string
}

func OS() *ProcedureBuilder {
	return &ProcedureBuilder{
		middleware: NewMiddlewareChain(),
	}
}

func (pb *ProcedureBuilder) Input(inputType interface{}) *ProcedureBuilder {
	if inputType != nil {
		pb.inputType = reflect.TypeOf(inputType)
		if pb.inputType.Kind() == reflect.Ptr {
			pb.inputType = pb.inputType.Elem()
		}
	}
	return pb
}

func (pb *ProcedureBuilder) Handler(handler HandlerFunc) *ProcedureBuilder {
	pb.handler = handler
	return pb
}

func (pb *ProcedureBuilder) Use(middleware MiddlewareFunc) *ProcedureBuilder {
	if middleware != nil {
		pb.middleware.Add(middleware)
	}
	return pb
}

func (pb *ProcedureBuilder) Context(contextType interface{}) *ProcedureBuilder {
	if contextType != nil {
		pb.contextType = reflect.TypeOf(contextType)
	}
	return pb
}

func (pb *ProcedureBuilder) Output(outputType interface{}) *ProcedureBuilder {
	if outputType != nil {
		pb.outputType = reflect.TypeOf(outputType)
		if pb.outputType.Kind() == reflect.Ptr {
			pb.outputType = pb.outputType.Elem()
		}
	}
	return pb
}

func (pb *ProcedureBuilder) Meta(meta Meta) *ProcedureBuilder {
	pb.meta = meta
	return pb
}

func (pb *ProcedureBuilder) Route(route Route) *ProcedureBuilder {
	pb.route = &route
	return pb
}

func (pb *ProcedureBuilder) Tag(tags ...string) *ProcedureBuilder {
	pb.tags = append(pb.tags, tags...)
	return pb
}

// Build panics if handler or route is missing
func (pb *ProcedureBuilder) Build() *Procedure {
	if pb.handler == nil {
		panic("procedure handler is required")
	}

	if pb.route == nil {
		panic("procedure route is required")
	}

	return &Procedure{
		Handler:     pb.handler,
		InputType:   pb.inputType,
		OutputType:  pb.outputType,
		Middleware:  pb.middleware,
		ContextType: pb.contextType,
		Meta:        pb.meta,
		Route:       pb.route,
		Tags:        pb.tags,
	}
}
