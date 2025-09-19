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
	switch fieldType {
	case "Bool":
		parsed, ok = boolParser(value)
	case "String":
		parsed, ok = stringParser(value)
	case "IPv4":
		parsed, ok = ipv4Parser(value)
	case "Int":
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

		if value == "" && optional {
			if hasDefault(field.Tag.Get("env")) {
				value = getDefault(field.Tag.Get("env"))
			} else {
				environment[field.Name] = envVarType{
					reflect.Zero(field.Type),
					field.Type.Kind(),
				}
				continue
			}
		}

		if value == "" && !optional {
			missing = append(missing, field.Name)
		}

		var ok error
		var parsed any
		kind := field.Type.Kind().String()
		if kind == "slice" {
			sliceOf := reflect.SliceOf(field.Type).String()
			sep := getSeparator(field.Tag.Get("env"))
			start := len("[][]env.")
			parsed, ok = validateAndParseSlice(field.Name, sliceOf[start:], value, sep)
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
