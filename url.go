package env

import (
	"fmt"
	"net/url"
	"strings"
)

// Type: URL
func urlValidator(value string) (*url.URL, error) {
	originalValue := strings.TrimSpace(value)

	if originalValue == "" {
		return nil, fmt.Errorf("empty URL")
	}

	u, err := url.ParseRequestURI(originalValue)
	if err != nil {
		return nil, fmt.Errorf("invalid URL format: %v", err)
	}

	// We return the parsed URL so that `httpURLValidator` can check the
	// protocol and does not need to parse the URL again
	return u, nil
}

func urlParser(value string) (string, error) {
	_, err := urlValidator(value)
	if err != nil {
		return "", err
	}

	return value, nil
}

type URL string

// Type: HTTPURL
func httpURLValidator(value string) error {
	u, err := urlValidator(value)
	if err != nil {
		return err
	}

	// Check if protocol is HTTP or HTTPS only
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("HTTPURL only accepts http:// or https:// protocols, got: %s://", u.Scheme)
	}

	return nil
}

func httpURLParser(value string) (string, error) {
	ok := httpURLValidator(value)
	if ok != nil {
		return "", ok
	}

	// Return the validated URL as-is
	return value, nil
}

type HTTPURL string
