package g_test

import (
	"math"
	"math/big"
	"testing"

	. "github.com/enetx/g"
)

func TestFloatBytes(t *testing.T) {
	// Test case for positive float
	f := Float(3.14)
	expected := []byte{64, 9, 30, 184, 81, 235, 133, 31} // Bytes representation of 3.14 in big-endian
	actual := f.Bytes()
	if actual.Ne(expected) {
		t.Errorf("Bytes representation of positive float incorrect. Expected: %v, Got: %v", expected, actual)
	}

	// Test case for negative float
	f = Float(-3.14)
	expected = []byte{192, 9, 30, 184, 81, 235, 133, 31} // Bytes representation of -3.14 in big-endian
	actual = f.Bytes()
	if actual.Ne(expected) {
		t.Errorf("Bytes representation of negative float incorrect. Expected: %v, Got: %v", expected, actual)
	}

	// Test case for infinity
	f = Float(math.Inf(1))
	expected = []byte{127, 240, 0, 0, 0, 0, 0, 0} // Bytes representation of positive infinity in big-endian
	actual = f.Bytes()
	if actual.Ne(expected) {
		t.Errorf("Bytes representation of positive infinity incorrect. Expected: %v, Got: %v", expected, actual)
	}

	// Test case for negative infinity
	f = Float(math.Inf(-1))
	expected = []byte{255, 240, 0, 0, 0, 0, 0, 0} // Bytes representation of negative infinity in big-endian
	actual = f.Bytes()
	if actual.Ne(expected) {
		t.Errorf("Bytes representation of negative infinity incorrect. Expected: %v, Got: %v", expected, actual)
	}
}

func TestFloatCompare(t *testing.T) {
	testCases := []struct {
		f1       Float
		f2       Float
		expected int
	}{
		{3.14, 6.28, -1},
		{6.28, 3.14, 1},
		{1.23, 1.23, 0},
		{-2.5, 2.5, -1},
	}

	for _, tc := range testCases {
		result := int(tc.f1.Cmp(tc.f2))
		if result != tc.expected {
			t.Errorf("Compare(%f, %f): expected %d, got %d", tc.f1, tc.f2, tc.expected, result)
		}
	}
}

