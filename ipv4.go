package env

import (
	"fmt"
	"net"
)

func ipv4Validator(value string) error {
	ip := net.ParseIP(value)
	if ip == nil {
		return fmt.Errorf("invalid IPv4 address: %s", value)
	}
	if ip.To4() == nil {
		return fmt.Errorf("invalid IPv4 address: %s", value)
	}

	return nil
}

func ipv4Parser(value string) (string, error) {
	ok := ipv4Validator(value)
	if ok != nil {
		return "", ok
	}
	return value, nil
}

type IPv4 string
