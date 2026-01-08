package g_test

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"testing"

	. "github.com/enetx/g"
)

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

func wantFloatBytesBE(f float64) Bytes {
	var buf [8]byte
	bits := math.Float64bits(f)
	binary.BigEndian.PutUint64(buf[:], bits)
	return Bytes(buf[:])
}

func wantFloatBytesLE(f float64) Bytes {
	var buf [8]byte
	bits := math.Float64bits(f)
	binary.LittleEndian.PutUint64(buf[:], bits)
	return Bytes(buf[:])
}

func TestFloatBytes_Orders(t *testing.T) {
	type tc struct {
		name string
		in   float64
	}
	cases := []tc{
		{name: "zero", in: 0.0},
		{name: "negative zero", in: math.Copysign(0, -1)},
		{name: "positive small", in: 3.14159},
		{name: "negative small", in: -3.14159},
		{name: "positive integer", in: 42.0},
		{name: "negative integer", in: -42.0},
		{name: "very small positive", in: 1e-10},
		{name: "very small negative", in: -1e-10},
		{name: "very large positive", in: 1e10},
		{name: "very large negative", in: -1e10},
		{name: "max float64", in: math.MaxFloat64},
		{name: "smallest positive float64", in: math.SmallestNonzeroFloat64},
		{name: "positive infinity", in: math.Inf(1)},
		{name: "negative infinity", in: math.Inf(-1)},
		{name: "NaN", in: math.NaN()},
		{name: "pi", in: math.Pi},
		{name: "e", in: math.E},
		{name: "fraction", in: 1.0 / 3.0},
		{name: "large fraction", in: 123456789.987654321},
		{name: "scientific notation", in: 1.23456789e-100},
	}

	for _, c := range cases {
		t.Run("BytesBE/"+c.name, func(t *testing.T) {
			want := wantFloatBytesBE(c.in)
			got := Float(c.in).BytesBE()
			if len(got) != len(want) {
				t.Fatalf("BytesBE(%g): length mismatch, want %v (len=%d), got %v (len=%d)",
					c.in, want, len(want), got, len(got))
			}
			for i := range got {
				if got[i] != want[i] {
					t.Fatalf("BytesBE(%g): want %v, got %v", c.in, want, got)
				}
			}
		})

		t.Run("BytesLE/"+c.name, func(t *testing.T) {
			want := wantFloatBytesLE(c.in)
			got := Float(c.in).BytesLE()
			if len(got) != len(want) {
				t.Fatalf("BytesLE(%g): length mismatch, want %v (len=%d), got %v (len=%d)",
					c.in, want, len(want), got, len(got))
			}
			for i := range got {
				if got[i] != want[i] {
					t.Fatalf("BytesLE(%g): want %v, got %v", c.in, want, got)
				}
			}
		})
	}
}

// Round-trip tests to ensure BytesBE/LE and FloatBE/LE are inverses
func TestFloatBytes_RoundTrip(t *testing.T) {
	testValues := []float64{
		0.0, -0.0, 1.0, -1.0, 0.5, -0.5,
		math.Pi, math.E, math.Sqrt2, math.Ln2, math.Log2E,
		42.0, -42.0, 123.456, -123.456,
		1e-10, -1e-10, 1e10, -1e10,
		1e-100, -1e-100, 1e100, -1e100,
		math.MaxFloat64, math.SmallestNonzeroFloat64,
		math.Inf(1), math.Inf(-1), math.NaN(),
		1.0 / 3.0, 2.0 / 3.0, -1.0 / 3.0, -2.0 / 3.0,
		3.141592653589793, 2.718281828459045,
		1.23456789e-50, -9.87654321e50,
		0.1, 0.2, 0.3, 0.7, 0.9, // Decimal fractions
	}

	for i, val := range testValues {
		t.Run(fmt.Sprintf("BE_%d_%.10g", i, val), func(t *testing.T) {
			bytes := Float(val).BytesBE()
			back := float64(bytes.FloatBE())

			// Special handling for NaN since NaN != NaN
			if math.IsNaN(val) {
				if !math.IsNaN(back) {
					t.Fatalf("Round-trip BE failed: NaN -> %v -> %g (not NaN)", bytes, back)
				}
				return
			}

			// For negative zero, check signbit
			if val == 0 && math.Signbit(val) {
				if back != 0 || !math.Signbit(back) {
					t.Fatalf("Round-trip BE failed: -0 -> %v -> %g (signbit: %t)", bytes, back, math.Signbit(back))
				}
				return
			}

			if back != val {
				t.Fatalf("Round-trip BE failed: %g -> %v -> %g", val, bytes, back)
			}
		})

		t.Run(fmt.Sprintf("LE_%d_%.10g", i, val), func(t *testing.T) {
			bytes := Float(val).BytesLE()
			back := float64(bytes.FloatLE())

			// Special handling for NaN since NaN != NaN
			if math.IsNaN(val) {
				if !math.IsNaN(back) {
					t.Fatalf("Round-trip LE failed: NaN -> %v -> %g (not NaN)", bytes, back)
				}
				return
			}

			// For negative zero, check signbit
			if val == 0 && math.Signbit(val) {
				if back != 0 || !math.Signbit(back) {
					t.Fatalf("Round-trip LE failed: -0 -> %v -> %g (signbit: %t)", bytes, back, math.Signbit(back))
				}
				return
			}

			if back != val {
				t.Fatalf("Round-trip LE failed: %g -> %v -> %g", val, bytes, back)
			}
		})
	}
}

