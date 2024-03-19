package g_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/enetx/g"
)

func TestFloatBytes(t *testing.T) {
	// Test case for positive float
	f := g.Float(3.14)
	expected := []byte{64, 9, 30, 184, 81, 235, 133, 31} // Bytes representation of 3.14 in big-endian
	actual := f.Bytes()
	if actual.Ne(expected) {
		t.Errorf("Bytes representation of positive float incorrect. Expected: %v, Got: %v", expected, actual)
	}

	// Test case for negative float
	f = g.Float(-3.14)
	expected = []byte{192, 9, 30, 184, 81, 235, 133, 31} // Bytes representation of -3.14 in big-endian
	actual = f.Bytes()
	if actual.Ne(expected) {
		t.Errorf("Bytes representation of negative float incorrect. Expected: %v, Got: %v", expected, actual)
	}

	// Test case for infinity
	f = g.Float(math.Inf(1))
	expected = []byte{127, 240, 0, 0, 0, 0, 0, 0} // Bytes representation of positive infinity in big-endian
	actual = f.Bytes()
	if actual.Ne(expected) {
		t.Errorf("Bytes representation of positive infinity incorrect. Expected: %v, Got: %v", expected, actual)
	}

	// Test case for negative infinity
	f = g.Float(math.Inf(-1))
	expected = []byte{255, 240, 0, 0, 0, 0, 0, 0} // Bytes representation of negative infinity in big-endian
	actual = f.Bytes()
	if actual.Ne(expected) {
		t.Errorf("Bytes representation of negative infinity incorrect. Expected: %v, Got: %v", expected, actual)
	}
}

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

func TestFloatAbs(t *testing.T) {
	// Test case for positive float
	f := g.Float(3.14)
	expected := g.Float(3.14)
	actual := f.Abs()
	if actual != expected {
		t.Errorf("Absolute value of positive float incorrect. Expected: %v, Got: %v", expected, actual)
	}

	// Test case for negative float
	f = g.Float(-3.14)
	expected = g.Float(3.14)
	actual = f.Abs()
	if actual != expected {
		t.Errorf("Absolute value of negative float incorrect. Expected: %v, Got: %v", expected, actual)
	}

	// Test case for zero float
	f = g.Float(0)
	expected = g.Float(0)
	actual = f.Abs()
	if actual != expected {
		t.Errorf("Absolute value of zero float incorrect. Expected: %v, Got: %v", expected, actual)
	}
}

func TestFloatAdd(t *testing.T) {
	// Test case for addition of positive floats
	f1 := g.Float(3.14)
	f2 := g.Float(1.23)
	expected := g.Float(4.37)
	actual := f1.Add(f2)
	if actual != expected {
		t.Errorf("Addition of positive floats incorrect. Expected: %v, Got: %v", expected, actual)
	}

	// Test case for addition of negative floats
	f1 = g.Float(-3.14)
	f2 = g.Float(-1.23)
	expected = g.Float(-4.37)
	actual = f1.Add(f2)
	if actual != expected {
		t.Errorf("Addition of negative floats incorrect. Expected: %v, Got: %v", expected, actual)
	}
}

func TestFloatToBigFloat(t *testing.T) {
	// Test case for converting positive float to *big.Float
	f := g.Float(3.14)
	expected := big.NewFloat(3.14)
	actual := f.ToBigFloat()
	if actual.Cmp(expected) != 0 {
		t.Errorf("Conversion of positive float to *big.Float incorrect. Expected: %v, Got: %v", expected, actual)
	}

	// Test case for converting negative float to *big.Float
	f = g.Float(-3.14)
	expected = big.NewFloat(-3.14)
	actual = f.ToBigFloat()
	if actual.Cmp(expected) != 0 {
		t.Errorf("Conversion of negative float to *big.Float incorrect. Expected: %v, Got: %v", expected, actual)
	}

	// Test case for converting zero float to *big.Float
	f = g.Float(0)
	expected = big.NewFloat(0)
	actual = f.ToBigFloat()
	if actual.Cmp(expected) != 0 {
		t.Errorf("Conversion of zero float to *big.Float incorrect. Expected: %v, Got: %v", expected, actual)
	}
}

func TestFloatIsZero(t *testing.T) {
	// Test case for zero float
	f := g.Float(0)
	if !f.IsZero() {
		t.Errorf("IsZero method failed to identify zero float.")
	}

	// Test case for positive non-zero float
	f = g.Float(3.14)
	if f.IsZero() {
		t.Errorf("IsZero method incorrectly identified positive non-zero float as zero.")
	}

	// Test case for negative non-zero float
	f = g.Float(-3.14)
	if f.IsZero() {
		t.Errorf("IsZero method incorrectly identified negative non-zero float as zero.")
	}
}

