package data

import (
	"errors"
	"io"
	"sync"
)

// ErrRingBufferClosed is returned when a Read/Pop or Write/Push is called on a [RingBuffer] that has already had its
// [*RingBuffer.Close] method called
var ErrRingBufferClosed = errors.New("buffer closed")

// RingBuffer is a list with a maximum length. When adding an element that would cause the list to grow past its maximum,
// items from the front of the list will be dropped, avoiding out of control list growth. It is safe for use across
// multiple goroutines
type RingBuffer[T any] struct {
	buf      []T
	capacity uint
	rw       *sync.Mutex
	closed   bool
}

// NewRingBuffer will return a ring buffer of the requested type and with the requested capacity. If capacity is 0,
// NewRingBuffer will panic. If you create NewRingBuffer[byte], the returned instance implements [io.ReadWriteCloser]
func NewRingBuffer[T any](capacity uint) *RingBuffer[T] {
	if capacity == 0 {
		panic("capacity must be > 0")
	}

	return &RingBuffer[T]{
		buf:      make([]T, 0, capacity),
		capacity: capacity,
		rw:       &sync.Mutex{},
	}
}

func (r *RingBuffer[T]) Read(p []T) (int, error) {
	if r.closed {
		return -1, ErrRingBufferClosed
	}
	r.rw.Lock()
	defer r.rw.Unlock()

	if len(p) == 0 {
		return 0, nil
	}

	if len(r.buf) == 0 {
		return 0, io.EOF
	}

	if len(p) > len(r.buf) {
		n := len(r.buf)
		copy(p, r.buf)
		clear(r.buf)
		return n, nil
	}

	n := len(p)
	copy(p, r.buf[:n])
	r.buf = r.buf[n:]

	return n, nil
}

func (r *RingBuffer[T]) Write(p []T) (int, error) {
	if r.closed {
		return -1, ErrRingBufferClosed
	}
	r.rw.Lock()
	defer r.rw.Unlock()

	if len(p) == 0 {
		return 0, nil
	}

	if len(p) >= int(r.capacity) {
		start := len(p) - int(r.capacity)
		r.buf = p[start:]

		return int(r.capacity), nil
	}

	r.buf = append(r.buf[len(p):int(r.capacity)], p...)
	return len(p), nil
}

// Close will render the RingBuffer unusable
func (r *RingBuffer[T]) Close() error {
	if r.closed {
		return ErrRingBufferClosed
	}

	r.closed = true
	r.buf = make([]T, 0, 1)
	r.capacity = 0

	return nil
}

// Pop is a convenience  method that wraps Read and returns a slice rather doing a pass-by-value
func (r *RingBuffer[T]) Pop(n int) ([]T, error) {
	if n < 1 {
		return nil, errors.New("n must be positive")
	}

	p := make([]T, n)
	x, err := r.Read(p)
	if err != nil {
		return nil, err
	}

	return p[:x], nil
}

// Push is a convenience method that wraps Write and allows the user to pass in a variadic list of items
func (r *RingBuffer[T]) Push(vals ...T) (int, error) {
	return r.Write(vals)
}

// Len is equivalent to len(slice)
func (r *RingBuffer[T]) Len() int {
	return len(r.buf)
}

// Cap returns the capacity of the buffer and is immutable
func (r *RingBuffer[T]) Cap() uint {
	return r.capacity
}