// Test endianness consistency
func TestFloatBytes_Endianness(t *testing.T) {
	testValues := []float64{
		0.0, 1.0, -1.0, math.Pi, math.E, 123.456, -123.456,
		math.Inf(1), math.Inf(-1), math.NaN(),
		math.MaxFloat64, math.SmallestNonzeroFloat64,
	}

	for i, val := range testValues {
		t.Run(fmt.Sprintf("Endianness_%d_%.6g", i, val), func(t *testing.T) {
			bytesBE := Float(val).BytesBE()
			bytesLE := Float(val).BytesLE()

			// Verify that BE and LE are byte-reversed versions of each other
			if len(bytesBE) != 8 || len(bytesLE) != 8 {
				t.Fatalf("Expected 8 bytes, got BE: %d, LE: %d", len(bytesBE), len(bytesLE))
			}

			for j := 0; j < 8; j++ {
				if bytesBE[j] != bytesLE[7-j] {
					t.Fatalf("Byte reversal failed at position %d: BE[%d]=%d, LE[%d]=%d",
						j, j, bytesBE[j], 7-j, bytesLE[7-j])
				}
			}

			// Verify both convert back to the same value
			backBE := float64(bytesBE.FloatBE())
			backLE := float64(bytesLE.FloatLE())

			// Handle NaN specially
			if math.IsNaN(val) {
				if !math.IsNaN(backBE) || !math.IsNaN(backLE) {
					t.Fatalf("NaN not preserved: original NaN, backBE: %g (isNaN: %t), backLE: %g (isNaN: %t)",
						backBE, math.IsNaN(backBE), backLE, math.IsNaN(backLE))
				}
				return
			}

			// Handle negative zero specially
			if val == 0 && math.Signbit(val) {
				if backBE != 0 || !math.Signbit(backBE) || backLE != 0 || !math.Signbit(backLE) {
					t.Fatalf("Negative zero not preserved: backBE: %g (signbit: %t), backLE: %g (signbit: %t)",
						backBE, math.Signbit(backBE), backLE, math.Signbit(backLE))
				}
				return
			}

			if backBE != val || backLE != val {
				t.Fatalf("Value not preserved: original: %g, backBE: %g, backLE: %g", val, backBE, backLE)
			}
		})
	}
}

func TestFloatBits(t *testing.T) {
	testCases := []Float{
		0.0,
		1.0,
		-1.0,
		3.14159,
		-3.14159,
		Float(math.Inf(1)),
		Float(math.Inf(-1)),
		Float(math.NaN()),
		Float(math.SmallestNonzeroFloat64),
		Float(math.MaxFloat64),
	}

	for _, f := range testCases {
		expected := math.Float64bits(f.Std())
		result := f.Bits()

		if result != expected {
			t.Errorf("Bits() for %v: expected %d, got %d", f, expected, result)
		}
	}
}

func TestFloatScan(t *testing.T) {
	var f Float

	if err := f.Scan(nil); err != nil {
		t.Fatalf("Scan(nil) error: %v", err)
	}
	if f != 0 {
		t.Fatalf("Expected 0, got %v", f)
	}

	if err := f.Scan(3.14); err != nil {
		t.Fatalf("Scan(3.14) error: %v", err)
	}
	if f != 3.14 {
		t.Fatalf("Expected 3.14, got %v", f)
	}

	err := f.Scan("not a float")
	if err == nil {
		t.Fatal("Expected error for unsupported type")
	}
}

func TestFloatValue(t *testing.T) {
	f := Float(2.71)
	val, err := f.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if fv, ok := val.(float64); !ok || fv != float64(f) {
		t.Fatalf("Expected %v, got %v", f, val)
	}

	var zero Float
	val2, err := zero.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if val2 != float64(0) {
		t.Fatalf("Expected 0, got %v", val2)
	}
}
