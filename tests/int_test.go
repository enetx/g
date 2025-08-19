package g_test

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"testing"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func TestIntAbs(t *testing.T) {
	// Test positive integer
	posInt := Int(5)
	posAbs := posInt.Abs()
	if posAbs != posInt {
		t.Errorf("Abs function incorrect for positive integer. Expected: %d, Got: %d", posInt, posAbs)
	}

	// Test negative integer
	negInt := Int(-5)
	negAbs := negInt.Abs()
	if negAbs != posInt {
		t.Errorf("Abs function incorrect for negative integer. Expected: %d, Got: %d", posInt, negAbs)
	}

	// Test zero
	zero := Int(0)
	zeroAbs := zero.Abs()
	if zeroAbs != zero {
		t.Errorf("Abs function incorrect for zero. Expected: %d, Got: %d", zero, zeroAbs)
	}
}

func TestIntBigInt(t *testing.T) {
	intVal := Int(123)
	bigInt := intVal.BigInt()
	expectedBigInt := big.NewInt(123)
	if bigInt.Cmp(expectedBigInt) != 0 {
		t.Errorf("BigInt function incorrect. Expected: %s, Got: %s", expectedBigInt, bigInt)
	}

	// Test Div
	dividend := Int(10)
	divisor := Int(2)
	quotient := dividend.Div(divisor)
	expectedQuotient := Int(5)
	if quotient != expectedQuotient {
		t.Errorf("Div function incorrect. Expected: %d, Got: %d", expectedQuotient, quotient)
	}
}

func TestIntString(t *testing.T) {
	// Test positive integer
	posInt := Int(123)
	posStr := posInt.String()
	expectedPosStr := String("123")
	if posStr != expectedPosStr {
		t.Errorf("String function incorrect for positive integer. Expected: %s, Got: %s", expectedPosStr, posStr)
	}

	// Test negative integer
	negInt := Int(-123)
	negStr := negInt.String()
	expectedNegStr := String("-123")
	if negStr != expectedNegStr {
		t.Errorf("String function incorrect for negative integer. Expected: %s, Got: %s", expectedNegStr, negStr)
	}

	// Test zero
	zero := Int(0)
	zeroStr := zero.String()
	expectedZeroStr := String("0")
	if zeroStr != expectedZeroStr {
		t.Errorf("ToString function incorrect for zero. Expected: %s, Got: %s", expectedZeroStr, zeroStr)
	}
}

func TestIntInt16(t *testing.T) {
	// Test positive integer within int16 range
	posInt := Int(123)
	posInt16 := posInt.Int16()
	expectedPosInt16 := int16(123)
	if posInt16 != expectedPosInt16 {
		t.Errorf("Int16 function incorrect for positive integer. Expected: %d, Got: %d", expectedPosInt16, posInt16)
	}
}

func TestIntInt32(t *testing.T) {
	// Test positive integer within int32 range
	posInt := Int(123)
	posInt32 := posInt.Int32()
	expectedPosInt32 := int32(123)
	if posInt32 != expectedPosInt32 {
		t.Errorf("Int32 function incorrect for positive integer. Expected: %d, Got: %d", expectedPosInt32, posInt32)
	}

	// Test negative integer within int32 range
	negInt := Int(-123)
	negInt32 := negInt.Int32()
	expectedNegInt32 := int32(-123)
	if negInt32 != expectedNegInt32 {
		t.Errorf("Int32 function incorrect for negative integer. Expected: %d, Got: %d", expectedNegInt32, negInt32)
	}
}

func TestIntInt8(t *testing.T) {
	// Test positive integer within int8 range
	posInt := Int(123)
	posInt8 := posInt.Int8()
	expectedPosInt8 := int8(123)
	if posInt8 != expectedPosInt8 {
		t.Errorf("Int8 function incorrect for positive integer. Expected: %d, Got: %d", expectedPosInt8, posInt8)
	}

	// Test negative integer within int8 range
	negInt := Int(-123)
	negInt8 := negInt.Int8()
	expectedNegInt8 := int8(-123)
	if negInt8 != expectedNegInt8 {
		t.Errorf("Int8 function incorrect for negative integer. Expected: %d, Got: %d", expectedNegInt8, negInt8)
	}

	// Test integer outside int8 range
	bigInt := Int(2000) // larger than int8 max value
	bigInt8 := bigInt.Int8()
	expectedBigInt8 := int8(-48) // expected value after overflow
	if bigInt8 != expectedBigInt8 {
		t.Errorf(
			"Int8 function incorrect for integer outside int8 range. Expected: %d, Got: %d",
			expectedBigInt8,
			bigInt8,
		)
	}
}

