package openapi

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ProcedureInfo struct {
	Path       string
	Method     string
	Route      interface{}
	Meta       interface{}
	Tags       []string
	InputType  reflect.Type
	OutputType reflect.Type
	ErrorCodes []int
	PathParams []string
}

// GenerateSpec generates an OpenAPI 3.0 specification from procedures
func GenerateSpec(procedures []ProcedureInfo) (map[string]interface{}, error) {

	spec := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]interface{}{
			"title":   "API Documentation",
			"version": "1.0.0",
		},
		"paths": make(map[string]interface{}),
	}

	paths := spec["paths"].(map[string]interface{})

	for _, proc := range procedures {
		path := proc.Path
		if path == "" && proc.Route != nil {
			routeValue := reflect.ValueOf(proc.Route)
			if routeValue.Kind() == reflect.Ptr && !routeValue.IsNil() {
				routeValue = routeValue.Elem()
			}
			if routeValue.IsValid() {
				pathField := routeValue.FieldByName("Path")
				if pathField.IsValid() && pathField.Kind() == reflect.String {
					path = pathField.String()
				}
			}
		}
		if path != "" && path[0] != '/' {
			path = "/" + path
		}

		// Convert :param to {param} for OpenAPI spec
		path = convertToOpenAPIPath(path)

		if path == "" {
			continue
		}

		method := proc.Method
		if method == "" && proc.Route != nil {
			routeValue := reflect.ValueOf(proc.Route)
			if routeValue.Kind() == reflect.Ptr && !routeValue.IsNil() {
				routeValue = routeValue.Elem()
			}
			if routeValue.IsValid() {
				methodField := routeValue.FieldByName("Method")
				if methodField.IsValid() && methodField.Kind() == reflect.String {
					method = methodField.String()
				}
			}
		}
		if method == "" {
			continue
		}
		method = toLower(method)

		pathItem, exists := paths[path]
		if !exists {
			pathItem = make(map[string]interface{})
			paths[path] = pathItem
		}

		pathItemMap := pathItem.(map[string]interface{})

		operation := map[string]interface{}{}

		if len(proc.PathParams) > 0 {
			params := []interface{}{}
			for _, param := range proc.PathParams {
				params = append(params, map[string]interface{}{
					"name":     param,
					"in":       "path",
					"required": true,
					"schema": map[string]interface{}{
						"type": "integer",
					},
				})
			}
			operation["parameters"] = params
		}

		metaValue := reflect.ValueOf(proc.Meta)
		if metaValue.IsValid() {
			summaryField := metaValue.FieldByName("Summary")
			if summaryField.IsValid() && summaryField.Kind() == reflect.String {
				if summary := summaryField.String(); summary != "" {
					operation["summary"] = summary
				}
			}
			descField := metaValue.FieldByName("Description")
			if descField.IsValid() && descField.Kind() == reflect.String {
				if desc := descField.String(); desc != "" {
					operation["description"] = desc
				}
			}
		}
		if len(proc.Tags) > 0 {
			operation["tags"] = proc.Tags
		}

		// Only add requestBody if input type has fields (not an empty struct{})
		if proc.InputType != nil && proc.InputType.Kind() == reflect.Struct && proc.InputType.NumField() > 0 {
			requestBody := map[string]interface{}{
				"required": true,
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": typeToSchema(proc.InputType),
					},
				},
			}
			operation["requestBody"] = requestBody
		}

		responses := map[string]interface{}{}

		if method == "post" {
			responses["201"] = map[string]interface{}{
				"description": "Created",
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": typeToSchema(proc.OutputType),
					},
				},
			}
		} else {
			responses["200"] = map[string]interface{}{
				"description": "Success",
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": typeToSchema(proc.OutputType),
					},
				},
			}
		}

		errorCodes := proc.ErrorCodes
		has400 := false
		for _, code := range errorCodes {
			if code == 400 {
				has400 = true
				break
			}
		}
		if !has400 && proc.InputType != nil && proc.InputType.Kind() != reflect.Invalid {
			errorCodes = append([]int{400}, errorCodes...)
		}

		for _, code := range errorCodes {
			codeStr := fmt.Sprintf("%d", code)
			description := getErrorDescription(code)
			responses[codeStr] = map[string]interface{}{
				"description": description,
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": map[string]interface{}{
							"type": "object",
							"properties": map[string]interface{}{
								"message": map[string]interface{}{
									"type": "string",
								},
							},
							"required": []string{"message"},
						},
					},
				},
			}
		}

		operation["responses"] = responses

		pathItemMap[method] = operation
	}

	return spec, nil
}

