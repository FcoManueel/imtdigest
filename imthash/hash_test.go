package imthash

import (
	"testing"
)

func TestHash_Write(t *testing.T) {
	hash := &Hash{}
	hash.Write([]byte{12})
	expected := [8]byte{24, 108, 90, 204, 81, 189, 102, 126}
	if hash.h != expected {
		t.Errorf("hash.Write() = %v, want %v", hash.h, expected)
	}
}
