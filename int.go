package env

import (
	s "strconv"
)

func intParser(value string) (int, error) {
	return s.Atoi(value)
}

type Int int
