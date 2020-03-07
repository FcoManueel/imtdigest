package imthash

import "fmt"

type Hash struct {
	h [8]byte
}

var coefficients = [8]int{2, 3, 5, 7, 11, 13, 17, 19}

func (hash *Hash) Write(p []byte) (n int, err error) {
	for _, ib := range p {
		hash.h[0] = byte((int(ib) * coefficients[0]) % 255) // This will result in 0 for any 0 byte and cascade over the next items
		for i := 1; i < len(hash.h); i++ {
			hash.h[i] = byte(((int(hash.h[i-1]) + int(ib)) * coefficients[i]) % 255)
		}
	}
	return len(p), nil
}

func (hash *Hash) Hex() string {
	return fmt.Sprintf("%x", hash.h)
}
