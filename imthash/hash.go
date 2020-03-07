package imthash

type Hash struct {
	h [8]byte
}

var coefficients = [8]int{2, 3, 5, 7, 11, 13, 17, 19}

func (hash *Hash) Write(p []byte) (n int, err error) {
	for _, ib := range p {
		hash.h[0] = byte((int(ib) * coefficients[0]) % 255)
		for i := 1; i < len(hash.h); i++ {
			hash.h[i] = byte(((int(hash.h[i-1]) + int(ib)) * coefficients[i]) % 255)
		}
	}
	return len(p), nil
}

func (h *Hash) Hex() string {
	// TODO Expose hash in hexadecimal
	return "0123456789abcdef"
}
