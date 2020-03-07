package main

import "fmt"

type Hash struct{}

func (h *Hash) Write(p []byte) (n int, err error) {
	// TODO Implement stream hasher
	fmt.Print(string(p))
	return len(p), nil
}

func (h *Hash) Hex() string {
	// TODO Expose hash in hexadecimal
	return "0123456789abcdef"
}
