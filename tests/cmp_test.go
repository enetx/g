package g_test

import (
	"testing"

	"github.com/enetx/g/cmp"
)

func TestOrderedThen(t *testing.T) {
	tests := []struct {
		name     string
		o        cmp.Ordered
		other    cmp.Ordered
		expected cmp.Ordered
	}{
		{"Non-zero receiver", cmp.Ordered(2), cmp.Ordered(3), cmp.Ordered(2)},
		{"Zero receiver", cmp.Ordered(0), cmp.Ordered(3), cmp.Ordered(3)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.o.Then(tt.other)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}

func TestCmp(t *testing.T) {
	tests := []struct {
		name     string
		x        cmp.Ordered
		y        cmp.Ordered
		expected cmp.Ordered
	}{
		{"x < y", cmp.Ordered(2), cmp.Ordered(3), cmp.Ordered(-1)},
		{"x = y", cmp.Ordered(2), cmp.Ordered(2), cmp.Ordered(0)},
		{"x > y", cmp.Ordered(3), cmp.Ordered(2), cmp.Ordered(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cmp.Cmp(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
