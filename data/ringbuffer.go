package data

import (
	"errors"
	"io"
	"sync"
	"time"
)

// ErrRingBufferClosed is returned when a Read/Pop or Write/Push is called on a [RingBuffer] that has already had its
// [*RingBuffer.Close] method called
var ErrRingBufferClosed = errors.New("buffer closed")

// RingBuffer is a list with a maximum length. When adding an element that would cause the list to grow past its maximum,
// items from the front of the list will be dropped, avoiding out of control list growth. It is safe for use across
// multiple goroutines
type RingBuffer[T any] struct {
	buf      []T
	capacity int
	rw       *sync.RWMutex
	closed   bool
}

// NewRingBuffer will return a ring buffer of the requested type and with the requested capacity. If capacity is 0,
// NewRingBuffer will panic. If you create NewRingBuffer[byte], the returned instance implements [io.ReadWriteCloser]
func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	if capacity == 0 {
		panic("capacity must be > 0")
	}

	return &RingBuffer[T]{
		buf:      make([]T, 0, capacity),
		capacity: capacity,
		rw:       &sync.RWMutex{},
	}
}

func (r *RingBuffer[T]) isClosed() bool {
	r.rw.RLock()
	defer r.rw.RUnlock()
	return r.closed
}

// Read will remove and return the first len(p) items from the buffer. If the buffer is closed, it will return
// [ErrRingBufferClosed]. If the buffer is empty, it will block until data becomes avaailable, and it is the responsibilty
// of the user to Close the buffer if they wish to end the wait. If Close is called while a Read call is waiting, it will
// return -1, [io.EOF]. In the successful case, it will return the number of items returned (which will be min(len(p), buffer.capacity)) and a nil error
func (r *RingBuffer[T]) Read(p []T) (int, error) {
	if r.isClosed() {
		return -1, ErrRingBufferClosed
	}

	if len(p) == 0 {
		return 0, nil
	}

	bufLen := r.Len()
	for bufLen == 0 {
		time.Sleep(5 * time.Millisecond)
		if r.isClosed() {
			return -1, io.EOF
		}
		bufLen = r.Len()
	}

	if len(p) > bufLen {
		n := bufLen
		r.rw.RLock()
		copy(p, r.buf)
		r.rw.RUnlock()
		r.rw.Lock()
		clear(r.buf)
		r.rw.Unlock()
		return n, nil
	}

	n := len(p)
	r.rw.RLock()
	copy(p, r.buf[:n])
	r.rw.RUnlock()
	r.rw.Lock()
	r.buf = r.buf[n:]
	r.rw.Unlock()

	return n, nil
}

// Write will push items into the buffer. If the buffer is closed, it will return [ErrRingBufferClosed]. That is the only
// case in which Write will return an error. If the given list of items would cause the buffer to grow, a sufficient number
// of items will be dropped from the front of the buffer. If the list of items is larger than the buffer itself, then
// when it is complete, the list will be the last buffer.capacity items of the passed list
func (r *RingBuffer[T]) Write(p []T) (int, error) {
	if r.isClosed() {
		return -1, ErrRingBufferClosed
	}

	if len(p) == 0 {
		return 0, nil
	}

	if len(p) >= r.Cap() {
		start := len(p) - r.Cap()
		r.rw.Lock()
		r.buf = p[start:]
		r.rw.Unlock()

		return r.Cap(), nil
	}

	availableSpace := r.Cap() - r.Len()
	if len(p) <= availableSpace {
		r.rw.Lock()
		r.buf = append(r.buf, p...)
		r.rw.Unlock()
		return len(p), nil
	} else {
		bumpCount := len(p) - availableSpace
		r.rw.Lock()
		r.buf = append(r.buf[bumpCount:], p...)
		r.rw.Unlock()
		return len(p), nil
	}
}

// Close will render the RingBuffer unusable
func (r *RingBuffer[T]) Close() error {
	if r.isClosed() {
		return ErrRingBufferClosed
	}
	r.rw.Lock()
	defer r.rw.Unlock()

	r.closed = true
	r.buf = make([]T, 0, 1)

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
	r.rw.RLock()
	defer r.rw.RUnlock()
	return len(r.buf)
}

// Cap returns the capacity of the buffer and is immutable
func (r *RingBuffer[T]) Cap() int {
	return r.capacity
}

// Peek allows the user to view the first n items of the ring buffer without affecting the data
func (r *RingBuffer[T]) Peek(n int) ([]T, error) {
	if r.isClosed() {
		return nil, ErrRingBufferClosed
	}

	if n < 0 {
		return nil, errors.New("n must be positive")
	}

	p := make([]T, n)
	copyLen := min(n, r.Len())

	r.rw.RLock()
	copy(p, r.buf[:copyLen])
	r.rw.RUnlock()
	return p, nil
}

// PeekAll is a convenience method equivalent to [*RingBuffer.Peek(*RingBuffer.Len())]
func (r *RingBuffer[T]) PeekAll() ([]T, error) {
	return r.Peek(r.Len())
}

// Clear will empty the buffer without closing it
func (r *RingBuffer[T]) Clear() error {
	if r.isClosed() {
		return ErrRingBufferClosed
	}

	r.rw.Lock()
	clear(r.buf)
	r.rw.Unlock()

	return nil
}
