package data_test

import (
	"errors"

	"github.com/jghiloni/go-commonutils/v3/data"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ringbuffer", func() {
	It("Never grows beyond capacity", func() {
		buf := data.NewRingBuffer[int](10)

		for i := range 100 {
			n, err := buf.Push(i)
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(Equal(1))
			Expect(buf.Len()).To(BeNumerically("<=", 10))
			Expect(buf.Cap()).To(BeEquivalentTo(10))
		}

		vals, err := buf.Pop(10)
		Expect(err).NotTo(HaveOccurred())
		Expect(vals).To(Equal([]int{90, 91, 92, 93, 94, 95, 96, 97, 98, 99}))
	})

	It("fails on read and write after close", func() {
		buf := data.NewRingBuffer[bool](3)
		_, err := buf.Push(true, false)
		Expect(err).NotTo(HaveOccurred())
		err = buf.Close()
		Expect(err).NotTo(HaveOccurred())
		_, err = buf.Pop(1)
		Expect(errors.Is(err, data.ErrRingBufferClosed)).To(BeTrue())
		err = nil
		_, err = buf.Push(true)
		Expect(errors.Is(err, data.ErrRingBufferClosed)).To(BeTrue())
	})
})
