package slices_test

import (
	"time"

	"github.com/jghiloni/go-commonutils/v3/values"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jghiloni/go-commonutils/v3/slices"
)

type mapTest struct {
	Val    int
	IsEven *bool
}

type uniqTest struct {
	val int
}

func newUniqTest(val int) uniqTest {
	return uniqTest{val}
}

func utID(u uniqTest) int {
	return u.val
}

var _ = Describe("Slice", func() {
	DescribeTable("AnySlice", func(actual []any, expected []any) {
		Expect(actual).To(HaveLen(len(expected)))
		for i, a := range actual {
			Expect(a).To(BeEquivalentTo(expected[i]))
		}
	},
		Entry("full", slices.AnySlice([]string{"a", "b", "c"}), []any{"a", "b", "c"}),
		Entry("empty", slices.AnySlice([]bool{}), []any{}),
	)

	DescribeTable("FromAnySlice", func(actual []string, expected []string) {
		Expect(actual).To(HaveLen(len(expected)))
		for i, a := range actual {
			Expect(a).To(BeEquivalentTo(expected[i]))
		}
	},
		Entry("full", slices.FromAnySlice[string]([]any{"a", "b", "c"}), []string{"a", "b", "c"}),
		Entry("empty", slices.FromAnySlice[string]([]any{}), []string{}),
	)

	DescribeTable("SubsliceUntil", func(actual []bool, expected []bool) {
		Expect(actual).To(Equal(expected))
	},
		Entry("full", slices.SubsliceUntil([]bool{false, false, true}, func(b bool) bool { return b }), []bool{false, false}),
		Entry("empty", slices.SubsliceUntil([]bool{}, func(b bool) bool { return !b }), []bool{}),
	)

	DescribeTable("Map", func(actual []mapTest, expected []mapTest) {
		Expect(actual).To(Equal(expected))
	},
		Entry("full", slices.Map([]int{1, 2, 4, 6}, func(i int) mapTest {
			return mapTest{Val: i, IsEven: values.Ref(i%2 == 0)}
		}), []mapTest{{1, values.Ref(false)}, {2, values.Ref(true)}, {4, values.Ref(true)}, {6, values.Ref(true)}}),
		Entry("empty", slices.Map([]int{}, func(int) mapTest { return mapTest{} }), []mapTest{}),
	)

	DescribeTable("Filter", func(actual []time.Duration, expected []time.Duration) {
		Expect(actual).To(Equal(expected))
	},
		Entry("full", slices.Filter([]time.Duration{time.Minute, time.Millisecond, 3601 * time.Second, time.Hour, 59_000 * time.Millisecond}, func(d time.Duration) bool {
			return d < time.Hour
		}), []time.Duration{time.Minute, time.Millisecond, 59_000_000_000 * time.Nanosecond}),
		Entry("empty result", slices.Filter([]time.Duration{61 * time.Minute, 2 * time.Hour}, func(d time.Duration) bool {
			return d < time.Hour
		}), []time.Duration{}),
		Entry("empty src", slices.Filter([]time.Duration{}, nil), []time.Duration{}),
	)

	DescribeTable("Reduce", func(x int, y int) { Expect(x).To(Equal(y)) },
		Entry("fac", slices.Reduce([]int{1, 2, 3, 4, 5}, func(cur int, val int) int {
			return cur * val
		}, 1), 120),
		Entry("empty factorial", slices.Reduce([]int{}, func(c int, v int) int { return c * v }, 1), 1),
	)

	DescribeTable("Uniq", func(expected, initial []int) {
		actual := slices.Uniq(initial)
		Expect(actual).To(Equal(expected))
	},
		Entry("no dupes", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}),
		Entry("all dupes", []int{1}, []int{1, 1, 1, 1, 1, 1}),
		Entry("some dupes", []int{4, 3, 1, 2}, []int{4, 3, 3, 4, 1, 2, 3, 2, 4, 1}),
	)

	DescribeTable("UniqFunc", func(expected, initial []uniqTest) {
		actual := slices.UniqFunc(initial, utID)
		Expect(actual).To(Equal(expected))
	},
		Entry("no dupes", []uniqTest{newUniqTest(1), newUniqTest(2), newUniqTest(3), newUniqTest(4), newUniqTest(5)}, []uniqTest{newUniqTest(1), newUniqTest(2), newUniqTest(3), newUniqTest(4), newUniqTest(5)}),
		Entry("all dupes", []uniqTest{newUniqTest(1)}, []uniqTest{newUniqTest(1), newUniqTest(1), newUniqTest(1), newUniqTest(1), newUniqTest(1), newUniqTest(1)}),
		Entry("some dupes", []uniqTest{newUniqTest(4), newUniqTest(3), newUniqTest(1), newUniqTest(2)}, []uniqTest{newUniqTest(4), newUniqTest(3), newUniqTest(3), newUniqTest(4), newUniqTest(1), newUniqTest(2), newUniqTest(3), newUniqTest(2), newUniqTest(4), newUniqTest(1)}),
	)

	Describe("Shuffle", func() {
		It("Returns things randomized", func() {
			test := make([]int, 1000)
			for i := range 1000 {
				test[i] = i
			}

			used := make([]int, 0, 1000)
			for i := range slices.Shuffle(test) {
				Expect(slices.Contains(used, i)).To(BeFalse())
				used = append(used, i)
			}
		})
	})
})
