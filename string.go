package env

func stringValidator(value string) bool {
	return true
}

func stringParser(value string) string {
	return value
}

type String string
