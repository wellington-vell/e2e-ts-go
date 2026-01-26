package gorpc

import (
	"net/http"
	"strings"
)

type radixNode struct {
	prefix   string
	children map[string]*radixNode
	// paramChild is the child node for parameter segments (e.g., :id, :userId)
	paramChild *radixNode
	// paramName is the name of the parameter (e.g., "id", "userId") - only set if this is a parameter node
	paramName string
	handlers  map[string]ProcedureAny
}

// radixRouter is a custom HTTP router using a radix tree for efficient path matching
type radixRouter struct {
	root *radixNode
}

func NewRouter() *radixRouter {
	return &radixRouter{
		root: &radixNode{
			children: make(map[string]*radixNode),
			handlers: make(map[string]ProcedureAny),
		},
	}
}

func (r *radixRouter) Insert(pattern, method string, procedure ProcedureAny) {
	segments := splitPath(pattern)
	r.insertRecursive(r.root, segments, method, procedure, 0)
}

func (r *radixRouter) insertRecursive(node *radixNode, segments []string, method string, procedure ProcedureAny, depth int) {
	if depth >= len(segments) {
		if node.handlers == nil {
			node.handlers = make(map[string]ProcedureAny)
		}
		node.handlers[method] = procedure
		return
	}

	segment := segments[depth]
	isParam := strings.HasPrefix(segment, ":")
	paramName := ""
	if isParam {
		paramName = segment[1:]
	}

	if isParam {
		if node.paramChild == nil {
			node.paramChild = &radixNode{
				prefix:    segment,
				paramName: paramName,
				children:  make(map[string]*radixNode),
				handlers:  make(map[string]ProcedureAny),
			}
		}
		// Ensure paramName matches if paramChild already exists
		if node.paramChild.paramName != paramName {
			node.paramChild.paramName = paramName
		}
		r.insertRecursive(node.paramChild, segments, method, procedure, depth+1)
	} else {
		child, exists := node.children[segment]
		if !exists {
			child = &radixNode{
				prefix:   segment,
				children: make(map[string]*radixNode),
				handlers: make(map[string]ProcedureAny),
			}
			node.children[segment] = child
		}
		r.insertRecursive(child, segments, method, procedure, depth+1)
	}
}

type MatchResult struct {
	Procedure ProcedureAny
	Params    map[string]string
}

func (r *radixRouter) Match(path, method string) (*MatchResult, bool) {
	segments := splitPath(path)
	params := make(map[string]string)
	procedure := r.matchRecursive(r.root, segments, method, params, 0)
	if procedure == nil {
		return nil, false
	}
	return &MatchResult{
		Procedure: procedure,
		Params:    params,
	}, true
}

func (r *radixRouter) PathExists(path string) bool {
	segments := splitPath(path)
	return r.pathExistsRecursive(r.root, segments, 0)
}

func (r *radixRouter) pathExistsRecursive(node *radixNode, segments []string, depth int) bool {
	if depth >= len(segments) {
		return len(node.handlers) > 0
	}

	segment := segments[depth]

	if child, exists := node.children[segment]; exists {
		if r.pathExistsRecursive(child, segments, depth+1) {
			return true
		}
	}

	if node.paramChild != nil {
		if r.pathExistsRecursive(node.paramChild, segments, depth+1) {
			return true
		}
	}

	return false
}

func (r *radixRouter) matchRecursive(node *radixNode, segments []string, method string, params map[string]string, depth int) ProcedureAny {
	if depth >= len(segments) {
		if node.handlers != nil {
			if proc, ok := node.handlers[method]; ok {
				return proc
			}
		}
		return nil
	}

	segment := segments[depth]

	// Try exact match first (static segment)
	if child, exists := node.children[segment]; exists {
		if proc := r.matchRecursive(child, segments, method, params, depth+1); proc != nil {
			return proc
		}
	}

	if node.paramChild != nil {
		if node.paramChild.paramName != "" {
			params[node.paramChild.paramName] = segment
		}
		if proc := r.matchRecursive(node.paramChild, segments, method, params, depth+1); proc != nil {
			return proc
		}
		// Remove parameter if match failed (backtrack)
		if node.paramChild.paramName != "" {
			delete(params, node.paramChild.paramName)
		}
	}

	return nil
}

func (r *radixRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	result, found := r.Match(req.URL.Path, req.Method)
	if !found {
		if r.PathExists(req.URL.Path) {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		http.NotFound(w, req)
		return
	}

	ctx := &Context{
		Req:    req,
		Res:    w,
		Params: result.Params,
	}

	result.Procedure.HandleRequestWithContext(ctx)
}

func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}
