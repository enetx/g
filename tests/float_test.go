package g_test

import (
	"testing"

	"gitlab.com/x0xO/g"
)

func TestFloatCompare(t *testing.T) {
	testCases := []struct {
		f1       g.Float
		f2       g.Float
		expected g.Int
	}{
		{3.14, 6.28, -1},
		{6.28, 3.14, 1},
		{1.23, 1.23, 0},
		{-2.5, 2.5, -1},
	}

	for _, tc := range testCases {
		result := tc.f1.Compare(tc.f2)
		if !result.Eq(tc.expected) {
			t.Errorf("Compare(%f, %f): expected %d, got %d", tc.f1, tc.f2, tc.expected, result)
		}
	}
}

func TestFloatEq(t *testing.T) {
	testCases := []struct {
		f1       g.Float
		f2       g.Float
		expected bool
	}{
		{3.14, 6.28, false},
		{1.23, 1.23, true},
		{0.0, 0.0, true},
		{-2.5, 2.5, false},
	}

	for _, tc := range testCases {
		result := tc.f1.Eq(tc.f2)
		if result != tc.expected {
			t.Errorf("Eq(%f, %f): expected %t, got %t", tc.f1, tc.f2, tc.expected, result)
		}
	}
}

func TestFloatNe(t *testing.T) {
	testCases := []struct {
		f1       g.Float
		f2       g.Float
		expected bool
	}{
		{3.14, 6.28, true},
		{1.23, 1.23, false},
		{0.0, 0.0, false},
		{-2.5, 2.5, true},
	}

	for _, tc := range testCases {
		result := tc.f1.Ne(tc.f2)
		if result != tc.expected {
			t.Errorf("Ne(%f, %f): expected %t, got %t", tc.f1, tc.f2, tc.expected, result)
		}
	}
}

func TestFloatGt(t *testing.T) {
	testCases := []struct {
		f1       g.Float
		f2       g.Float
		expected bool
	}{
		{3.14, 6.28, false},
		{6.28, 3.14, true},
		{1.23, 1.23, false},
		{-2.5, 2.5, false},
	}

	for _, tc := range testCases {
		result := tc.f1.Gt(tc.f2)
		if result != tc.expected {
			t.Errorf("Gt(%f, %f): expected %t, got %t", tc.f1, tc.f2, tc.expected, result)
		}
	}
}

func TestFloatLt(t *testing.T) {
	testCases := []struct {
		f1       g.Float
		f2       g.Float
		expected bool
	}{
		{3.14, 6.28, true},
		{6.28, 3.14, false},
		{1.23, 1.23, false},
		{-2.5, 2.5, true},
	}
	for _, tc := range testCases {
		result := tc.f1.Lt(tc.f2)
		if result != tc.expected {
			t.Errorf("Lt(%f, %f): expected %t, got %t", tc.f1, tc.f2, tc.expected, result)
		}
	}
}

func TestFloatRoundDecimal(t *testing.T) {
	testCases := []struct {
		value    g.Float
		decimals int
		expected g.Float
	}{
		{3.1415926535, 2, 3.14},
		{3.1415926535, 3, 3.142},
		{100.123456789, 4, 100.1235},
		{-5.6789, 1, -5.7},
		{12345.6789, 0, 12346},
	}

	for _, testCase := range testCases {
		result := testCase.value.RoundDecimal(testCase.decimals)
		if result != testCase.expected {
			t.Errorf(
				"Failed: value=%.10f decimals=%d, expected=%.10f, got=%.10f\n",
				testCase.value,
				testCase.decimals,
				testCase.expected,
				result,
			)
		}
	}
}

func TestFloatMax(t *testing.T) {
	if max := g.NewFloat(2.2).Max(2.8, 2.1, 2.7); max != 2.8 {
		t.Errorf("Max() = %f, want: %f.", max, 2.8)
	}
}

func TestFloatMin(t *testing.T) {
	if min := g.NewFloat(2.2).Min(2.8, 2.1, 2.7); min != 2.1 {
		t.Errorf("Min() = %f; want: %f", min, 2.1)
	}
}
