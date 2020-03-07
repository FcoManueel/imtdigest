package imthash

import (
	"fmt"
	"testing"
)

func TestHash_Write(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  [8]byte
	}{
		{
			name:  "nil input",
			input: nil,
			want:  [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:  "empty input",
			input: []byte{},
			want:  [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:  "provided example",
			input: []byte{12},
			want:  [8]byte{24, 108, 90, 204, 81, 189, 102, 126},
		},
		{
			// Leading 0s don't affect the hash
			name:  "leading 0s",
			input: []byte{0, 0, 0, 12},
			want:  [8]byte{24, 108, 90, 204, 81, 189, 102, 126},
		},
		{
			// Any 0 byte in the input effectively flushes the hash collected so far
			name:  "middle 0s",
			input: []byte{1, 2, 3, 4, 5, 0, 12},
			want:  [8]byte{24, 108, 90, 204, 81, 189, 102, 126},
		},
		{
			// Because a 0 flushes the hash, a traling 0 always results in an all-0 byte array
			name:  "trailing 0",
			input: []byte{12, 0},
			want:  [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := &Hash{}
			n, err := hash.Write(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if n != len(tt.input) {
				t.Errorf("hash.Write() read %d bytes, want %v", n, len(tt.input))
			}
			if hash.h != tt.want {
				t.Errorf("hash.Write() = %v, want %v", hash.h, tt.want)
			}
		})
	}
}

func TestHash_Hex(t *testing.T) {
	tests := []struct {
		input [8]byte
		want  string
	}{
		{
			input: [8]byte{0, 0, 0, 0, 0, 0, 0, 0},
			want:  "0000000000000000",
		},
		{
			input: [8]byte{255, 255, 255, 255, 255, 255, 255, 255},
			want:  "ffffffffffffffff",
		},
		{
			input: [8]byte{0, 1, 2, 3, 4, 5, 6, 7},
			want:  "0001020304050607",
		},
		{
			input: [8]byte{8, 9, 10, 11, 12, 13, 14, 15},
			want:  "08090a0b0c0d0e0f",
		},
		{
			input: [8]byte{24, 108, 90, 204, 81, 189, 102, 126},
			want:  "186c5acc51bd667e",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i+1), func(t *testing.T) {
			hash := &Hash{
				h: tt.input,
			}
			hex := hash.Hex()
			if hex != tt.want {
				t.Errorf("Hex() = %v, want %v", hex, tt.want)
			}
		})
	}
}
