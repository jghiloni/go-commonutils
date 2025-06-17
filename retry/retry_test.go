package retry_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/jghiloni/go-commonutils/v2/retry"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func doWork(ctx context.Context) (int, error) {
	// let's have a 10% chance of failure
	time.Sleep(500 * time.Millisecond)
	select {
	case <-ctx.Done():
		return 0, retry.ErrDoNotRetry
	default:
		fmt.Fprintln(io.Discard, "moving on")
	}

	now := time.Now().Round(time.Second)
	nowSecond := now.Second()

	// check for mod 1 instead of mod 0 to avoid the chance that nowSeconds is 0,
	// which makes it harder to test to see if we're returning our value or the zero value
	if nowSecond%3 == 1 {
		return nowSecond, errors.New("oopsie doodles")
	}

	if nowSecond%19 == 1 {
		return nowSecond, fmt.Errorf("%w: very special case", retry.ErrDoNotRetry)
	}

	// check if nowSecond is zero here too, and return something else just in case
	if nowSecond == 0 {
		nowSecond = -1
	}

	return nowSecond, nil
}

var _ = Describe("Retry", func() {
	It("Works with no backoff", func() {
		ctx := context.Background()
		v, err := retry.DoWithRetry(ctx, doWork, retry.WithBackoffAlgorithm[int](retry.None), retry.WithMaxRetries[int](5))
		if err != nil {
			Expect(errors.Is(err, retry.ErrDoNotRetry) || errors.Is(err, retry.ErrTooManyRetries)).To(BeTrue())
			Expect(v).Should(BeZero())
		}

		Expect(v).ShouldNot(BeZero())
	})

	It("Works with linear backoff", func() {
		ctx := context.Background()
		v, err := retry.DoWithRetry(ctx, doWork, retry.WithBackoffAlgorithm[int](retry.Linear), retry.WithMaxRetries[int](5))
		if err != nil {
			Expect(errors.Is(err, retry.ErrDoNotRetry) || errors.Is(err, retry.ErrTooManyRetries)).To(BeTrue())
			Expect(v).Should(BeZero())
		}

		Expect(v).ShouldNot(BeZero())
	})

	It("times out gracefully", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
		defer cancel()

		v, err := retry.DoWithRetry(ctx, doWork, retry.WithBackoffAlgorithm[int](retry.Linear), retry.WithMaxRetries[int](5))
		Expect(v).To(BeZero())
		Expect(errors.Is(err, retry.ErrDoNotRetry)).To(BeTrue())
	})
})