func TestIntIsZero(t *testing.T) {
	// Test zero value
	zeroInt := Int(0)
	isZero := zeroInt.IsZero()
	if !isZero {
		t.Errorf("IsZero function incorrect for zero value. Expected: true, Got: %t", isZero)
	}

	// Test non-zero value
	nonZeroInt := Int(123)
	isZero = nonZeroInt.IsZero()
	if isZero {
		t.Errorf("IsZero function incorrect for non-zero value. Expected: false, Got: %t", isZero)
	}
}

func TestIntIsPositive(t *testing.T) {
	tests := []struct {
		name string
		i    Int
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
		i    Int
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
		min := Int(100).Random()
		max := Int(100).Random().Add(min)

		r := min.RandomRange(max)
		if r.Lt(min) || r.Gt(max) {
			t.Errorf("RandomRange(%d, %d) = %d, want in range [%d, %d]", min, max, r, min, max)
		}
	}
}

func TestIntMax(t *testing.T) {
	if max := Int(1).Max(1, 2, 3, 4, 5); max != 5 {
		t.Errorf("Max() = %d, want: %d.", max, 5)
	}
}

func TestIntMin(t *testing.T) {
	if min := Int(1).Min(2, 3, 4, 5); min != 1 {
		t.Errorf("Min() = %d; want: %d", min, 1)
	}
}

func TestIntLte(t *testing.T) {
	// Test for less than
	ltInt1 := Int(5)
	ltInt2 := Int(10)
	isLte := ltInt1.Lte(ltInt2)
	if !isLte {
		t.Errorf("Lte function incorrect for less than. Expected: true, Got: %t", isLte)
	}

	// Test for equal
	eqInt1 := Int(10)
	eqInt2 := Int(10)
	isLte = eqInt1.Lte(eqInt2)
	if !isLte {
		t.Errorf("Lte function incorrect for equal values. Expected: true, Got: %t", isLte)
	}

	// Test for greater than
	gtInt1 := Int(15)
	gtInt2 := Int(10)
	isLte = gtInt1.Lte(gtInt2)
	if isLte {
		t.Errorf("Lte function incorrect for greater than. Expected: false, Got: %t", isLte)
	}
}

func TestIntMul(t *testing.T) {
	// Test for positive multiplication
	posInt1 := Int(5)
	posInt2 := Int(10)
	result := posInt1.Mul(posInt2)
	expected := Int(50)
	if result != expected {
		t.Errorf("Mul function incorrect for positive multiplication. Expected: %d, Got: %d", expected, result)
	}

	// Test for negative multiplication
	negInt1 := Int(-5)
	negInt2 := Int(10)
	result = negInt1.Mul(negInt2)
	expected = Int(-50)
	if result != expected {
		t.Errorf("Mul function incorrect for negative multiplication. Expected: %d, Got: %d", expected, result)
	}
}

func TestIntNe(t *testing.T) {
	// Test for inequality
	ineqInt1 := Int(5)
	ineqInt2 := Int(10)
	isNe := ineqInt1.Ne(ineqInt2)
	if !isNe {
		t.Errorf("Ne function incorrect for inequality. Expected: true, Got: %t", isNe)
	}

	// Test for equality
	eqInt1 := Int(10)
	eqInt2 := Int(10)
	isNe = eqInt1.Ne(eqInt2)
	if isNe {
		t.Errorf("Ne function incorrect for equality. Expected: false, Got: %t", isNe)
	}
}

func TestIntBinary(t *testing.T) {
	// Test for positive integer
	posInt := Int(10)
	binary := posInt.Binary()
	expected := String("00001010")
	if binary != expected {
		t.Errorf("ToBinary function incorrect for positive integer. Expected: %s, Got: %s", expected, binary)
	}

	// Test for negative integer
	negInt := Int(-10)
	binary = negInt.Binary()
	expected = String("-0001010") // Two's complement representation
	if binary != expected {
		t.Errorf("ToBinary function incorrect for negative integer. Expected: %s, Got: %s", expected, binary)
	}

	// Test for zero
	zeroInt := Int(0)
	binary = zeroInt.Binary()
	expected = String("00000000")
	if binary != expected {
		t.Errorf("ToBinary function incorrect for zero. Expected: %s, Got: %s", expected, binary)
	}
}

