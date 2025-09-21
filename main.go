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

// Assert validates environment variables and returns a populated struct instance
func Assert[T any](config T) (T, error) {
	missing, invalid := Validate(config)

	var errors []string
	if missing != nil {
		errors = append(errors, fmt.Sprintf("Missing: %v", missing))
	}
	if invalid != nil {
		errors = append(errors, fmt.Sprintf("Invalid: %v", invalid))
	}

	if len(errors) > 0 {
		var zero T
		return zero, err.New(strings.Join(errors, "\n"))
	}

	// Create a new instance of the struct and populate it with parsed values
	result := reflect.New(reflect.TypeOf(config)).Elem()

	for n := 0; n < result.NumField(); n++ {
		field := result.Field(n)
		fieldName := result.Type().Field(n).Name

		// Get the parsed value from envMap
		if envVar, ok := envMap[fieldName]; ok {
			// Handle slice types specially
			if envVar.Type == reflect.Slice {
				sliceValue := envVar.Value
				if sliceValue.Kind() == reflect.Slice {
					// Convert []interface{} to the target slice type
					result := reflect.MakeSlice(field.Type(), sliceValue.Len(), sliceValue.Cap())
					for i := 0; i < sliceValue.Len(); i++ {
						elem := sliceValue.Index(i)
						if elem.CanInterface() {
							// Convert the element to the correct type
							elemValue := reflect.ValueOf(elem.Interface())
							if elemValue.Type().ConvertibleTo(field.Type().Elem()) {
								result.Index(i).Set(elemValue.Convert(field.Type().Elem()))
							} else {
								// If direct conversion fails, try to parse as string first
								if elemValue.Type() == reflect.TypeOf("") {
									// Element is a string, parse it
									elemStr := elem.Interface().(string)
									parsed, err := parseVariable(fieldName, field.Type().Elem().Name(), elemStr)
									if err == nil {
										parsedValue := reflect.ValueOf(parsed)
										if parsedValue.Type().ConvertibleTo(field.Type().Elem()) {
											result.Index(i).Set(parsedValue.Convert(field.Type().Elem()))
										} else {
											// For custom string types like IPv4, create from the parsed string
											if field.Type().Elem().Kind() == reflect.String {
												customType := reflect.New(field.Type().Elem()).Elem()
												customType.SetString(parsed.(string))
												result.Index(i).Set(customType)
											} else {
												result.Index(i).Set(parsedValue)
											}
										}
									} else {
										result.Index(i).Set(elemValue)
									}
								} else {
									result.Index(i).Set(elemValue)
								}
							}
						}
					}
					field.Set(result)
				}
			} else {
				// For non-slice types, set the value directly
				// But first check if we need to convert custom types
				if envVar.Value.Type().ConvertibleTo(field.Type()) {
					field.Set(envVar.Value.Convert(field.Type()))
				} else {
					// For custom string types like IPv4, create from the parsed string
					if field.Type().Kind() == reflect.String && envVar.Value.Type() == reflect.TypeOf("") {
						customType := reflect.New(field.Type()).Elem()
						customType.SetString(envVar.Value.Interface().(string))
						field.Set(customType)
					} else {
						field.Set(envVar.Value)
					}
				}
			}
		}
	}

	return result.Interface().(T), nil
}

// MustAssert validates environment variables and returns a populated struct instance
// It panics if validation fails, making it convenient for the common use case
func MustAssert[T any](config T) T {
	result, err := Assert(config)
	if err != nil {
		panic(fmt.Sprintf("Configuration error: %v", err))
	}
	return result
}

// Parses the value of the environment variable into the correct type
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
	case "url":
		parsed, ok = urlParser(value)
	case "httpurl":
		parsed, ok = httpURLParser(value)
	default:
		panic(fmt.Sprintf(
			"Unrecognized type '%s' for field '%s'", fieldType, fieldName))
	}

	return parsed, ok
}

// Checks to see if the field has a custom environment variable name, if not
// it returns the field name. For example, if the field is `DatabaseURL` and the
// environment variable name is `DB_URL`, it will return `DB_URL`.
func getEnvVarNameFromField(field reflect.StructField) string {
	name := getName(field.Tag.Get("env"))
	if name == "" {
		name = field.Name
	}

	return name
}

// Validates the environment variables and returns a list of missing and invalid
// variables. If the value is valid, it will be added to the environment map.
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
		name := strings.ToUpper(getEnvVarNameFromField(field))
		value := os.Getenv(name)
		optional := isOptional(field.Tag.Get("env"))

		if value == "" {
			if optional {
				// If the field is optional, we can use the default value if it exists
				if hasDefault(field.Tag.Get("env")) {
					value = getDefault(field.Tag.Get("env"))
				} else {
					// If the field is optional and has no default value, we can use a zero value
					environment[field.Name] = envVarType{
						reflect.Value{},
						field.Type.Kind(),
					}

					// We can continue to the next field, nothing to validate
					continue
				}
			} else {
				// If the field is required and has no value, we add it to the missing list
				missing = append(missing, fmt.Sprintf("%s (%s)", field.Name, name))

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
			invalid = append(invalid, invalidType{fmt.Sprintf("%s (%s)", field.Name, name), value})
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

// Given a slice field and its corresponding environment variable value, it will
// parse the value into the correct type and return a slice of the parsed values.
// If the value is invalid, it will return an error.
func validateAndParseSlice(fieldName string, fieldType string, value string, sep string) ([]any, error) {
	var values []any
	var allOk = true
	for slice := range strings.SplitSeq(value, sep) {
		parsed, ok := parseVariable(fieldName, fieldType, slice)
		values = append(values, parsed)
		allOk = allOk && ok == nil
	}

	if !allOk {
		return values, fmt.Errorf("invalid slice: %v", allOk)
	}
	return values, nil
}
