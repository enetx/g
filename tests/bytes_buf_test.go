package g

import (
	"testing"

	. "github.com/enetx/g"
)

func TestStringBuf(t *testing.T) {
	tests := []struct {
		value    Int
		expected String
	}{
		{0, "0"},
		{5, "5"},
		{9, "9"},
		{10, "10"},
		{42, "42"},
		{123456789, "123456789"},
		{-1, "-1"},
		{-42, "-42"},
	}

	for _, tt := range tests {
		{
			buf := NewBytes(0, 20)
			got := tt.value.StringBuf(&buf)
			if got != tt.expected {
				t.Errorf("StringBuf(%d) = %q, want %q", tt.value, got, tt.expected)
			}
		}

		{
			buf := NewBytes()
			got := tt.value.StringBuf(&buf)
			if got != tt.expected {
				t.Errorf("StringBufSafe(%d) = %q, want %q", tt.value, got, tt.expected)
			}
		}
	}
}