func TestFloatEq(t *testing.T) {
	testCases := []struct {
		f1       Float
		f2       Float
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
		f1       Float
		f2       Float
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
		f1       Float
		f2       Float
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
		f1       Float
		f2       Float
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

func TestFloatRound(t *testing.T) {
	// Test cases for positive numbers
	positiveTests := []struct {
		input    Float
		expected Int
	}{
		{1.1, 1},
		{1.5, 2},
		{1.9, 2},
	}

	for _, tc := range positiveTests {
		result := tc.input.Round()
		if result != tc.expected {
			t.Errorf("Round(%f) = %d; expected %d", tc.input, result, tc.expected)
		}
	}

	// Test cases for negative numbers
	negativeTests := []struct {
		input    Float
		expected Int
	}{
		{-1.1, -1},
		{-1.5, -2},
		{-1.9, -2},
	}

	for _, tc := range negativeTests {
		result := tc.input.Round()
		if result != tc.expected {
			t.Errorf("Round(%f) = %d; expected %d", tc.input, result, tc.expected)
		}
	}
}

func TestFloatRoundDecimal(t *testing.T) {
	testCases := []struct {
		value    Float
		decimals Int
		expected Float
	}{
		{3.1415926535, 2, 3.14},
		{3.1415926535, 3, 3.142},
		{100.123456789, 4, 100.1235},
		{-5.6789, 1, -5.7},
		{12345.6789, 0, 12346},
		{12345.6789, -1, 12345.6789},
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
	if max := Float(2.2).Max(2.8, 2.1, 2.7); max != 2.8 {
		t.Errorf("Max() = %f, want: %f.", max, 2.8)
	}
}

func TestFloatMin(t *testing.T) {
	if min := Float(2.2).Min(2.8, 2.1, 2.7); min != 2.1 {
		t.Errorf("Min() = %f; want: %f", min, 2.1)
	}
}

func TestFloatAbs(t *testing.T) {
	// Test case for positive float
	f := Float(3.14)
	expected := Float(3.14)
	actual := f.Abs()
	if actual != expected {
		t.Errorf("Absolute value of positive float incorrect. Expected: %v, Got: %v", expected, actual)
	}

	// Test case for negative float
	f = Float(-3.14)
	expected = Float(3.14)
	actual = f.Abs()
	if actual != expected {
		t.Errorf("Absolute value of negative float incorrect. Expected: %v, Got: %v", expected, actual)
	}

	// Test case for zero float
	f = Float(0)
	expected = Float(0)
	actual = f.Abs()
	if actual != expected {
		t.Errorf("Absolute value of zero float incorrect. Expected: %v, Got: %v", expected, actual)
	}
}

func TestFloatAdd(t *testing.T) {
	// Test case for addition of positive floats
	f1 := Float(3.14)
	f2 := Float(1.23)
	expected := Float(4.37)
	actual := f1.Add(f2)
	if actual != expected {
		t.Errorf("Addition of positive floats incorrect. Expected: %v, Got: %v", expected, actual)
	}

	// Test case for addition of negative floats
	f1 = Float(-3.14)
	f2 = Float(-1.23)
	expected = Float(-4.37)
	actual = f1.Add(f2)
	if actual != expected {
		t.Errorf("Addition of negative floats incorrect. Expected: %v, Got: %v", expected, actual)
	}
}

func TestFloatBigFloat(t *testing.T) {
	// Test case for converting positive float to *big.Float
	f := Float(3.14)
	expected := big.NewFloat(3.14)
	actual := f.BigFloat()
	if actual.Cmp(expected) != 0 {
		t.Errorf("Conversion of positive float to *big.Float incorrect. Expected: %v, Got: %v", expected, actual)
	}

	// Test case for converting negative float to *big.Float
	f = Float(-3.14)
	expected = big.NewFloat(-3.14)
	actual = f.BigFloat()
	if actual.Cmp(expected) != 0 {
		t.Errorf("Conversion of negative float to *big.Float incorrect. Expected: %v, Got: %v", expected, actual)
	}

	// Test case for converting zero float to *big.Float
	f = Float(0)
	expected = big.NewFloat(0)
	actual = f.BigFloat()
	if actual.Cmp(expected) != 0 {
		t.Errorf("Conversion of zero float to *big.Float incorrect. Expected: %v, Got: %v", expected, actual)
	}
}

func TestFloatIsZero(t *testing.T) {
	// Test case for zero float
	f := Float(0)
	if !f.IsZero() {
		t.Errorf("IsZero method failed to identify zero float.")
	}

	// Test case for positive non-zero float
	f = Float(3.14)
	if f.IsZero() {
		t.Errorf("IsZero method incorrectly identified positive non-zero float as zero.")
	}

	// Test case for negative non-zero float
	f = Float(-3.14)
	if f.IsZero() {
		t.Errorf("IsZero method incorrectly identified negative non-zero float as zero.")
	}
}

func TestFloatInt(t *testing.T) {
	// Test case for positive float
	f := Float(3.14)
	expected := Int(3)
	actual := f.Int()
	if actual != expected {
		t.Errorf("ToInt method failed to convert positive float. Expected: %d, Got: %d", expected, actual)
	}

	// Test case for negative float
	f = Float(-3.14)
	expected = Int(-3)
	actual = f.Int()
	if actual != expected {
		t.Errorf("ToInt method failed to convert negative float. Expected: %d, Got: %d", expected, actual)
	}

	// Test case for zero float
	f = Float(0)
	expected = Int(0)
	actual = f.Int()
	if actual != expected {
		t.Errorf("ToInt method failed to convert zero float. Expected: %d, Got: %d", expected, actual)
	}
}

func TestFloatString(t *testing.T) {
	// Test case for positive float
	f := Float(3.14)
	expected := String("3.14")
	actual := f.String()
	if actual != expected {
		t.Errorf("ToString method failed to convert positive float. Expected: %s, Got: %s", expected, actual)
	}

	// Test case for negative float
	f = Float(-3.14)
	expected = String("-3.14")
	actual = f.String()
	if actual != expected {
		t.Errorf("ToString method failed to convert negative float. Expected: %s, Got: %s", expected, actual)
	}

	// Test case for zero float
	f = Float(0)
	expected = String("0")
	actual = f.String()
	if actual != expected {
		t.Errorf("ToString method failed to convert zero float. Expected: %s, Got: %s", expected, actual)
	}
}

func TestFloatFloat32(t *testing.T) {
	// Test case for positive float
	f := Float(3.14)
	expected := float32(3.14)
	actual := f.Float32()
	if actual != expected {
		t.Errorf("AsFloat32 method failed to convert positive float. Expected: %f, Got: %f", expected, actual)
	}

	// Test case for negative float
	f = Float(-3.14)
	expected = float32(-3.14)
	actual = f.Float32()
	if actual != expected {
		t.Errorf("AsFloat32 method failed to convert negative float. Expected: %f, Got: %f", expected, actual)
	}

	// Test case for zero float
	f = Float(0)
	expected = float32(0)
	actual = f.Float32()
	if actual != expected {
		t.Errorf("AsFloat32 method failed to convert zero float. Expected: %f, Got: %f", expected, actual)
	}
}

func TestFloatHashing(t *testing.T) {
	// Test case for a positive float
	f := Float(3.14)
	fh := f.Hash()

	// Test MD5
	expectedMD5 := String("32200b8781d6e8f31543da4cf19ff307")
	actualMD5 := fh.MD5()
	if actualMD5 != expectedMD5 {
		t.Errorf("MD5 hash mismatch for positive float. Expected: %s, Got: %s", expectedMD5, actualMD5)
	}

	// Test SHA1
	expectedSHA1 := String("8d3ad0b5fdf81c2de3656ebe8d8b0f14e1431438")
	actualSHA1 := fh.SHA1()
	if actualSHA1 != expectedSHA1 {
		t.Errorf("SHA1 hash mismatch for positive float. Expected: %s, Got: %s", expectedSHA1, actualSHA1)
	}

	// Test SHA256
	expectedSHA256 := String("a7c511f4744a60f88b6a88fbbb1ed7c79820e028f841c50843963bbb1dcdd9f6")
	actualSHA256 := fh.SHA256()
	if actualSHA256 != expectedSHA256 {
		t.Errorf("SHA256 hash mismatch for positive float. Expected: %s, Got: %s", expectedSHA256, actualSHA256)
	}

	// Test SHA512
	expectedSHA512 := String(
		"a86ec42eec985ea198240622e13ddfdbd25bee28007d4ee7b17058292dc46ef51e5b107ab44d70ae14300d88bf71a4cda93851ab920f5eeef8bc1531cd451063",
	)
	actualSHA512 := fh.SHA512()
	if actualSHA512 != expectedSHA512 {
		t.Errorf("SHA512 hash mismatch for positive float. Expected: %s, Got: %s", expectedSHA512, actualSHA512)
	}
}

func TestFloatTransform(t *testing.T) {
	original := Float(3.14)

	multiplyByTwo := func(f Float) Float { return f * 2 }
	expected := Float(6.28)
	result := original.Transform(multiplyByTwo)

	if result != expected {
		t.Errorf("Transform failed: expected %f, got %f", expected, result)
	}

	addConstant := func(f Float) Float { return f + 1.86 }
	expectedWithAddition := Float(5.00)
	resultWithAddition := original.Transform(addConstant)

	if resultWithAddition != expectedWithAddition {
		t.Errorf("Transform with addition failed: expected %f, got %f", expectedWithAddition, resultWithAddition)
	}
}

func TestNewFloat(t *testing.T) {
	f := NewFloat(3.14)
	expected := Float(3.14)
	if f != expected {
		t.Errorf("NewFloat(3.14) should return Float(3.14), got %f", f)
	}
}

func TestFloatDiv(t *testing.T) {
	f := Float(10.0)
	result := f.Div(2.0)
	expected := Float(5.0)
	if result != expected {
		t.Errorf("Div(2.0) should return 5.0, got %f", result)
	}
}

func TestFloatPrint(t *testing.T) {
	f := Float(3.14)
	result := f.Print()
	if result != f {
		t.Errorf("Print() should return original float unchanged")
	}
}

func TestFloatPrintln(t *testing.T) {
	f := Float(3.14)
	result := f.Println()
	if result != f {
		t.Errorf("Println() should return original float unchanged")
	}
}
