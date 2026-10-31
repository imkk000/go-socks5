package bufferpool

import (
	"sync"
)

// BufPool is an interface for getting and returning temporary
// byte slices for use by io.CopyBuffer.
type BufPool interface {
	Get() []byte
	Put([]byte)
}

type pool struct {
	size int
	pool *sync.Pool
}

// NewPool new buffer pool for getting and returning temporary
// byte slices for use by io.CopyBuffer.
func NewPool(size int) BufPool {
	return &pool{
		size,
		&sync.Pool{
			New: func() any { return make([]byte, 0, size) }},
	}
}

// Get implement interface BufPool
func (sf *pool) Get() []byte {
	b, ok := sf.pool.Get().([]byte)
	if !ok {
		return make([]byte, 0, sf.size)
	}
	return b
}

// Put implement interface BufPool
func (sf *pool) Put(b []byte) {
	if cap(b) != sf.size {
		return // invalid buffer size that's put into leaky buffer
	}
	sf.pool.Put(b[:0]) //nolint:staticcheck
}
