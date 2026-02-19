package openapi

import (
	"reflect"
	"testing"
	"time"
)

func TestTypeToSchemaTimeTime(t *testing.T) {
	result := typeToSchema(reflect.TypeOf(time.Time{}))

	if result["type"] != "string" {
		t.Errorf("Expected type to be 'string', got %v", result["type"])
	}

	if result["format"] != "date-time" {
		t.Errorf("Expected format to be 'date-time', got %v", result["format"])
	}
}

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

func TestApplyValidation(t *testing.T) {
	tests := []struct {
		name        string
		schema      map[string]interface{}
		validateTag string
		expected    map[string]interface{}
	}{
		{
			name:        "min and max for string",
			schema:      map[string]interface{}{"type": "string"},
			validateTag: "min=5,max=250",
			expected: map[string]interface{}{
				"type":      "string",
				"minLength": 5,
				"maxLength": 250,
			},
		},
		{
			name:        "min and max for integer",
			schema:      map[string]interface{}{"type": "integer"},
			validateTag: "min=1,max=100",
			expected: map[string]interface{}{
				"type":    "integer",
				"minimum": float64(1),
				"maximum": float64(100),
			},
		},
		{
			name:        "min and max for number",
			schema:      map[string]interface{}{"type": "number"},
			validateTag: "min=0,max=10.5",
			expected: map[string]interface{}{
				"type":    "number",
				"minimum": float64(0),
				"maximum": float64(10.5),
			},
		},
		{
			name:        "email format",
			schema:      map[string]interface{}{"type": "string"},
			validateTag: "email",
			expected: map[string]interface{}{
				"type":   "string",
				"format": "email",
			},
		},
		{
			name:        "url format",
			schema:      map[string]interface{}{"type": "string"},
			validateTag: "url",
			expected: map[string]interface{}{
				"type":   "string",
				"format": "uri",
			},
		},
		{
			name:        "uuid format",
			schema:      map[string]interface{}{"type": "string"},
			validateTag: "uuid",
			expected: map[string]interface{}{
				"type":   "string",
				"format": "uuid",
			},
		},
		{
			name:        "datetime format",
			schema:      map[string]interface{}{"type": "string"},
			validateTag: "datetime",
			expected: map[string]interface{}{
				"type":   "string",
				"format": "date-time",
			},
		},
		{
			name:        "pattern",
			schema:      map[string]interface{}{"type": "string"},
			validateTag: "pattern=^[a-z]+$",
			expected: map[string]interface{}{
				"type":    "string",
				"pattern": "^[a-z]+$",
			},
		},
		{
			name:        "multiple tags combined",
			schema:      map[string]interface{}{"type": "string"},
			validateTag: "required,min=5,max=250,email",
			expected: map[string]interface{}{
				"type":      "string",
				"minLength": 5,
				"maxLength": 250,
				"format":    "email",
			},
		},
		{
			name:        "empty validate tag",
			schema:      map[string]interface{}{"type": "string"},
			validateTag: "",
			expected: map[string]interface{}{
				"type": "string",
			},
		},
		{
			name:        "unknown tags are ignored",
			schema:      map[string]interface{}{"type": "string"},
			validateTag: "required,unknowntag",
			expected: map[string]interface{}{
				"type": "string",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyValidation(tt.schema, tt.validateTag)
			for key, expectedVal := range tt.expected {
				if result[key] != expectedVal {
					t.Errorf("applyValidation() key %s = %v, want %v", key, result[key], expectedVal)
				}
			}
		})
	}
}
