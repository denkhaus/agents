package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateJSON performs JSON schema validation on the given data.
// It first transforms the data to resolve methods that act as fields
// and converts types implementing fmt.Stringer to their string representation.
// It returns the processed data, the validation result, and an error if any.
func ValidateJSON(data interface{}, schema *gojsonschema.Schema) (interface{}, *gojsonschema.Result, error) {
	processedData, err := resolveMethodsInStruct(data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to resolve methods in data: %w", err)
	}

	dataLoader := gojsonschema.NewGoLoader(processedData)
	result, err := schema.Validate(dataLoader)
	if err != nil {
		return nil, nil, fmt.Errorf("json schema validation failed: %w", err)
	}

	return processedData, result, nil
}

// resolveMethodsInStruct recursively traverses a Go data structure (struct, slice, map)
// and resolves methods that have no arguments and return a single value, treating them as fields.
// It also converts types implementing fmt.Stringer to their string representation.
// It returns a new data structure (map[string]interface{} or []interface{}) with resolved methods.
func resolveMethodsInStruct(input interface{}) (interface{}, error) {
	val := reflect.ValueOf(input)

	// Handle nil input
	if !val.IsValid() || (val.Kind() == reflect.Ptr && val.IsNil()) {
		return nil, nil
	}

	// Store the original value for method lookup
	originalVal := val

	// Dereference pointer if applicable for field processing
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// If it's a struct and not addressable, create an addressable copy
	if val.Kind() == reflect.Struct && !val.CanAddr() {
		newVal := reflect.New(val.Type())
		newVal.Elem().Set(val)
		val = newVal.Elem() // Use the Elem() of the new pointer to get the addressable struct value
	}

	switch val.Kind() {
	case reflect.Struct:
		output := make(map[string]interface{})
		typ := val.Type() // This typ is for the value type (AgentInfo)

		// Use originalVal.Type() for methods, as methods might be on the pointer type
		methodTyp := originalVal.Type()
		resolvedValues := make(map[string]interface{})

		// Process exported fields
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			fieldType := typ.Field(i)

			if !fieldType.IsExported() {
				continue // Skip unexported fields
			}

			// Handle embedded structs
			if field.Kind() == reflect.Struct && fieldType.Anonymous {
				embeddedResolved, err := resolveMethodsInStruct(field.Interface())
				if err != nil {
					return nil, err
				}
				// Merge embedded struct's resolved fields into the current output
				if embeddedMap, ok := embeddedResolved.(map[string]interface{}); ok {
					for k, v := range embeddedMap {
						resolvedValues[k] = v
					}
				}
				continue // Skip normal processing for embedded struct
			}

			fieldName := fieldType.Name
			jsonTag := fieldType.Tag.Get("json")
			if jsonTag != "" {
				parts := strings.Split(jsonTag, ",")
				fieldName = parts[0]
				if fieldName == "-" {
					continue
				}
			}
			resolvedValues[fieldName] = field.Interface()
		}

		// Process exported methods
		for i := 0; i < methodTyp.NumMethod(); i++ {
			method := methodTyp.Method(i)
			// Check if method takes no arguments (besides receiver) and returns one value
			if method.Type.NumIn() == 1 && method.Type.NumOut() == 1 {
				// Prioritize method over field if names clash
				results := originalVal.MethodByName(method.Name).Call([]reflect.Value{})
				if len(results) > 0 {
					resolvedValue := results[0].Interface()

					// Check if the resolvedValue implements fmt.Stringer
					if stringer, ok := resolvedValue.(fmt.Stringer); ok {
						resolvedValue = stringer.String()
					}

					resolvedValues[method.Name] = resolvedValue
				}
			}
		}

		// Recursively resolve methods for collected values
		for key, value := range resolvedValues {
			resolvedValue, err := resolveMethodsInStruct(value)
			if err != nil {
				return nil, err
			}
			output[key] = resolvedValue
		}
		return output, nil

	case reflect.Slice, reflect.Array:
		output := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			resolvedItem, err := resolveMethodsInStruct(val.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			output[i] = resolvedItem
		}
		return output, nil

	case reflect.Map:
		output := make(map[string]interface{})
		for _, key := range val.MapKeys() {
			resolvedValue, err := resolveMethodsInStruct(val.MapIndex(key).Interface())
			if err != nil {
				return nil, err
			}
			output[key.String()] = resolvedValue
		}
		return output, nil

	default:
		// Return primitives and other types as is
		return input, nil
	}
}
