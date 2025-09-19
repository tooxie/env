package env

import (
	err "errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type envVarType struct {
	Value reflect.Value
	Type  reflect.Kind
}
type envMapType map[string]envVarType
type invalidType struct {
	name  string
	value string
}

var envMap envMapType

func Get[T any](name string) T {
	envVar, ok := envMap[name]
	if !ok {
		panic(fmt.Sprintf(
			"Invalid key '%s': Can't access environment variables that were "+
				"not previously registered. Call `env.Assert()` first.", name))
	}

	// Handle slice types specially
	if envVar.Type == reflect.Slice {
		sliceValue := envVar.Value
		if sliceValue.Kind() == reflect.Slice {
			// Convert []interface{} to the target slice type
			result := reflect.MakeSlice(reflect.TypeOf((*T)(nil)).Elem(), sliceValue.Len(), sliceValue.Cap())
			for i := 0; i < sliceValue.Len(); i++ {
				elem := sliceValue.Index(i)
				if elem.CanInterface() {
					result.Index(i).Set(reflect.ValueOf(elem.Interface()))
				}
			}
			return result.Interface().(T)
		}
	}

	value, ok := envVar.Value.Interface().(T)
	if !ok {
		panic(fmt.Sprintf("Invalid type for key '%s': %T", name, value))
	}
	return value
}

func Assert(variables interface{}) error {
	missing, invalid := Validate(variables)

	var errors []string
	if missing != nil {
		errors = append(errors, fmt.Sprintf("Missing: %v", missing))
	}
	if invalid != nil {
		errors = append(errors, fmt.Sprintf("Invalid: %v", invalid))
	}

	if len(errors) > 0 {
		return err.New(strings.Join(errors, "\n"))
	}

	return nil
}

func parseVariable(fieldName string, fieldType string, value string) (any, error) {
	var ok error
	var parsed any
	normalizedFieldType := strings.ToLower(fieldType)
	switch normalizedFieldType {
	case "bool":
		parsed, ok = boolParser(value)
	case "string":
		parsed, ok = stringParser(value)
	case "ipv4":
		parsed, ok = ipv4Parser(value)
	case "int":
		parsed, ok = intParser(value)
	default:
		panic(fmt.Sprintf(
			"Unrecognized type '%s' for field '%s'", fieldType, fieldName))
	}

	return parsed, ok
}

func Validate(variables interface{}) ([]string, []invalidType) {
	environment := make(envMapType)
	t := reflect.TypeOf(variables)
	if t.Kind() != reflect.Struct {
		panic("Invalid parameter")
	}

	var missing []string
	var invalid []invalidType

	for n := 0; n < t.NumField(); n++ {
		field := t.Field(n)
		value := os.Getenv(field.Name)
		optional := isOptional(field.Tag.Get("env"))

		if value == "" {
			if optional {
				// If the field is optional, we can use the default value if it exists
				if hasDefault(field.Tag.Get("env")) {
					value = getDefault(field.Tag.Get("env"))
				} else {
					// If the field is optional and has no default value, we can use a zero value
					environment[field.Name] = envVarType{
						reflect.Zero(field.Type),
						field.Type.Kind(),
					}

					// We can continue to the next field, nothing to validate
					continue
				}
			} else {
				// If the field is required and has no value, we add it to the missing list
				missing = append(missing, field.Name)

				// We can continue to the next field, nothing to validate
				continue
			}
		}

		var ok error
		var parsed any
		kind := field.Type.Kind().String()
		if kind == "slice" {
			sep := getSeparator(field.Tag.Get("env"))
			// Get the element type from the slice
			elementType := field.Type.Elem()
			elementTypeName := elementType.Name()
			if elementTypeName == "" {
				elementTypeName = elementType.String()
			}
			parsed, ok = validateAndParseSlice(field.Name, elementTypeName, value, sep)
		} else {
			parsed, ok = parseVariable(field.Name, field.Type.Name(), value)
		}

		if ok != nil {
			invalid = append(invalid, invalidType{field.Name, value})
		} else {
			var varType reflect.Kind
			if kind == "slice" {
				varType = reflect.SliceOf(field.Type).Kind()
			} else {
				varType = field.Type.Kind()
			}
			environment[field.Name] = envVarType{
				reflect.ValueOf(parsed),
				varType,
			}
		}
	}

	envMap = environment
	return missing, invalid
}

func validateAndParseSlice(fieldName string, fieldType string, value string, sep string) ([]any, error) {
	var values []any
	var allOk = true
	for _, slice := range strings.Split(value, sep) {
		parsed, ok := parseVariable(fieldName, fieldType, slice)
		values = append(values, parsed)
		allOk = allOk && ok == nil
	}

	if !allOk {
		return values, fmt.Errorf("invalid slice: %v", allOk)
	}
	return values, nil
}