func TestFloatToInt(t *testing.T) {
	// Test case for positive float
	f := g.Float(3.14)
	expected := g.Int(3)
	actual := f.ToInt()
	if actual != expected {
		t.Errorf("ToInt method failed to convert positive float. Expected: %d, Got: %d", expected, actual)
	}

	// Test case for negative float
	f = g.Float(-3.14)
	expected = g.Int(-3)
	actual = f.ToInt()
	if actual != expected {
		t.Errorf("ToInt method failed to convert negative float. Expected: %d, Got: %d", expected, actual)
	}

	// Test case for zero float
	f = g.Float(0)
	expected = g.Int(0)
	actual = f.ToInt()
	if actual != expected {
		t.Errorf("ToInt method failed to convert zero float. Expected: %d, Got: %d", expected, actual)
	}
}

func TestFloatToString(t *testing.T) {
	// Test case for positive float
	f := g.Float(3.14)
	expected := g.String("3.14")
	actual := f.ToString()
	if actual != expected {
		t.Errorf("ToString method failed to convert positive float. Expected: %s, Got: %s", expected, actual)
	}

	// Test case for negative float
	f = g.Float(-3.14)
	expected = g.String("-3.14")
	actual = f.ToString()
	if actual != expected {
		t.Errorf("ToString method failed to convert negative float. Expected: %s, Got: %s", expected, actual)
	}

	// Test case for zero float
	f = g.Float(0)
	expected = g.String("0")
	actual = f.ToString()
	if actual != expected {
		t.Errorf("ToString method failed to convert zero float. Expected: %s, Got: %s", expected, actual)
	}
}

func TestFloatAsFloat32(t *testing.T) {
	// Test case for positive float
	f := g.Float(3.14)
	expected := float32(3.14)
	actual := f.AsFloat32()
	if actual != expected {
		t.Errorf("AsFloat32 method failed to convert positive float. Expected: %f, Got: %f", expected, actual)
	}

	// Test case for negative float
	f = g.Float(-3.14)
	expected = float32(-3.14)
	actual = f.AsFloat32()
	if actual != expected {
		t.Errorf("AsFloat32 method failed to convert negative float. Expected: %f, Got: %f", expected, actual)
	}

	// Test case for zero float
	f = g.Float(0)
	expected = float32(0)
	actual = f.AsFloat32()
	if actual != expected {
		t.Errorf("AsFloat32 method failed to convert zero float. Expected: %f, Got: %f", expected, actual)
	}
}

func TestFloatHashing(t *testing.T) {
	// Test case for a positive float
	f := g.Float(3.14)
	fh := f.Hash()

	// Test MD5
	expectedMD5 := g.String("32200b8781d6e8f31543da4cf19ff307")
	actualMD5 := fh.MD5()
	if actualMD5 != expectedMD5 {
		t.Errorf("MD5 hash mismatch for positive float. Expected: %s, Got: %s", expectedMD5, actualMD5)
	}

	// Test SHA1
	expectedSHA1 := g.String("8d3ad0b5fdf81c2de3656ebe8d8b0f14e1431438")
	actualSHA1 := fh.SHA1()
	if actualSHA1 != expectedSHA1 {
		t.Errorf("SHA1 hash mismatch for positive float. Expected: %s, Got: %s", expectedSHA1, actualSHA1)
	}

	// Test SHA256
	expectedSHA256 := g.String("a7c511f4744a60f88b6a88fbbb1ed7c79820e028f841c50843963bbb1dcdd9f6")
	actualSHA256 := fh.SHA256()
	if actualSHA256 != expectedSHA256 {
		t.Errorf("SHA256 hash mismatch for positive float. Expected: %s, Got: %s", expectedSHA256, actualSHA256)
	}

	// Test SHA512
	expectedSHA512 := g.String(
		"a86ec42eec985ea198240622e13ddfdbd25bee28007d4ee7b17058292dc46ef51e5b107ab44d70ae14300d88bf71a4cda93851ab920f5eeef8bc1531cd451063",
	)
	actualSHA512 := fh.SHA512()
	if actualSHA512 != expectedSHA512 {
		t.Errorf("SHA512 hash mismatch for positive float. Expected: %s, Got: %s", expectedSHA512, actualSHA512)
	}
}
