package g_test

import (
	"strings"
	"testing"

	. "github.com/enetx/g/iter"
)

func TestCmp(t *testing.T) {
	// Test equal sequences
	cmp1 := Cmp(FromSlice([]int{1, 2, 3}), FromSlice([]int{1, 2, 3}), func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if cmp1 != 0 {
		t.Errorf("Cmp(equal) = %d, want 0", cmp1)
	}

	// Test first less than second
	cmp2 := Cmp(FromSlice([]int{1, 2}), FromSlice([]int{1, 3}), func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if cmp2 != -1 {
		t.Errorf("Cmp(first < second) = %d, want -1", cmp2)
	}

	// Test first greater than second
	cmp3 := Cmp(FromSlice([]int{1, 3}), FromSlice([]int{1, 2}), func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if cmp3 != 1 {
		t.Errorf("Cmp(first > second) = %d, want 1", cmp3)
	}
}

func TestEqual(t *testing.T) {
	// Test equal sequences
	result := Equal(FromSlice([]int{1, 2, 3}), FromSlice([]int{1, 2, 3}))
	if !result {
		t.Errorf("Equal() = %v, want true", result)
	}

	// Test unequal sequences
	result2 := Equal(FromSlice([]int{1, 2, 3}), FromSlice([]int{1, 2, 4}))
	if result2 {
		t.Errorf("Equal(unequal) = %v, want false", result2)
	}

	// Test different lengths
	result3 := Equal(FromSlice([]int{1, 2}), FromSlice([]int{1, 2, 3}))
	if result3 {
		t.Errorf("Equal(different lengths) = %v, want false", result3)
	}

	// Test empty sequences
	result4 := Equal(FromSlice([]int{}), FromSlice([]int{}))
	if !result4 {
		t.Errorf("Equal(empty, empty) = %v, want true", result4)
	}

	// Test one empty, one not
	result5 := Equal(FromSlice([]int{}), FromSlice([]int{1}))
	if result5 {
		t.Errorf("Equal(empty, non-empty) = %v, want false", result5)
	}
}

func TestLt(t *testing.T) {
	// Test first less than second
	result := Lt(FromSlice([]int{1, 2}), FromSlice([]int{1, 3}), func(a, b int) bool {
		return a < b
	})
	if !result {
		t.Errorf("Lt() = %v, want true", result)
	}

	// Test first not less than second
	result2 := Lt(FromSlice([]int{1, 3}), FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if result2 {
		t.Errorf("Lt(not less) = %v, want false", result2)
	}
}

func TestLe(t *testing.T) {
	// Test first less than or equal to second
	result := Le(FromSlice([]int{1, 2}), FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result {
		t.Errorf("Le() = %v, want true", result)
	}
}

func TestGt(t *testing.T) {
	// Test first greater than second
	result := Gt(FromSlice([]int{1, 3}), FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result {
		t.Errorf("Gt() = %v, want true", result)
	}

	// Test first not greater than second
	result2 := Gt(FromSlice([]int{1, 2}), FromSlice([]int{1, 3}), func(a, b int) bool {
		return a < b
	})
	if result2 {
		t.Errorf("Gt(not greater) = %v, want false", result2)
	}
}

func TestGe(t *testing.T) {
	// Test first greater than or equal to second
	result := Ge(FromSlice([]int{1, 2}), FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result {
		t.Errorf("Ge() = %v, want true", result)
	}
}

func TestCmpDifferentLengths(t *testing.T) {
	// Test comparison with different lengths
	cmp1 := Cmp(FromSlice([]int{1, 2}), FromSlice([]int{1, 2, 3}), func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if cmp1 != -1 {
		t.Errorf("Cmp(different lengths, first shorter) = %d, want -1", cmp1)
	}
}

func TestCmpSecondLonger(t *testing.T) {
	// Test comparison when second is longer
	cmp1 := Cmp(FromSlice([]int{1, 2, 3}), FromSlice([]int{1, 2}), func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})
	if cmp1 != 1 {
		t.Errorf("Cmp(second shorter) = %d, want 1", cmp1)
	}
}

func TestLtWithEqualElements(t *testing.T) {
	// Test Lt with equal elements
	result := Lt(FromSlice([]int{1, 2}), FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if result {
		t.Errorf("Lt(equal) = %v, want false", result)
	}
}

func TestLeWithEqualElements(t *testing.T) {
	// Test Le with equal elements
	result := Le(FromSlice([]int{1, 2}), FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result {
		t.Errorf("Le(equal) = %v, want true", result)
	}
}

func TestGtWithEqualElements(t *testing.T) {
	// Test Gt with equal elements
	result := Gt(FromSlice([]int{1, 2}), FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if result {
		t.Errorf("Gt(equal) = %v, want false", result)
	}
}

func TestGeWithEqualElements(t *testing.T) {
	// Test Ge with equal elements
	result := Ge(FromSlice([]int{1, 2}), FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result {
		t.Errorf("Ge(equal) = %v, want true", result)
	}
}

func TestEqualBy(t *testing.T) {
	// Test EqualBy with custom comparison by length
	result := EqualBy(FromSlice([]string{"a", "bb"}), FromSlice([]string{"x", "yy"}), func(a, b string) bool {
		return len(a) == len(b)
	})
	if !result {
		t.Errorf("EqualBy(same lengths) = %v, want true", result)
	}

	// Test EqualBy with different lengths
	result2 := EqualBy(FromSlice([]string{"a", "bb"}), FromSlice([]string{"x", "yyy"}), func(a, b string) bool {
		return len(a) == len(b)
	})
	if result2 {
		t.Errorf("EqualBy(different lengths) = %v, want false", result2)
	}

	// Test EqualBy early termination
	result3 := EqualBy(
		FromSlice([]string{"a", "bb", "ccc"}),
		FromSlice([]string{"x", "yyy", "dddd"}),
		func(a, b string) bool {
			return len(a) == len(b)
		},
	)
	if result3 {
		t.Errorf("EqualBy(early mismatch) = %v, want false", result3)
	}

	// Test EqualBy different sequence lengths
	result4 := EqualBy(FromSlice([]string{"a", "bb"}), FromSlice([]string{"x", "yy", "zzz"}), func(a, b string) bool {
		return len(a) == len(b)
	})
	if result4 {
		t.Errorf("EqualBy(different seq lengths) = %v, want false", result4)
	}

	// Test EqualBy empty sequences
	result5 := EqualBy(FromSlice([]string{}), FromSlice([]string{}), func(a, b string) bool {
		return len(a) == len(b)
	})
	if !result5 {
		t.Errorf("EqualBy(empty) = %v, want true", result5)
	}

	// Test EqualBy case-insensitive comparison
	result6 := EqualBy(
		FromSlice([]string{"Hello", "WORLD"}),
		FromSlice([]string{"hello", "world"}),
		func(a, b string) bool {
			return len(a) == len(b) && strings.EqualFold(a, b)
		},
	)
	if !result6 {
		t.Errorf("EqualBy(case insensitive) = %v, want true", result6)
	}

	// Test EqualBy with early false return
	count := 0
	result7 := EqualBy(
		FromSlice([]string{"a", "bb", "ccc"}),
		FromSlice([]string{"x", "yy", "ddd"}),
		func(a, b string) bool {
			count++
			return len(a) == len(b)
		},
	)
	if !result7 {
		t.Errorf("EqualBy(all match) = %v, want true", result7)
	}
	if count != 3 {
		t.Errorf("EqualBy comparison count = %v, want 3", count)
	}
}

func TestLeAdvanced(t *testing.T) {
	// Test Le with different length sequences (first shorter)
	result := Le(FromSlice([]int{1, 2}), FromSlice([]int{1, 2, 3}), func(a, b int) bool {
		return a < b
	})
	if !result {
		t.Errorf("Le(first shorter) = %v, want true", result)
	}

	// Test Le with different length sequences (second shorter)
	result2 := Le(FromSlice([]int{1, 2, 3}), FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if result2 {
		t.Errorf("Le(second shorter) = %v, want false", result2)
	}

	// Test Le with first less than second
	result3 := Le(FromSlice([]int{1, 1}), FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result3 {
		t.Errorf("Le(first < second) = %v, want true", result3)
	}

	// Test Le with empty sequences
	result4 := Le(FromSlice([]int{}), FromSlice([]int{}), func(a, b int) bool {
		return a < b
	})
	if !result4 {
		t.Errorf("Le(empty sequences) = %v, want true", result4)
	}
}

func TestGeAdvanced(t *testing.T) {
	// Test Ge with different length sequences (first longer)
	result := Ge(FromSlice([]int{1, 2, 3}), FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result {
		t.Errorf("Ge(first longer) = %v, want true", result)
	}

	// Test Ge with different length sequences (second longer)
	result2 := Ge(FromSlice([]int{1, 2}), FromSlice([]int{1, 2, 3}), func(a, b int) bool {
		return a < b
	})
	if result2 {
		t.Errorf("Ge(second longer) = %v, want false", result2)
	}

	// Test Ge with first greater than second
	result3 := Ge(FromSlice([]int{1, 3}), FromSlice([]int{1, 2}), func(a, b int) bool {
		return a < b
	})
	if !result3 {
		t.Errorf("Ge(first > second) = %v, want true", result3)
	}

	// Test Ge with empty sequences
	result4 := Ge(FromSlice([]int{}), FromSlice([]int{}), func(a, b int) bool {
		return a < b
	})
	if !result4 {
		t.Errorf("Ge(empty sequences) = %v, want true", result4)
	}
}
