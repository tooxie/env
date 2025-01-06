package main

import (
	"regexp"
	"strings"
)

const defaultSeparator = " "

func toLower(tag string) string {
	return strings.ToLower(tag)
}

func isOptional(tag string) bool {
	return strings.Contains(toLower(tag), "optional")
}

func hasDefault(tag string) bool {
	r := regexp.MustCompile("default='(?P<Default>.*?)'")
	m := r.FindAllStringSubmatch(tag, -1)
	return len(m) > 0
}

func getDefault(tag string) string {
	r := regexp.MustCompile("default='(?P<Default>.*?)'")
	m := r.FindAllStringSubmatch(tag, -1)
	return m[0][1]
}

func getSeparator(tag string) string {
	r := regexp.MustCompile("separator='(?P<Sep>.)'")
	m := r.FindAllStringSubmatch(tag, -1)

	if len(m) == 0 {
		return defaultSeparator
	}

	if len(m) != 1 {
		panic("Too many separators in tag")
	}

	return m[0][1]
}
