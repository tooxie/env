package env

import (
	"fmt"
	"log"
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

	value := envVar.Value.Interface().(T)
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
		return fmt.Errorf(strings.Join(errors, "\n"))
	}

	return nil
}

func runValidatorAndParser(fieldName string, fieldType string, value string) (bool, any) {
	var isValid bool
	var parsed any
	switch fieldType {
	case "Bool":
		isValid = boolValidator(value)
		if isValid {
			parsed = boolParser(value)
		}
	case "String":
		isValid = stringValidator(value)
		if isValid {
			parsed = stringParser(value)
		}
	case "IPv4":
		isValid = ipv4Validator(value)
		if isValid {
			parsed = ipv4Parser(value)
		}
	case "Int":
		isValid = intValidator(value)
		if isValid {
			parsed = intParser(value)
		}
	default:
		panic(fmt.Sprintf(
			"Unrecognized type '%s' for field '%s'", fieldType, fieldName))
	}

	return isValid, parsed
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
			fmt.Println("OPTIONAL", field.Tag.Get("env"))
			if hasDefault(field.Tag.Get("env")) {
				fmt.Println("Has default!")
				value = getDefault(field.Tag.Get("env"))
				fmt.Println("New value:", value)
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

		var isValid bool
		var parsed any
		kind := field.Type.Kind().String()
		if kind == "slice" {
			sliceOf := reflect.SliceOf(field.Type).String()
			sep := getSeparator(field.Tag.Get("env"))
			start := len("[][]env.")
			isValid, parsed = validateAndParseSlice(field.Name, sliceOf[start:], value, sep)
		} else {
			isValid, parsed = runValidatorAndParser(field.Name, field.Type.Name(), value)
		}

		if !isValid {
			invalid = append(invalid, invalidType{field.Name, value})
		} else {
			if kind == "slice" {
				environment[field.Name] = envVarType{
					reflect.ValueOf(parsed),
					reflect.SliceOf(field.Type).Kind(),
				}
			} else {
				environment[field.Name] = envVarType{
					reflect.ValueOf(parsed),
					field.Type.Kind(),
				}
			}
		}
	}

	envMap = environment
	log.Println(envMap)
	return missing, invalid
}

func validateAndParseSlice(fieldName string, fieldType string, value string, sep string) (bool, []any) {
	var values []any
	isAllValid := true
	for _, slice := range strings.Split(value, sep) {
		isValid, parsed := runValidatorAndParser(fieldName, fieldType, slice)
		values = append(values, parsed)
		isAllValid = isAllValid && isValid
	}

	return isAllValid, values
}