func TestIntUInt16(t *testing.T) {
	// Test for positive integer
	posInt := Int(100)
	uint16Val := posInt.UInt16()
	expected := uint16(100)
	if uint16Val != expected {
		t.Errorf("UInt16 function incorrect for positive integer. Expected: %d, Got: %d", expected, uint16Val)
	}

	// Test for negative integer
	negInt := Int(-100)
	uint16Val = negInt.UInt16()
	expected = 65436 // Conversion to uint16 of negative number results in 0
	if uint16Val != expected {
		t.Errorf("UInt16 function incorrect for negative integer. Expected: %d, Got: %d", expected, uint16Val)
	}
}

func TestIntUInt32(t *testing.T) {
	// Test for positive integer
	posInt := Int(100)
	uint32Val := posInt.UInt32()
	expected := uint32(100)
	if uint32Val != expected {
		t.Errorf("UInt32 function incorrect for positive integer. Expected: %d, Got: %d", expected, uint32Val)
	}

	// Test for negative integer
	negInt := Int(-100)
	uint32Val = negInt.UInt32()
	expected = 4294967196 // Conversion to uint32 of negative number results in 0
	if uint32Val != expected {
		t.Errorf("UInt32 function incorrect for negative integer. Expected: %d, Got: %d", expected, uint32Val)
	}
}

func TestIntUInt8(t *testing.T) {
	// Test for positive integer within range
	posInt := Int(100)
	uint8Val := posInt.UInt8()
	expected := uint8(100)
	if uint8Val != expected {
		t.Errorf(
			"UInt8 function incorrect for positive integer within range. Expected: %d, Got: %d",
			expected,
			uint8Val,
		)
	}

	// Test for positive integer outside range
	posInt = Int(300)
	uint8Val = posInt.UInt8()
	expected = 44 // Overflow results in 44
	if uint8Val != expected {
		t.Errorf(
			"UInt8 function incorrect for positive integer outside range. Expected: %d, Got: %d",
			expected,
			uint8Val,
		)
	}

	// Test for negative integer
	negInt := Int(-100)
	uint8Val = negInt.UInt8()
	expected = 156 // Conversion to uint8 of negative number results in 156
	if uint8Val != expected {
		t.Errorf("UInt8 function incorrect for negative integer. Expected: %d, Got: %d", expected, uint8Val)
	}
}

func TestIntRem(t *testing.T) {
	// Test cases
	testCases := []struct {
		dividend int
		divisor  int
		expected int
	}{
		{10, 3, 1},    // 10 % 3 = 1
		{15, 7, 1},    // 15 % 7 = 1
		{20, 5, 0},    // 20 % 5 = 0
		{100, 17, 15}, // 100 % 17 = 15
		{35, 11, 2},   // 35 % 11 = 2
		{7, 3, 1},     // 7 % 3 = 1
		{8, 4, 0},     // 8 % 4 = 0
	}

	// Test each case
	for _, tc := range testCases {
		// Wrap the input integers
		i := Int(tc.dividend)
		b := Int(tc.divisor)

		// Call the Rem method
		result := i.Rem(b)

		if result.Std() != tc.expected {
			t.Errorf(
				"Rem function incorrect for %d %% %d. Expected: %d, Got: %d",
				tc.dividend,
				tc.divisor,
				tc.expected,
				result,
			)
		}
	}
}

func TestIntSub(t *testing.T) {
	// Testing subtraction with positive integers
	result := Int(5).Sub(3)
	expected := Int(2)
	if result != expected {
		t.Errorf("Subtraction failed: expected %v, got %v", expected, result)
	}

	// Testing subtraction with negative integers
	result = Int(-5).Sub(-3)
	expected = Int(-2)
	if result != expected {
		t.Errorf("Subtraction failed: expected %v, got %v", expected, result)
	}

	// Testing subtraction with positive and negative integers
	result = Int(5).Sub(-3)
	expected = Int(8)
	if result != expected {
		t.Errorf("Subtraction failed: expected %v, got %v", expected, result)
	}

	// Testing subtraction with negative and positive integers
	result = Int(-5).Sub(3)
	expected = Int(-8)
	if result != expected {
		t.Errorf("Subtraction failed: expected %v, got %v", expected, result)
	}

	// Testing subtraction with zero
	result = Int(0).Sub(0)
	expected = Int(0)
	if result != expected {
		t.Errorf("Subtraction failed: expected %v, got %v", expected, result)
	}
}

