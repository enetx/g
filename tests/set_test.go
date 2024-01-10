package g_test

import (
	"testing"

	"gitlab.com/x0xO/g"
)

func TestSetDifference(t *testing.T) {
	set1 := g.SetOf(1, 2, 3, 4)
	set2 := g.SetOf(3, 4, 5, 6)
	set5 := g.SetOf(1, 2)
	set6 := g.SetOf(2, 3, 4)

	set3 := set1.Difference(set2).Collect()
	set4 := set2.Difference(set1).Collect()
	set7 := set5.Difference(set6).Collect()
	set8 := set6.Difference(set5).Collect()

	if set3.Len() != 2 || set3.Ne(g.SetOf(1, 2)) {
		t.Errorf("Unexpected result: %v", set3)
	}

	if set4.Len() != 2 || set4.Ne(g.SetOf(5, 6)) {
		t.Errorf("Unexpected result: %v", set4)
	}

	if set7.Len() != 1 || set7.Ne(g.SetOf(1)) {
		t.Errorf("Unexpected result: %v", set7)
	}

	if set8.Len() != 2 || set8.Ne(g.SetOf(3, 4)) {
		t.Errorf("Unexpected result: %v", set8)
	}
}

func TestSetSymmetricDifference(t *testing.T) {
	set1 := g.NewSet[int](10)
	set2 := set1.Clone()
	result := set1.SymmetricDifference(set2).Collect()

	if !result.Empty() {
		t.Errorf("SymmetricDifference between equal sets should be empty, got %v", result)
	}

	set1 = g.SetOf(0, 1, 2, 3, 4)
	set2 = g.SetOf(5, 6, 7, 8, 9)
	result = set1.SymmetricDifference(set2).Collect()
	expected := set1.Union(set2).Collect()

	if !result.Eq(expected) {
		t.Errorf(
			"SymmetricDifference between disjoint sets should be their union, expected %v but got %v",
			expected,
			result,
		)
	}

	set1 = g.SetOf(0, 1, 2, 3, 4, 5)
	set2 = g.SetOf(4, 5, 6, 7, 8)
	result = set1.SymmetricDifference(set2).Collect()
	expected2 := g.SetOf(0, 1, 2, 3, 6, 7, 8)

	if !result.Eq(expected2) {
		t.Errorf(
			"SymmetricDifference between sets with common elements should be correct, expected %v but got %v",
			expected,
			result,
		)
	}
}

func TestSetIntersection(t *testing.T) {
	set1 := g.Set[int]{}
	set2 := g.Set[int]{}

	set1 = set1.Add(1, 2, 3)
	set2 = set2.Add(2, 3, 4)

	set3 := set1.Intersection(set2).Collect()

	if !set3.Contains(2) || !set3.Contains(3) {
		t.Error("Intersection failed")
	}
}

func TestSetUnion(t *testing.T) {
	set1 := g.NewSet[int]().Add(1, 2, 3)
	set2 := g.NewSet[int]().Add(2, 3, 4)
	set3 := g.NewSet[int]().Add(1, 2, 3, 4)

	result := set1.Union(set2).Collect()

	if result.Len() != 4 {
		t.Errorf("Union(%v, %v) returned %v; expected %v", set1, set2, result, set3)
	}

	for v := range set3 {
		if !result.Contains(v) {
			t.Errorf("Union(%v, %v) missing element %v", set1, set2, v)
		}
	}
}

func TestSetSubset(t *testing.T) {
	tests := []struct {
		name  string
		s     g.Set[int]
		other g.Set[int]
		want  bool
	}{
		{
			name:  "test_subset_1",
			s:     g.SetOf(1, 2, 3),
			other: g.SetOf(1, 2, 3, 4, 5),
			want:  true,
		},
		{
			name:  "test_subset_2",
			s:     g.SetOf(1, 2, 3, 4),
			other: g.SetOf(1, 2, 3),
			want:  false,
		},
		{
			name:  "test_subset_3",
			s:     g.SetOf(5, 4, 3, 2, 1),
			other: g.SetOf(1, 2, 3, 4, 5),
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Subset(tt.other); got != tt.want {
				t.Errorf("Subset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetSuperset(t *testing.T) {
	tests := []struct {
		name  string
		s     g.Set[int]
		other g.Set[int]
		want  bool
	}{
		{
			name:  "test_superset_1",
			s:     g.SetOf(1, 2, 3, 4, 5),
			other: g.SetOf(1, 2, 3),
			want:  true,
		},
		{
			name:  "test_superset_2",
			s:     g.SetOf(1, 2, 3),
			other: g.SetOf(1, 2, 3, 4),
			want:  false,
		},
		{
			name:  "test_superset_3",
			s:     g.SetOf(1, 2, 3, 4, 5),
			other: g.SetOf(5, 4, 3, 2, 1),
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Superset(tt.other); got != tt.want {
				t.Errorf("Superset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetEq(t *testing.T) {
	tests := []struct {
		name  string
		s     g.Set[int]
		other g.Set[int]
		want  bool
	}{
		{
			name:  "test_eq_1",
			s:     g.SetOf(1, 2, 3),
			other: g.SetOf(1, 2, 3),
			want:  true,
		},
		{
			name:  "test_eq_2",
			s:     g.SetOf(1, 2, 3),
			other: g.SetOf(1, 2, 4),
			want:  false,
		},
		{
			name:  "test_eq_3",
			s:     g.SetOf(1, 2, 3),
			other: g.SetOf(3, 2, 1),
			want:  true,
		},
		{
			name:  "test_eq_4",
			s:     g.SetOf(1, 2, 3, 4),
			other: g.SetOf(1, 2, 3),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Eq(tt.other); got != tt.want {
				t.Errorf("Eq() = %v, want %v", got, tt.want)
			}
		})
	}
}
