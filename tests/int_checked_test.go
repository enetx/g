package g_test

import (
	"math"
	"testing"

	. "github.com/enetx/g"
)

func TestIntCheckedAdd(t *testing.T) {
	tests := []struct {
		name string
		i, b Int
		want Option[Int]
	}{
		{"simple", 2, 3, Some(Int(5))},
		{"negative", -2, -3, Some(Int(-5))},
		{"max plus zero", math.MaxInt, 0, Some(Int(math.MaxInt))},
		{"max plus one overflows", math.MaxInt, 1, None[Int]()},
		{"min plus minus one overflows", math.MinInt, -1, None[Int]()},
		{"min plus max", math.MinInt, math.MaxInt, Some(Int(-1))},
		{"max plus max overflows", math.MaxInt, math.MaxInt, None[Int]()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.CheckedAdd(tt.b); got != tt.want {
				t.Errorf("Int(%d).CheckedAdd(%d) = %v, want %v", tt.i, tt.b, got, tt.want)
			}
		})
	}
}

func TestIntCheckedSub(t *testing.T) {
	tests := []struct {
		name string
		i, b Int
		want Option[Int]
	}{
		{"simple", 5, 3, Some(Int(2))},
		{"min minus zero", math.MinInt, 0, Some(Int(math.MinInt))},
		{"min minus one overflows", math.MinInt, 1, None[Int]()},
		{"max minus minus one overflows", math.MaxInt, -1, None[Int]()},
		{"max minus max", math.MaxInt, math.MaxInt, Some(Int(0))},
		{"zero minus min overflows", 0, math.MinInt, None[Int]()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.CheckedSub(tt.b); got != tt.want {
				t.Errorf("Int(%d).CheckedSub(%d) = %v, want %v", tt.i, tt.b, got, tt.want)
			}
		})
	}
}

func TestIntCheckedMul(t *testing.T) {
	tests := []struct {
		name string
		i, b Int
		want Option[Int]
	}{
		{"simple", 6, 7, Some(Int(42))},
		{"zero left", 0, math.MinInt, Some(Int(0))},
		{"zero right", math.MinInt, 0, Some(Int(0))},
		{"min times one", math.MinInt, 1, Some(Int(math.MinInt))},
		{"min times minus one overflows", math.MinInt, -1, None[Int]()},
		{"minus one times min overflows", -1, math.MinInt, None[Int]()},
		{"minus one times max", -1, math.MaxInt, Some(Int(-math.MaxInt))},
		{"max times two overflows", math.MaxInt, 2, None[Int]()},
		{"half max times two", math.MaxInt / 2, 2, Some(Int(math.MaxInt - 1))},
		{"min half times two", math.MinInt / 2, 2, Some(Int(math.MinInt))},
		{"negative overflow", math.MinInt / 2, -2, None[Int]()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.CheckedMul(tt.b); got != tt.want {
				t.Errorf("Int(%d).CheckedMul(%d) = %v, want %v", tt.i, tt.b, got, tt.want)
			}
		})
	}
}

func TestIntCheckedDiv(t *testing.T) {
	tests := []struct {
		name string
		i, b Int
		want Option[Int]
	}{
		{"simple", 10, 2, Some(Int(5))},
		{"divide by zero", 1, 0, None[Int]()},
		{"min div minus one overflows", math.MinInt, -1, None[Int]()},
		{"min div one", math.MinInt, 1, Some(Int(math.MinInt))},
		{"max div minus one", math.MaxInt, -1, Some(Int(-math.MaxInt))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.CheckedDiv(tt.b); got != tt.want {
				t.Errorf("Int(%d).CheckedDiv(%d) = %v, want %v", tt.i, tt.b, got, tt.want)
			}
		})
	}
}

func TestIntCheckedRem(t *testing.T) {
	tests := []struct {
		name string
		i, b Int
		want Option[Int]
	}{
		{"simple", 10, 3, Some(Int(1))},
		{"rem by zero", 1, 0, None[Int]()},
		{"min rem minus one none", math.MinInt, -1, None[Int]()},
		{"min rem one", math.MinInt, 1, Some(Int(0))},
		{"negative dividend", -7, 3, Some(Int(-1))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.CheckedRem(tt.b); got != tt.want {
				t.Errorf("Int(%d).CheckedRem(%d) = %v, want %v", tt.i, tt.b, got, tt.want)
			}
		})
	}
}

func TestIntCheckedNeg(t *testing.T) {
	tests := []struct {
		name string
		i    Int
		want Option[Int]
	}{
		{"positive", 5, Some(Int(-5))},
		{"negative", -5, Some(Int(5))},
		{"zero", 0, Some(Int(0))},
		{"max", math.MaxInt, Some(Int(-math.MaxInt))},
		{"min overflows", math.MinInt, None[Int]()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.CheckedNeg(); got != tt.want {
				t.Errorf("Int(%d).CheckedNeg() = %v, want %v", tt.i, got, tt.want)
			}
		})
	}
}

