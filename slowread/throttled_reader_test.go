package slowread

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
	"time"
)

const KB = 1 << 10
const MB = KB << 10

type testCloser struct {
	closed bool
}

func (tc *testCloser) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (tc *testCloser) Close() error {
	tc.closed = true
	return nil
}

func Test_reader_Close(t *testing.T) {
	tc := &testCloser{}
	sr := &slowReader{
		r: tc,
	}
	if err := sr.Close(); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !tc.closed {
		t.Errorf("Wrapped Reader was not closed")
	}
}

func Test_reader_Close__noop(t *testing.T) {
	sr := &slowReader{}
	if err := sr.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func Test_reader_Read(t *testing.T) {
	tests := []struct {
		rate  int
		total int
	}{
		{
			rate:  1,
			total: 2,
		},
		{
			rate:  KB,
			total: KB,
		},
		{
			rate:  KB,
			total: 256,
		},
		{
			rate:  KB,
			total: 32,
		},
		{
			rate:  MB,
			total: 256 * KB,
		},
		{
			rate:  MB,
			total: 32 * KB,
		},
		{
			rate:  MB,
			total: 256,
		},
		{
			rate:  MB,
			total: 32,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%db at %db/s", tt.total, tt.rate), func(t *testing.T) {
			r := bytes.NewReader(bytes.Repeat([]byte{255}, tt.total))
			sr := NewReader(r, tt.rate)

			// measure how long it takes to read dummy data and see if it's being throttled correctly
			start := time.Now()
			nRead, err := io.Copy(ioutil.Discard, sr)
			elapsed := time.Now().Sub(start)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if nRead != int64(tt.total) {
				t.Errorf("Total bytes read = %d, want %d", nRead, tt.total)
			}

			effectiveRate := float64(tt.total) / elapsed.Seconds()
			if effectiveRate > float64(tt.rate) {
				t.Errorf("Effective read rate = %.2f b/s, maximum wanted %v b/s", effectiveRate, tt.rate)
			}
		})
	}
}

func Test_reader_Read__nothrottling(t *testing.T) {
	tests := []struct {
		rate  int
		total int
	}{
		{
			rate:  0,
			total: 1,
		},
		{
			rate:  0,
			total: KB,
		},
		{
			rate:  0,
			total: MB,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%db with no throttling", tt.total), func(t *testing.T) {
			r := bytes.NewReader(bytes.Repeat([]byte{255}, tt.total))
			sr := NewReader(r, tt.rate)

			start := time.Now()
			nRead, err := io.Copy(ioutil.Discard, sr)
			elapsed := time.Now().Sub(start)

			if err != nil {
				t.Errorf("Read() unexpected error = %v", err)
				return
			}
			if nRead != int64(tt.total) {
				t.Errorf("Total bytes read = %v, want %v", nRead, tt.total)
			}
			// In a healthy environment this should take less than a millisecond. Compare to a sufficiently big value to avoid false positives
			if elapsed.Milliseconds() > 50 {
				t.Errorf("Reading %v bytes lasted %d milliseconds, it's likely being throttled", tt.total, elapsed.Milliseconds())
			}
		})
	}
}
