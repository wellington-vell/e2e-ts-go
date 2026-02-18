package openapi

import (
	"reflect"
	"testing"
)

func TestConvertToOpenAPIPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/api/todos/:id", "/api/todos/{id}"},
		{"/todos/:id", "/todos/{id}"},
		{"/api/users/:userId/posts/:postId", "/api/users/{userId}/posts/{postId}"},
		{"/api/todos", "/api/todos"},
		{"/", "/"},
		{"", ""},
		{"/api/:nested/:params/here", "/api/{nested}/{params}/here"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := convertToOpenAPIPath(tt.input)
			if result != tt.expected {
				t.Errorf("convertToOpenAPIPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateSpec(t *testing.T) {
	procedures := []ProcedureInfo{
		{
			Path:       "/api/todos/:id",
			Method:     "GET",
			PathParams: []string{"id"},
		},
		{
			Path:       "/api/todos",
			Method:     "GET",
			PathParams: nil,
		},
	}

	spec, err := GenerateSpec(procedures)
	if err != nil {
		t.Fatalf("GenerateSpec failed: %v", err)
	}

	paths := spec["paths"].(map[string]interface{})

	// Check that the path with :id is converted to {id}
	if _, exists := paths["/api/todos/{id}"]; !exists {
		t.Errorf("Expected path /api/todos/{id} to exist, got: %v", reflect.ValueOf(paths).MapKeys())
	}

	// Check that the original path with :id does NOT exist
	if _, exists := paths["/api/todos/:id"]; exists {
		t.Errorf("Path /api/todos/:id should not exist in OpenAPI spec")
	}

	// Check that parameters are included
	getOp := paths["/api/todos/{id}"].(map[string]interface{})["get"].(map[string]interface{})
	params := getOp["parameters"]
	if params == nil {
		t.Error("Expected parameters to be present for /api/todos/{id}")
	} else {
		paramsList := params.([]interface{})
		if len(paramsList) != 1 {
			t.Errorf("Expected 1 parameter, got %d", len(paramsList))
		}
		param := paramsList[0].(map[string]interface{})
		if param["name"] != "id" {
			t.Errorf("Expected parameter name 'id', got %v", param["name"])
		}
		if param["in"] != "path" {
			t.Errorf("Expected parameter in 'path', got %v", param["in"])
		}
	}
}
