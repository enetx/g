package g_test

import (
	"testing"

	"github.com/enetx/g/cmp"
)

func TestOrderingThen(t *testing.T) {
	tests := []struct {
		name     string
		o        cmp.Ordering
		other    cmp.Ordering
		expected cmp.Ordering
	}{
		{"Non-zero receiver", cmp.Ordering(2), cmp.Ordering(3), cmp.Ordering(2)},
		{"Zero receiver", cmp.Ordering(0), cmp.Ordering(3), cmp.Ordering(3)},
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

func TestOrderingCmp(t *testing.T) {
	tests := []struct {
		name     string
		x        cmp.Ordering
		y        cmp.Ordering
		expected cmp.Ordering
	}{
		{"x < y", cmp.Ordering(2), cmp.Ordering(3), cmp.Ordering(-1)},
		{"x = y", cmp.Ordering(2), cmp.Ordering(2), cmp.Ordering(0)},
		{"x > y", cmp.Ordering(3), cmp.Ordering(2), cmp.Ordering(1)},
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

func TestOrderingReverse(t *testing.T) {
	tests := []struct {
		name     string
		input    cmp.Ordering
		expected cmp.Ordering
	}{
		{"Reverse_Less", cmp.Less, cmp.Greater},
		{"Reverse_Equal", cmp.Equal, cmp.Equal},
		{"Reverse_Greater", cmp.Greater, cmp.Less},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Reverse()
			if result != tt.expected {
				t.Errorf("Expected reverse of %v to be %v, but got %v", tt.input, tt.expected, result)
			}
		})
	}
}

func TestOrderingString(t *testing.T) {
	tests := []struct {
		name     string
		input    cmp.Ordering
		expected string
	}{
		{"Less", cmp.Less, "Less"},
		{"Equal", cmp.Equal, "Equal"},
		{"Greater", cmp.Greater, "Greater"},
		{"Unknown", 10, "Unknown Ordering value: 10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			if result != tt.expected {
				t.Errorf("Expected %s for Ordering value %d, but got %s", tt.expected, tt.input, result)
			}
		})
	}
}
