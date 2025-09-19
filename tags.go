package env

import (
	"regexp"
	"strings"
)

const defaultSeparator = " "

var (
	defaultRegex   = regexp.MustCompile("default='(?P<Default>.*?)'")
	separatorRegex = regexp.MustCompile("separator='(?P<Sep>.)'")
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
