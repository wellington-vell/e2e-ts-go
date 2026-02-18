package gorpc

import (
	"net/http"
	"strings"
)

// radixNode represents a node in the radix tree used for efficient path matching.
// Each node stores a prefix (path segment), static children for exact matches,
// a parameter child for dynamic segments (e.g., :id), and HTTP method handlers.
// The parameter child allows efficient matching of routes with variables while
// keeping static paths fast through direct map lookup.
type radixNode struct {
	prefix     string
	children   map[string]*radixNode
	paramChild *radixNode
	paramName  string
	handlers   map[string]ProcedureAny
}

// radixRouter is a custom HTTP router using a radix tree for efficient path matching.
// It supports both static paths and dynamic path parameters (e.g., /todos/:id).
// The radix tree provides O(k) lookup time where k is the path length, making it
// efficient for large route tables. Static paths are matched via map lookup,
// while parameters are handled by a single child node with name capture.
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

// Insert adds a route to the radix tree. The path is split into segments,
// and each segment is inserted recursively. Segments starting with ":" are
// treated as parameters and stored in a dedicated paramChild node, allowing
// a single node to match any value for that parameter position.
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

// Match finds a procedure and extracts path parameters for a given path and HTTP method.
// It traverses the radix tree, first checking for exact segment matches, then falling
// back to parameter matching. When a parameter segment is matched, the captured value
// is stored in the params map using the parameter name as key. This allows handlers to
// access dynamic path segments via the Context.Params field.
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
		if node.paramChild.paramName != "" {
			delete(params, node.paramChild.paramName)
		}
	}

	return nil
}

// ServeHTTP implements the http.Handler interface. It attempts to match the
// incoming request path and method against registered routes. If no match is
// found but the path exists with a different method, it returns 405 Method Not
// Allowed. Otherwise, it returns 404 Not Found. On successful match, it creates
// a Context with path parameters and delegates to the procedure's handler.
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

// splitPath splits a URL path into segments for tree insertion and matching.
// It trims leading and trailing slashes, then splits on "/". An empty path
// returns an empty slice rather than a slice containing an empty string.
// This normalization ensures consistent handling of routes regardless of
// how they are formatted (with or without leading/trailing slashes).
func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}

// GetAllProcedures traverses the entire radix tree and returns all registered
// procedures with their paths and HTTP methods. This is used for generating
// OpenAPI documentation and other introspection features. The procedure
// handles are returned as ProcedureAny for type-erased access.
func (r *radixRouter) GetAllProcedures() []struct {
	Path      string
	Method    string
	Procedure ProcedureAny
} {
	var procedures []struct {
		Path      string
		Method    string
		Procedure ProcedureAny
	}
	r.collectProceduresRecursive(r.root, "", &procedures)
	return procedures
}

// collectProceduresRecursive walks the radix tree depth-first, building the
// full path as it descends. At each node with handlers, it records all HTTP
// method/procedure pairs. For parameter nodes, the path includes the ":paramName"
// syntax to indicate dynamic segments. This allows documentation generators
// to reconstruct the original route patterns from the tree structure.
func (r *radixRouter) collectProceduresRecursive(node *radixNode, currentPath string, procedures *[]struct {
	Path      string
	Method    string
	Procedure ProcedureAny
}) {
	if node.handlers != nil {
		for method, proc := range node.handlers {
			*procedures = append(*procedures, struct {
				Path      string
				Method    string
				Procedure ProcedureAny
			}{
				Path:      currentPath,
				Method:    method,
				Procedure: proc,
			})
		}
	}

	for segment, child := range node.children {
		childPath := currentPath + "/" + segment
		r.collectProceduresRecursive(child, childPath, procedures)
	}

	if node.paramChild != nil {
		paramPath := currentPath + "/:" + node.paramChild.paramName
		r.collectProceduresRecursive(node.paramChild, paramPath, procedures)
	}
}