func TestIntCmp(t *testing.T) {
	tests := []struct {
		name     string
		i, other Int
		expected cmp.Ordering
	}{
		{"LessThan", Int(5), Int(10), cmp.Less},
		{"GreaterThan", Int(15), Int(10), cmp.Greater},
		{"EqualTo", Int(10), Int(10), cmp.Equal},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.i.Cmp(test.other)
			if result != test.expected {
				t.Errorf("%s: Expected %v, Got %v", test.name, test.expected, result)
			}
		})
	}
}

func TestIntTransform(t *testing.T) {
	original := Int(10)

	multiplyByThree := func(i Int) Int { return i * 3 }
	expected := Int(30)
	result := original.Transform(multiplyByThree)

	if result != expected {
		t.Errorf("Transform failed: expected %d, got %d", expected, result)
	}

	subtractFive := func(i Int) Int { return i - 5 }
	expectedWithSubtraction := Int(5)
	resultWithSubtraction := original.Transform(subtractFive)

	if resultWithSubtraction != expectedWithSubtraction {
		t.Errorf(
			"Transform with subtraction failed: expected %d, got %d",
			expectedWithSubtraction,
			resultWithSubtraction,
		)
	}
}

func TestIntPrint(t *testing.T) {
	i := Int(42)
	result := i.Print()
	if result != i {
		t.Errorf("Print() should return original int unchanged")
	}
}

func TestIntPrintln(t *testing.T) {
	i := Int(42)
	result := i.Println()
	if result != i {
		t.Errorf("Println() should return original int unchanged")
	}
}

func TestRandomRange_Distribution_Sanity(t *testing.T) {
	const (
		lo    = Int(-3)
		hi    = Int(3)
		iters = 1_000_00 // 100k
	)
	counts := make(map[Int]int)

	for range iters {
		x := lo.RandomRange(hi)
		counts[x]++
	}

	want := float64(iters) / float64(hi-lo+1)
	tol := want * 0.05
	for v := lo; v <= hi; v++ {
		got := float64(counts[v])
		if diff := got - want; diff < -tol || diff > tol {
			t.Fatalf("value %d: got %d, want ~%.0f (+/-%0.f)", v, counts[v], want, tol)
		}
	}
}

func TestRandomRange_FullRange_Moves(t *testing.T) {
	lo := Int(math.MinInt64)
	hi := Int(math.MaxInt64)
	a := lo.RandomRange(hi)
	b := lo.RandomRange(hi)
	if a == b {
		t.Logf("two draws equal (ok but unlikely); a=%d b=%d", a, b)
	}
}

func wantBytesBE(i int64) Bytes {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(i))

	// Remove leading zeros but preserve sign bit
	start := 0
	for start < 7 && buf[start] == 0 {
		start++
	}
	// For positive numbers, if MSB is set, need extra zero byte
	if i >= 0 && buf[start]&0x80 != 0 {
		start--
	}
	// For negative numbers, if we removed too many bytes and lost sign bit
	if i < 0 && start > 0 && buf[start]&0x80 == 0 {
		start--
	}
	return Bytes(buf[start:])
}

func wantBytesLE(i int64) Bytes {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(i))

	// Remove trailing zeros but preserve sign bit
	end := 8
	for end > 1 && buf[end-1] == 0 {
		end--
	}
	// For positive numbers, if MSB is set, need extra zero byte
	if i >= 0 && buf[end-1]&0x80 != 0 {
		end++
	}
	// For negative numbers, if we removed too many bytes and lost sign bit
	if i < 0 && end < 8 && buf[end-1]&0x80 == 0 {
		end++
	}
	return Bytes(buf[:end])
}