// typeToSchema converts a Go reflect.Type to an OpenAPI JSON schema
func typeToSchema(t reflect.Type) map[string]interface{} {
	if t == nil {
		return map[string]interface{}{"type": "object"}
	}

	schema := make(map[string]interface{})

	switch t.Kind() {
	case reflect.String:
		schema["type"] = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		schema["type"] = "integer"
		schema["format"] = "int64"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		schema["type"] = "integer"
		schema["format"] = "int64"
	case reflect.Float32, reflect.Float64:
		schema["type"] = "number"
		schema["format"] = "double"
	case reflect.Bool:
		schema["type"] = "boolean"
	case reflect.Slice, reflect.Array:
		schema["type"] = "array"
		if t.Elem() != nil {
			schema["items"] = typeToSchema(t.Elem())
		}
	case reflect.Map:
		schema["type"] = "object"
		if t.Key().Kind() == reflect.String {
			if t.Elem() != nil {
				schema["additionalProperties"] = typeToSchema(t.Elem())
			}
		}
	case reflect.Ptr:
		innerSchema := typeToSchema(t.Elem())
		innerSchema["nullable"] = true
		return innerSchema
	case reflect.Struct:
		if t == reflect.TypeOf(time.Time{}) {
			schema["type"] = "string"
			schema["format"] = "date-time"
			return schema
		}
		schema["type"] = "object"
		properties := make(map[string]interface{})
		required := []string{}

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}

			fieldName := field.Name
			jsonTag := field.Tag.Get("json")
			if jsonTag != "" && jsonTag != "-" {
				if idx := indexOf(jsonTag, ','); idx != -1 {
					jsonTag = jsonTag[:idx]
				}
				if jsonTag != "" {
					fieldName = jsonTag
				}
			}

			fieldSchema := typeToSchema(field.Type)
			properties[fieldName] = fieldSchema

			validateTag := field.Tag.Get("validate")
			if validateTag != "" {
				fieldSchema = applyValidation(fieldSchema, validateTag)
				properties[fieldName] = fieldSchema
			}

			if field.Type.Kind() != reflect.Ptr {
				required = append(required, fieldName)
			}
		}

		if len(properties) > 0 {
			schema["properties"] = properties
		}
		if len(required) > 0 {
			schema["required"] = required
		}
	default:
		schema["type"] = "object"
	}

	return schema
}

func toLower(s string) string {
	if len(s) == 0 {
		return s
	}
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func indexOf(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

func getErrorDescription(code int) string {
	switch code {
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 403:
		return "Forbidden"
	case 404:
		return "Not Found"
	case 409:
		return "Conflict"
	case 422:
		return "Unprocessable Entity"
	case 500:
		return "Internal Server Error"
	case 502:
		return "Bad Gateway"
	case 503:
		return "Service Unavailable"
	default:
		return fmt.Sprintf("Error %d", code)
	}
}

func convertToOpenAPIPath(path string) string {
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		if len(segment) > 0 && segment[0] == ':' {
			segments[i] = "{" + segment[1:] + "}"
		}
	}
	return strings.Join(segments, "/")
}

func applyValidation(schema map[string]interface{}, validateTag string) map[string]interface{} {
	tags := strings.Split(validateTag, ",")
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)

		switch {
		case strings.HasPrefix(tag, "min="):
			val := strings.TrimPrefix(tag, "min=")
			if minVal, err := strconv.ParseFloat(val, 64); err == nil {
				switch schema["type"] {
				case "string":
					schema["minLength"] = int(minVal)
				case "integer", "number":
					schema["minimum"] = minVal
				}
			}
		case strings.HasPrefix(tag, "max="):
			val := strings.TrimPrefix(tag, "max=")
			if maxVal, err := strconv.ParseFloat(val, 64); err == nil {
				switch schema["type"] {
				case "string":
					schema["maxLength"] = int(maxVal)
				case "integer", "number":
					schema["maximum"] = maxVal
				}
			}
		case tag == "email":
			schema["format"] = "email"
		case tag == "url":
			schema["format"] = "uri"
		case tag == "uuid":
			schema["format"] = "uuid"
		case tag == "datetime":
			schema["format"] = "date-time"
		case strings.HasPrefix(tag, "pattern="):
			val := strings.TrimPrefix(tag, "pattern=")
			schema["pattern"] = val
		}
	}
	return schema
}
