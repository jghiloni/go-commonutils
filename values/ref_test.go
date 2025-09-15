package values_test

import (
	"github.com/jghiloni/go-commonutils/v3/values"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type refTestStruct struct {
	Foo string
	Bar string
}

var (
	refTestString     *any = values.Ref(any("foo"))
	refTestInt        *any = values.Ref(any(3))
	refTestBool       *any = values.Ref(any(true))
	refTestStructInst *any = values.Ref(any(refTestStruct{"one", "two"}))
)

var _ = Describe("Pointer Utilities", func() {
	DescribeTable("Ref", func(actual *any, expected any, negativeAssign bool) {
		Expect(actual).NotTo(BeNil())

		if negativeAssign {
			Expect(*actual).NotTo(BeAssignableToTypeOf(expected))
		} else {
			Expect(*actual).To(BeAssignableToTypeOf(expected))
			Expect(*actual).To(Equal(expected))
		}
	},
		Entry("string", refTestString, any("foo"), false),
		Entry("int", refTestInt, any(3), false),
		Entry("bool", refTestBool, any(true), false),
		Entry("invalid", refTestInt, any(4.2), true),
		Entry("struct", refTestStructInst, any(refTestStruct{"one", "two"}), false),
	)

	Describe("NilRefIfZero", func() {
		DescribeTable("Nil on zeros", func(actual *any, expected any) {
			Expect(actual).To(BeNil())
			Expect(expected).To(BeZero())
		},
			Entry(nil, nil, any("")),
			Entry(nil, nil, any(0)),
			Entry(nil, nil, any(false)),
			Entry(nil, nil, any(refTestStruct{})),
			Entry(nil, nil, any([]any(nil))),
		)

		DescribeTable("Non-nil on non-zeros", func(actual *any, expected any, negativeAssign bool) {
			Expect(actual).NotTo(BeNil())

			if negativeAssign {
				Expect(*actual).NotTo(BeAssignableToTypeOf(expected))
			} else {
				Expect(*actual).To(BeAssignableToTypeOf(expected))
				Expect(*actual).To(Equal(expected))
			}
		},
			Entry("string", refTestString, any("foo"), false),
			Entry("int", refTestInt, any(3), false),
			Entry("bool", refTestBool, any(true), false),
			Entry("invalid", refTestInt, any(4.2), true),
			Entry("struct", refTestStructInst, any(refTestStruct{"one", "two"}), false),
		)
	})
})
