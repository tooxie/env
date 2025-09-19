package env

import "fmt"

func in(needle string, haystack []string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}
	return false
}

func boolValidator(value string) error {
	booleanValues := []string{
		"true", "false",
		"yes", "no",
		"1", "0",
	}
	if !in(value, booleanValues) {
		return fmt.Errorf("invalid boolean value: %s", value)
	}
	return nil
}

func boolParser(value string) (bool, error) {
	ok := boolValidator(value)
	casted := in(value, []string{"true", "yes", "1"})
	return casted, ok
}

type Bool bool
