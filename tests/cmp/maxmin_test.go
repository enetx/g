package g_test

import (
	"testing"

	"github.com/enetx/g/cmp"
)

func TestMin(t *testing.T) {
	result1 := cmp.Min(5, 2, 8, 1, 9)
	if result1 != 1 {
		t.Errorf("Min(5, 2, 8, 1, 9) = %d, want 1", result1)
	}

	result2 := cmp.Min(3.5, 1.2, 7.8)
	if result2 != 1.2 {
		t.Errorf("Min(3.5, 1.2, 7.8) = %f, want 1.2", result2)
	}

	result3 := cmp.Min("zebra", "apple", "banana")
	if result3 != "apple" {
		t.Errorf("Min strings = %s, want apple", result3)
	}
}

func TestMax(t *testing.T) {
	result1 := cmp.Max(5, 2, 8, 1, 9)
	if result1 != 9 {
		t.Errorf("Max(5, 2, 8, 1, 9) = %d, want 9", result1)
	}

	result2 := cmp.Max(3.5, 1.2, 7.8)
	if result2 != 7.8 {
		t.Errorf("Max(3.5, 1.2, 7.8) = %f, want 7.8", result2)
	}

	result3 := cmp.Max("zebra", "apple", "banana")
	if result3 != "zebra" {
		t.Errorf("Max strings = %s, want zebra", result3)
	}
}

func TestMinByCustom(t *testing.T) {
	// Test with custom comparison function - reverse order
	reverseCompare := func(x, y int) cmp.Ordering {
		if x < y {
			return cmp.Greater
		} else if x > y {
			return cmp.Less
		}
		return cmp.Equal
	}

	result := cmp.MinBy(reverseCompare, 5, 2, 8, 1, 9)
	if result != 9 {
		t.Errorf("MinBy with reverse compare = %d, want 9", result)
	}
}

func TestMaxByCustom(t *testing.T) {
	// Test with custom comparison function - reverse order
	reverseCompare := func(x, y int) cmp.Ordering {
		if x < y {
			return cmp.Greater
		} else if x > y {
			return cmp.Less
		}
		return cmp.Equal
	}

	result := cmp.MaxBy(reverseCompare, 5, 2, 8, 1, 9)
	if result != 1 {
		t.Errorf("MaxBy with reverse compare = %d, want 1", result)
	}
}

func TestMin_EmptySlice(t *testing.T) {
	result := cmp.Min[int]()
	if result != 0 {
		t.Errorf("Min with no arguments should return zero value, got %d", result)
	}
}

func TestMax_EmptySlice(t *testing.T) {
	result := cmp.Max[int]()
	if result != 0 {
		t.Errorf("Max with no arguments should return zero value, got %d", result)
	}
}

func TestMinBy_EmptySlice(t *testing.T) {
	compareFunc := func(x, y int) cmp.Ordering {
		if x < y {
			return cmp.Less
		} else if x > y {
			return cmp.Greater
		}
		return cmp.Equal
	}

	result := cmp.MinBy(compareFunc)
	if result != 0 {
		t.Errorf("MinBy with no arguments should return zero value, got %d", result)
	}
}

func TestMaxBy_EmptySlice(t *testing.T) {
	compareFunc := func(x, y int) cmp.Ordering {
		if x < y {
			return cmp.Less
		} else if x > y {
			return cmp.Greater
		}
		return cmp.Equal
	}

	result := cmp.MaxBy(compareFunc)
	if result != 0 {
		t.Errorf("MaxBy with no arguments should return zero value, got %d", result)
	}
}

func TestMin_SingleValue(t *testing.T) {
	result := cmp.Min(42)
	if result != 42 {
		t.Errorf("Min with single value should return that value, got %d", result)
	}
}

func TestMax_SingleValue(t *testing.T) {
	result := cmp.Max(42)
	if result != 42 {
		t.Errorf("Max with single value should return that value, got %d", result)
	}
}
