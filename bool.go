package env

func in(needle string, haystack []string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}
	return false
}

func boolValidator(value string) bool {
	booleanValues := []string{
		"true", "false",
		"yes", "no",
		"1", "0",
	}
	return in(value, booleanValues)
}

func boolParser(value string) bool {
	return in(value, []string{"true", "yes", "1"})
}

type Bool bool
