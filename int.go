package main

import (
	s "strconv"
)

func intValidator(value string) bool {
	_, err := s.Atoi(value)
	return err == nil
}

func intParser(value string) int {
	casted, _ := s.Atoi(value)
	return casted
}

type Int int
