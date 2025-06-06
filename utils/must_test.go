package utils_test

import (
	"errors"

	"github.com/jghiloni/go-commonutils/v2/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func testFunc(i int) (bool, error) {
	if i%2 == 0 {
		return true, nil
	}

	return false, errors.New("odd")
}

var _ = Describe("Must", func() {
	It("doesn't panic when error is nil", func() {
		Expect(utils.Must(testFunc(2))).To(BeTrue())
	})

	It("panics when error is not nil", func() {
		actual := func() {
			utils.Must(testFunc(1))
		}

		Expect(actual).To(Panic())
	})
})
