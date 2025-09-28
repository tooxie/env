package env

import (
	"regexp"
	"strings"
)

const defaultSeparator = " "

var (
	defaultRegex   = regexp.MustCompile("default='(?P<Default>.*?)'")
	separatorRegex = regexp.MustCompile("separator='(?P<Sep>.)'")
	nameRegex      = regexp.MustCompile("name='(?P<Name>.*?)'")
	valuesRegex    = regexp.MustCompile("values='(?P<Values>.*?)'")
)

func toLower(tag string) string {
	return strings.ToLower(tag)
}

func isOptional(tag string) bool {
	return strings.Contains(toLower(tag), "optional")
}

func hasDefault(tag string) bool {
	m := defaultRegex.FindAllStringSubmatch(tag, -1)
	return len(m) > 0
}

func getDefault(tag string) string {
	m := defaultRegex.FindAllStringSubmatch(tag, -1)
	return m[0][1]
}

func getSeparator(tag string) string {
	m := separatorRegex.FindAllStringSubmatch(tag, -1)

	if len(m) == 0 {
		return defaultSeparator
	}

	if len(m) != 1 {
		panic("Too many separators in tag")
	}

	return m[0][1]
}

func getName(tag string) string {
	m := nameRegex.FindAllStringSubmatch(tag, -1)

	if len(m) == 0 {
		return ""
	}

	if len(m) != 1 {
		panic("Too many name specifications in tag")
	}

	return m[0][1]
}

func hasValues(tag string) bool {
	m := valuesRegex.FindAllStringSubmatch(tag, -1)
	return len(m) > 0
}

func getValues(tag string) []string {
	m := valuesRegex.FindAllStringSubmatch(tag, -1)

	if len(m) == 0 {
		return nil
	}

	if len(m) != 1 {
		panic("Too many values specifications in tag")
	}

	valuesStr := m[0][1]
	// Split by comma and trim whitespace
	values := strings.Split(valuesStr, ",")
	for i, v := range values {
		values[i] = strings.TrimSpace(v)
	}

	return values
}
