package g_test

import (
	"testing"

	"gitlab.com/x0xO/g"
)

func TestIntIsPositive(t *testing.T) {
	tests := []struct {
		name string
		i    g.Int
		want bool
	}{
		{"positive", 1, true},
		{"negative", -1, false},
		{"zero", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.IsPositive(); got != tt.want {
				t.Errorf("Int.IsPositive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntNegative(t *testing.T) {
	tests := []struct {
		name string
		i    g.Int
		want bool
	}{
		{"positive", 1, false},
		{"negative", -1, true},
		{"zero", 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.IsNegative(); got != tt.want {
				t.Errorf("Int.IsNegative() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntRandomRange(t *testing.T) {
	for range 100 {
		min := g.NewInt(100).Random()
		max := g.NewInt(100).Random().Add(min)

		r := min.RandomRange(max)
		if r.Lt(min) || r.Gt(max) {
			t.Errorf("RandomRange(%d, %d) = %d, want in range [%d, %d]", min, max, r, min, max)
		}
	}
}

func TestIntMax(t *testing.T) {
	if max := g.NewInt(1).Max(1, 2, 3, 4, 5); max != 5 {
		t.Errorf("Max() = %d, want: %d.", max, 5)
	}
}

func TestIntMin(t *testing.T) {
	if min := g.NewInt(1).Min(2, 3, 4, 5); min != 1 {
		t.Errorf("Min() = %d; want: %d", min, 1)
	}
}
