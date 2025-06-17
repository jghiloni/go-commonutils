package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

type retryer[T any] interface {
	do(context.Context, Retryable[T]) (T, error)
}

type Retryable[T any] func(context.Context) (T, error)

// RetryOption will affect the behavior of the retry
type RetryOption[T any] func(r retryer[T]) error

type algorithmFunc func(context.Context, int, time.Duration, time.Duration) error

type BackoffAlgorithm uint8

const (
	None BackoffAlgorithm = iota
	Linear
	Exponential
)

var (
	ErrUnknownAlgorithm = errors.New("unknown backoff algorithm")
	ErrInvalidDuration  = errors.New("duration must be > 0ns")
	ErrTooManyRetries   = errors.New("too many retries")
	ErrDoNotRetry       = errors.New("do not retry")

	algs = map[BackoffAlgorithm]algorithmFunc{
		None:        noneBackoff,
		Linear:      linearBackoff,
		Exponential: exponentialBackoff,
	}
)

type retryImpl[T any] struct {
	alg           algorithmFunc
	startInterval time.Duration
	maxDuration   time.Duration
	maxRetries    int
}

func WithBackoffAlgorithm[T any](alg BackoffAlgorithm) RetryOption[T] {
	return func(r retryer[T]) error {
		algFn, ok := algs[alg]
		if !ok {
			return ErrUnknownAlgorithm
		}

		ri := r.(*retryImpl[T])

		ri.alg = algFn
		return nil
	}
}

func WithStartInterval[T any, R Retryable[T]](interval time.Duration) RetryOption[T] {
	return func(r retryer[T]) error {
		if interval <= 0 {
			return ErrInvalidDuration
		}

		ri := r.(*retryImpl[T])

		ri.startInterval = interval
		return nil
	}
}

func WithMaxInterval[T any, R Retryable[T]](interval time.Duration) RetryOption[T] {
	return func(r retryer[T]) error {
		if interval <= 0 {
			return ErrInvalidDuration
		}

		ri := r.(*retryImpl[T])
		ri.maxDuration = interval
		return nil
	}
}

func WithMaxRetries[T any, R Retryable[T]](maxRetries int) RetryOption[T] {
	return func(r retryer[T]) error {
		m := maxRetries
		if m < 1 {
			m = math.MaxInt
		}

		ri := r.(*retryImpl[T])
		ri.maxRetries = m
		return nil
	}
}

func noneBackoff(ctx context.Context, attempt int, _ time.Duration, _ time.Duration) error {
	select {
	case <-ctx.Done():
		return ErrDoNotRetry
	default:
		return nil
	}
}

func linearBackoff(ctx context.Context, attempt int, interval time.Duration, maxDuration time.Duration) error {
	delay := min(time.Duration(attempt+1)*interval, maxDuration)
	select {
	case <-ctx.Done():
		return ErrDoNotRetry
	case <-time.After(delay):
		return nil
	}
}

func exponentialBackoff(ctx context.Context, attempt int, interval time.Duration, maxDuration time.Duration) error {
	var delay time.Duration

	// this will avoid overflow
	if attempt > 63 {
		delay = maxDuration
	}

	if delay == 0 {
		requestedDelay := uint64(interval) * uint64(1<<uint64(attempt))
		if requestedDelay > math.MaxInt64 {
			requestedDelay = math.MaxInt64
		}

		delay = min(time.Duration(requestedDelay), maxDuration)
	}

	select {
	case <-ctx.Done():
		return ErrDoNotRetry
	case <-time.After(delay):
		return nil
	}
}

func (r *retryImpl[T]) do(ctx context.Context, uow Retryable[T]) (val T, err error) {
	var zeroT T

	for i := range r.maxRetries {
		retVal, err := uow(ctx)
		if err == nil {
			return retVal, nil
		}

		if errors.Is(err, ErrDoNotRetry) {
			return zeroT, err
		}

		if err := r.alg(ctx, i, r.startInterval, r.maxDuration); errors.Is(err, ErrDoNotRetry) {
			return zeroT, err
		}
	}

	return zeroT, ErrTooManyRetries
}

func DoWithRetry[T any](ctx context.Context, unitOfWork Retryable[T], options ...RetryOption[T]) (T, error) {
	r := &retryImpl[T]{
		maxDuration:   time.Minute,
		maxRetries:    10,
		startInterval: 250 * time.Millisecond,
		alg:           linearBackoff,
	}

	var zeroT T

	for _, opt := range options {
		if err := opt(r); err != nil {
			return zeroT, err
		}
	}

	return r.do(ctx, unitOfWork)
}
