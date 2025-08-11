package g_test

import (
	"testing"

	"github.com/enetx/g/cmp"
)

func TestOrdering_Then(t *testing.T) {
	tests := []struct {
		name     string
		receiver cmp.Ordering
		other    cmp.Ordering
		want     cmp.Ordering
	}{
		{"Equal then Less", cmp.Equal, cmp.Less, cmp.Less},
		{"Equal then Greater", cmp.Equal, cmp.Greater, cmp.Greater},
		{"Equal then Equal", cmp.Equal, cmp.Equal, cmp.Equal},
		{"Less then anything", cmp.Less, cmp.Greater, cmp.Less},
		{"Greater then anything", cmp.Greater, cmp.Less, cmp.Greater},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.receiver.Then(tt.other)
			if got != tt.want {
				t.Errorf("Ordering.Then() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrdering_Reverse(t *testing.T) {
	tests := []struct {
		name     string
		ordering cmp.Ordering
		want     cmp.Ordering
	}{
		{"Less reversed", cmp.Less, cmp.Greater},
		{"Greater reversed", cmp.Greater, cmp.Less},
		{"Equal reversed", cmp.Equal, cmp.Equal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ordering.Reverse()
			if got != tt.want {
				t.Errorf("Ordering.Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrdering_IsLt(t *testing.T) {
	tests := []struct {
		name     string
		ordering cmp.Ordering
		want     bool
	}{
		{"Less is Lt", cmp.Less, true},
		{"Equal is not Lt", cmp.Equal, false},
		{"Greater is not Lt", cmp.Greater, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ordering.IsLt()
			if got != tt.want {
				t.Errorf("Ordering.IsLt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrdering_IsEq(t *testing.T) {
	tests := []struct {
		name     string
		ordering cmp.Ordering
		want     bool
	}{
		{"Less is not Eq", cmp.Less, false},
		{"Equal is Eq", cmp.Equal, true},
		{"Greater is not Eq", cmp.Greater, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ordering.IsEq()
			if got != tt.want {
				t.Errorf("Ordering.IsEq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrdering_IsGt(t *testing.T) {
	tests := []struct {
		name     string
		ordering cmp.Ordering
		want     bool
	}{
		{"Less is not Gt", cmp.Less, false},
		{"Equal is not Gt", cmp.Equal, false},
		{"Greater is Gt", cmp.Greater, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ordering.IsGt()
			if got != tt.want {
				t.Errorf("Ordering.IsGt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrdering_String(t *testing.T) {
	tests := []struct {
		name     string
		ordering cmp.Ordering
		want     string
	}{
		{"Less string", cmp.Less, "Less"},
		{"Equal string", cmp.Equal, "Equal"},
		{"Greater string", cmp.Greater, "Greater"},
		{"Unknown ordering", cmp.Ordering(99), "Unknown Ordering value: 99"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ordering.String()
			if got != tt.want {
				t.Errorf("Ordering.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
