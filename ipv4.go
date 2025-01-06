package main

import "net"

func ipv4Validator(value string) bool {
	ip := net.ParseIP(value)
	if ip == nil {
		return false
	}
	if ip.To4() == nil {
		return false
	}

	return true
}

func ipv4Parser(value string) string {
	return value
}

type IPv4 string