func TestIntBytes_Orders(t *testing.T) {
	type tc struct {
		name string
		in   int64
	}
	cases := []tc{
		{name: "zero", in: 0},
		{name: "single positive", in: 5},
		{name: "single negative -1", in: -1},
		{name: "single negative -128", in: -128},
		{name: "positive needs sign bit 255", in: 255},
		{name: "positive needs sign bit 32767", in: 32767},
		{name: "negative -2", in: -2},
		{name: "positive 3 bytes worth", in: 0x010203},
		{name: "negative 3 bytes worth", in: -0x010203},
		{name: "large positive", in: 0x123456789ABCDEF0},
		{name: "large negative", in: -0x123456789ABCDEF0},
		{name: "max int64", in: 9223372036854775807},        // 0x7FFFFFFFFFFFFFFF
		{name: "min int64", in: -9223372036854775808},       // 0x8000000000000000
		{name: "positive MSB set 32768", in: 32768},         // 0x8000
		{name: "negative -32768", in: -32768},               // 0x8000 as negative
		{name: "positive MSB set single byte 128", in: 128}, // 0x80
	}

	for _, c := range cases {
		t.Run("BE/"+c.name, func(t *testing.T) {
			want := wantBytesBE(c.in)
			got := Int(c.in).BytesBE()
			if len(got) != len(want) {
				t.Fatalf("BytesBE(%d): length mismatch, want %v (len=%d), got %v (len=%d)",
					c.in, want, len(want), got, len(got))
			}
			for i := range got {
				if got[i] != want[i] {
					t.Fatalf("BytesBE(%d): want %v, got %v", c.in, want, got)
				}
			}
		})
		t.Run("LE/"+c.name, func(t *testing.T) {
			want := wantBytesLE(c.in)
			got := Int(c.in).BytesLE()
			if len(got) != len(want) {
				t.Fatalf("BytesLE(%d): length mismatch, want %v (len=%d), got %v (len=%d)",
					c.in, want, len(want), got, len(got))
			}
			for i := range got {
				if got[i] != want[i] {
					t.Fatalf("BytesLE(%d): want %v, got %v", c.in, want, got)
				}
			}
		})
	}
}

// Round-trip tests to ensure BytesBE/LE and IntBE/LE are inverses
func TestIntBytes_RoundTrip(t *testing.T) {
	testValues := []int64{
		0, 1, -1, 127, 128, 255, 256, -128, -129, -255, -256,
		32767, 32768, -32768, -32769,
		0x7FFFFF, 0x800000, -0x800000, -0x800001,
		0x7FFFFFFF, 0x80000000, -0x80000000, -0x80000001,
		0x7FFFFFFFFFFF, 0x800000000000, -0x800000000000, -0x800000000001,
		0x7FFFFFFFFFFFFFFF, -0x8000000000000000, // max/min int64
		0x123456789ABCDEF0, -0x123456789ABCDEF0,
	}

	for _, val := range testValues {
		t.Run(fmt.Sprintf("BE_%d", val), func(t *testing.T) {
			bytes := Int(val).BytesBE()
			back := bytes.IntBE()
			if int64(back) != val {
				t.Fatalf("Round-trip BE failed: %d -> %v -> %d", val, bytes, back)
			}
		})

		t.Run(fmt.Sprintf("LE_%d", val), func(t *testing.T) {
			bytes := Int(val).BytesLE()
			back := bytes.IntLE()
			if int64(back) != val {
				t.Fatalf("Round-trip LE failed: %d -> %v -> %d", val, bytes, back)
			}
		})
	}
}

func TestIntRandom_EdgeCase(t *testing.T) {
	// Test edge case where i <= 0 should return 0
	zero := Int(0)
	result := zero.Random()
	if result != 0 {
		t.Errorf("Random(0) should return 0, got %d", result)
	}

	negative := Int(-5)
	result = negative.Random()
	if result != 0 {
		t.Errorf("Random(-5) should return 0, got %d", result)
	}
}

func TestIntRandom_NormalCase(t *testing.T) {
	// Test normal case where i > 0 should return random number in range [0, i)
	ten := Int(10)
	result := ten.Random()
	if result < 0 || result >= 10 {
		t.Errorf("Random(10) should return value in [0, 10), got %d", result)
	}

	one := Int(1)
	result = one.Random()
	if result != 0 {
		t.Errorf("Random(1) should return 0, got %d", result)
	}
}

func TestIntUInt64(t *testing.T) {
	testCases := []struct {
		input    Int
		expected uint64
	}{
		{Int(0), uint64(0)},
		{Int(42), uint64(42)},
		{Int(123456789), uint64(123456789)},
		{Int(9223372036854775807), uint64(9223372036854775807)}, // max int64
		{Int(-1), uint64(18446744073709551615)},                 // -1 as uint64
		{Int(-42), uint64(18446744073709551574)},                // -42 as uint64
	}

	for _, tc := range testCases {
		result := tc.input.UInt64()
		if result != tc.expected {
			t.Errorf("UInt64() for %d: expected %d, got %d", tc.input, tc.expected, result)
		}
	}
}
