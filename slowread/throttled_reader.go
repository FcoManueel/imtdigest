package slowread

import (
	"context"
	"io"
	"math"
	"time"

	"golang.org/x/time/rate"
)

type slowReader struct {
	r       io.Reader
	limiter *rate.Limiter
}

// We use a large burst size because we don't want to limit the speed at which the slow reader receives
// data, but rather the speed at which the slow reader sends that data to the wrapped reader
const burstSize = math.MaxInt32

// NewReader returns a new io.ReaderCloser that wraps r and Reads at a max rate of bytesPerSecond
func NewReader(r io.Reader, bytesPerSecond int) io.ReadCloser {
	var limiter *rate.Limiter
	if bytesPerSecond > 0 {
		limiter = rate.NewLimiter(rate.Limit(float64(bytesPerSecond)), burstSize)
		limiter.AllowN(time.Now(), burstSize)
	}
	return &slowReader{
		r:       r,
		limiter: limiter,
	}
}

// Read is the throttled version of the wrapped Reader
func (s *slowReader) Read(p []byte) (int, error) {
	if s.limiter == nil {
		return s.r.Read(p)
	}
	n, err := s.r.Read(p)
	if err != nil {
		return n, err
	}
	// consume n tokens from our bucket, where n is the amount of bytes read
	if err := s.limiter.WaitN(context.Background(), n); err != nil {
		return n, err
	}
	return n, nil
}

// Close defers the Close call to the wrapped reader. If the wrapped reader has no Close method, then it's a noop
func (s *slowReader) Close() error {
	if c, ok := s.r.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