func TestIntCheckedAbs(t *testing.T) {
	tests := []struct {
		name string
		i    Int
		want Option[Int]
	}{
		{"positive", 5, Some(Int(5))},
		{"negative", -5, Some(Int(5))},
		{"zero", 0, Some(Int(0))},
		{"min plus one", math.MinInt + 1, Some(Int(math.MaxInt))},
		{"min overflows", math.MinInt, None[Int]()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.CheckedAbs(); got != tt.want {
				t.Errorf("Int(%d).CheckedAbs() = %v, want %v", tt.i, got, tt.want)
			}
		})
	}
}

func TestIntCheckedPow(t *testing.T) {
	tests := []struct {
		name string
		i    Int
		exp  Int
		want Option[Int]
	}{
		{"zero exp", 10, 0, Some(Int(1))},
		{"zero base zero exp", 0, 0, Some(Int(1))},
		{"min base zero exp", math.MinInt, 0, Some(Int(1))},
		{"negative exp", 2, -1, None[Int]()},
		{"simple", 2, 10, Some(Int(1024))},
		{"negative base odd exp", -3, 3, Some(Int(-27))},
		{"negative base even exp", -3, 4, Some(Int(81))},
		{"ten pow eighteen", 10, 18, Some(Int(1000000000000000000))},
		{"ten pow forty overflows", 10, 40, None[Int]()},
		{"two pow sixty two", 2, 62, Some(Int(1) << 62)},
		{"two pow sixty three overflows", 2, 63, None[Int]()},
		{"minus two pow sixty three", -2, 63, Some(Int(math.MinInt))},
		{"minus two pow sixty four overflows", -2, 64, None[Int]()},
		{"one pow huge", 1, 1000, Some(Int(1))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.CheckedPow(tt.exp); got != tt.want {
				t.Errorf("Int(%d).CheckedPow(%d) = %v, want %v", tt.i, tt.exp, got, tt.want)
			}
		})
	}
}

func TestIntSaturatingAdd(t *testing.T) {
	tests := []struct {
		name string
		i, b Int
		want Int
	}{
		{"simple", 2, 3, 5},
		{"saturate max", math.MaxInt, 1, math.MaxInt},
		{"saturate min", math.MinInt, -1, math.MinInt},
		{"no saturation", math.MaxInt, -1, math.MaxInt - 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.SaturatingAdd(tt.b); got != tt.want {
				t.Errorf("Int(%d).SaturatingAdd(%d) = %d, want %d", tt.i, tt.b, got, tt.want)
			}
		})
	}
}

func TestIntSaturatingSub(t *testing.T) {
	tests := []struct {
		name string
		i, b Int
		want Int
	}{
		{"simple", 5, 3, 2},
		{"saturate min", math.MinInt, 1, math.MinInt},
		{"saturate max", math.MaxInt, -1, math.MaxInt},
		{"zero minus min saturates", 0, math.MinInt, math.MaxInt},
		{"no saturation", math.MinInt, -1, math.MinInt + 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.SaturatingSub(tt.b); got != tt.want {
				t.Errorf("Int(%d).SaturatingSub(%d) = %d, want %d", tt.i, tt.b, got, tt.want)
			}
		})
	}
}

func TestIntSaturatingMul(t *testing.T) {
	tests := []struct {
		name string
		i, b Int
		want Int
	}{
		{"simple", 6, 7, 42},
		{"zero", math.MinInt, 0, 0},
		{"saturate max same signs", math.MaxInt, 2, math.MaxInt},
		{"saturate max both negative", math.MinInt, -1, math.MaxInt},
		{"saturate min mixed signs", math.MaxInt, -2, math.MinInt},
		{"saturate min mixed signs swapped", -2, math.MaxInt, math.MinInt},
		{"min times one", math.MinInt, 1, math.MinInt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.SaturatingMul(tt.b); got != tt.want {
				t.Errorf("Int(%d).SaturatingMul(%d) = %d, want %d", tt.i, tt.b, got, tt.want)
			}
		})
	}
}

// Plain Int arithmetic is Go-native and wraps on overflow (two's complement).
func TestIntPlainOpsWrap(t *testing.T) {
	t.Run("add wraps max", func(t *testing.T) {
		if got := Int(math.MaxInt).Add(1); got != math.MinInt {
			t.Errorf("Add = %d, want %d", got, Int(math.MinInt))
		}
	})

	t.Run("sub wraps min", func(t *testing.T) {
		if got := Int(math.MinInt).Sub(1); got != math.MaxInt {
			t.Errorf("Sub = %d, want %d", got, Int(math.MaxInt))
		}
	})

	t.Run("mul wraps min times minus one", func(t *testing.T) {
		if got := Int(math.MinInt).Mul(-1); got != math.MinInt {
			t.Errorf("Mul = %d, want %d", got, Int(math.MinInt))
		}
	})

	t.Run("mul wraps max times two", func(t *testing.T) {
		if got := Int(math.MaxInt).Mul(2); got != -2 {
			t.Errorf("Mul = %d, want -2", got)
		}
	})

	t.Run("neg min wraps to min", func(t *testing.T) {
		if got := Int(math.MinInt).Neg(); got != math.MinInt {
			t.Errorf("Neg = %d, want %d", got, Int(math.MinInt))
		}
	})

	t.Run("neg simple", func(t *testing.T) {
		if got := Int(5).Neg(); got != -5 {
			t.Errorf("Neg = %d, want -5", got)
		}
	})

	t.Run("abs min wraps to min", func(t *testing.T) {
		if got := Int(math.MinInt).Abs(); got != math.MinInt {
			t.Errorf("Abs = %d, want %d", got, Int(math.MinInt))
		}
	})

	t.Run("div min by minus one wraps", func(t *testing.T) {
		if got := Int(math.MinInt).Div(-1); got != math.MinInt {
			t.Errorf("Div = %d, want %d", got, Int(math.MinInt))
		}
	})
}

func TestIntOverflowingAdd(t *testing.T) {
	tests := []struct {
		name     string
		i, b     Int
		want     Int
		overflow bool
	}{
		{"simple", 2, 3, 5, false},
		{"max plus one", math.MaxInt, 1, math.MinInt, true},
		{"min plus minus one", math.MinInt, -1, math.MaxInt, true},
		{"no overflow at max", math.MaxInt, 0, math.MaxInt, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, overflow := tt.i.OverflowingAdd(tt.b)
			if got != tt.want || overflow != tt.overflow {
				t.Errorf(
					"Int(%d).OverflowingAdd(%d) = (%d, %v), want (%d, %v)",
					tt.i,
					tt.b,
					got,
					overflow,
					tt.want,
					tt.overflow,
				)
			}
		})
	}
}

func TestIntOverflowingSub(t *testing.T) {
	tests := []struct {
		name     string
		i, b     Int
		want     Int
		overflow bool
	}{
		{"simple", 5, 3, 2, false},
		{"min minus one", math.MinInt, 1, math.MaxInt, true},
		{"max minus minus one", math.MaxInt, -1, math.MinInt, true},
		{"no overflow at min", math.MinInt, 0, math.MinInt, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, overflow := tt.i.OverflowingSub(tt.b)
			if got != tt.want || overflow != tt.overflow {
				t.Errorf(
					"Int(%d).OverflowingSub(%d) = (%d, %v), want (%d, %v)",
					tt.i,
					tt.b,
					got,
					overflow,
					tt.want,
					tt.overflow,
				)
			}
		})
	}
}

func TestIntOverflowingMul(t *testing.T) {
	tests := []struct {
		name     string
		i, b     Int
		want     Int
		overflow bool
	}{
		{"simple", 6, 7, 42, false},
		{"zero", math.MinInt, 0, 0, false},
		{"min times minus one wraps", math.MinInt, -1, math.MinInt, true},
		{"max times two wraps", math.MaxInt, 2, -2, true},
		{"min times one", math.MinInt, 1, math.MinInt, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, overflow := tt.i.OverflowingMul(tt.b)
			if got != tt.want || overflow != tt.overflow {
				t.Errorf(
					"Int(%d).OverflowingMul(%d) = (%d, %v), want (%d, %v)",
					tt.i,
					tt.b,
					got,
					overflow,
					tt.want,
					tt.overflow,
				)
			}
		})
	}
}

func TestIntClamp(t *testing.T) {
	tests := []struct {
		name     string
		i        Int
		min, max Int
		want     Int
	}{
		{"inside range", 5, 0, 10, 5},
		{"below min", -5, 0, 10, 0},
		{"above max", 15, 0, 10, 10},
		{"equal to min", 0, 0, 10, 0},
		{"equal to max", 10, 0, 10, 10},
		{"min bound extremes", math.MinInt, math.MinInt, math.MaxInt, math.MinInt},
		{"max bound extremes", math.MaxInt, math.MinInt, math.MaxInt, math.MaxInt},
		{"clamp min to zero", math.MinInt, 0, 10, 0},
		{"clamp max to ten", math.MaxInt, 0, 10, 10},
		{"degenerate range", 5, 7, 7, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Clamp(tt.min, tt.max); got != tt.want {
				t.Errorf("Int(%d).Clamp(%d, %d) = %d, want %d", tt.i, tt.min, tt.max, got, tt.want)
			}
		})
	}
}
