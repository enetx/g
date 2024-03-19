package g_test

import (
	"math/big"
	"testing"

	"github.com/enetx/g"
)

func TestIntAbs(t *testing.T) {
	// Test positive integer
	posInt := g.Int(5)
	posAbs := posInt.Abs()
	if posAbs != posInt {
		t.Errorf("Abs function incorrect for positive integer. Expected: %d, Got: %d", posInt, posAbs)
	}

	// Test negative integer
	negInt := g.Int(-5)
	negAbs := negInt.Abs()
	if negAbs != posInt {
		t.Errorf("Abs function incorrect for negative integer. Expected: %d, Got: %d", posInt, negAbs)
	}

	// Test zero
	zero := g.Int(0)
	zeroAbs := zero.Abs()
	if zeroAbs != zero {
		t.Errorf("Abs function incorrect for zero. Expected: %d, Got: %d", zero, zeroAbs)
	}
}

func TestIntToInt(t *testing.T) {
	// Test ToBigInt
	intVal := g.Int(123)
	bigInt := intVal.ToBigInt()
	expectedBigInt := big.NewInt(123)
	if bigInt.Cmp(expectedBigInt) != 0 {
		t.Errorf("ToBigInt function incorrect. Expected: %s, Got: %s", expectedBigInt, bigInt)
	}

	// Test Div
	dividend := g.Int(10)
	divisor := g.Int(2)
	quotient := dividend.Div(divisor)
	expectedQuotient := g.Int(5)
	if quotient != expectedQuotient {
		t.Errorf("Div function incorrect. Expected: %d, Got: %d", expectedQuotient, quotient)
	}
}

func TestIntToString(t *testing.T) {
	// Test positive integer
	posInt := g.Int(123)
	posStr := posInt.ToString()
	expectedPosStr := g.NewString("123")
	if posStr != expectedPosStr {
		t.Errorf("ToString function incorrect for positive integer. Expected: %s, Got: %s", expectedPosStr, posStr)
	}

	// Test negative integer
	negInt := g.Int(-123)
	negStr := negInt.ToString()
	expectedNegStr := g.NewString("-123")
	if negStr != expectedNegStr {
		t.Errorf("ToString function incorrect for negative integer. Expected: %s, Got: %s", expectedNegStr, negStr)
	}

	// Test zero
	zero := g.Int(0)
	zeroStr := zero.ToString()
	expectedZeroStr := g.NewString("0")
	if zeroStr != expectedZeroStr {
		t.Errorf("ToString function incorrect for zero. Expected: %s, Got: %s", expectedZeroStr, zeroStr)
	}
}

func TestIntAsInt16(t *testing.T) {
	// Test positive integer within int16 range
	posInt := g.Int(123)
	posInt16 := posInt.AsInt16()
	expectedPosInt16 := int16(123)
	if posInt16 != expectedPosInt16 {
		t.Errorf("AsInt16 function incorrect for positive integer. Expected: %d, Got: %d", expectedPosInt16, posInt16)
	}
}

func TestIntAsInt32(t *testing.T) {
	// Test positive integer within int32 range
	posInt := g.Int(123)
	posInt32 := posInt.AsInt32()
	expectedPosInt32 := int32(123)
	if posInt32 != expectedPosInt32 {
		t.Errorf("AsInt32 function incorrect for positive integer. Expected: %d, Got: %d", expectedPosInt32, posInt32)
	}

	// Test negative integer within int32 range
	negInt := g.Int(-123)
	negInt32 := negInt.AsInt32()
	expectedNegInt32 := int32(-123)
	if negInt32 != expectedNegInt32 {
		t.Errorf("AsInt32 function incorrect for negative integer. Expected: %d, Got: %d", expectedNegInt32, negInt32)
	}
}

func TestIntAsInt8(t *testing.T) {
	// Test positive integer within int8 range
	posInt := g.Int(123)
	posInt8 := posInt.AsInt8()
	expectedPosInt8 := int8(123)
	if posInt8 != expectedPosInt8 {
		t.Errorf("AsInt8 function incorrect for positive integer. Expected: %d, Got: %d", expectedPosInt8, posInt8)
	}

	// Test negative integer within int8 range
	negInt := g.Int(-123)
	negInt8 := negInt.AsInt8()
	expectedNegInt8 := int8(-123)
	if negInt8 != expectedNegInt8 {
		t.Errorf("AsInt8 function incorrect for negative integer. Expected: %d, Got: %d", expectedNegInt8, negInt8)
	}

	// Test integer outside int8 range
	bigInt := g.Int(2000) // larger than int8 max value
	bigInt8 := bigInt.AsInt8()
	expectedBigInt8 := int8(-48) // expected value after overflow
	if bigInt8 != expectedBigInt8 {
		t.Errorf(
			"AsInt8 function incorrect for integer outside int8 range. Expected: %d, Got: %d",
			expectedBigInt8,
			bigInt8,
		)
	}
}

func TestIntIsZero(t *testing.T) {
	// Test zero value
	zeroInt := g.Int(0)
	isZero := zeroInt.IsZero()
	if !isZero {
		t.Errorf("IsZero function incorrect for zero value. Expected: true, Got: %t", isZero)
	}

	// Test non-zero value
	nonZeroInt := g.Int(123)
	isZero = nonZeroInt.IsZero()
	if isZero {
		t.Errorf("IsZero function incorrect for non-zero value. Expected: false, Got: %t", isZero)
	}
}

func TestIntIsPositive(t *testing.T) {
	tests := []struct {
		name string
		i    g.Int
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
		i    g.Int
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
		min := g.NewInt(100).Random()
		max := g.NewInt(100).Random().Add(min)

		r := min.RandomRange(max)
		if r.Lt(min) || r.Gt(max) {
			t.Errorf("RandomRange(%d, %d) = %d, want in range [%d, %d]", min, max, r, min, max)
		}
	}
}

func TestIntMax(t *testing.T) {
	if max := g.NewInt(1).Max(1, 2, 3, 4, 5); max != 5 {
		t.Errorf("Max() = %d, want: %d.", max, 5)
	}
}

func TestIntMin(t *testing.T) {
	if min := g.NewInt(1).Min(2, 3, 4, 5); min != 1 {
		t.Errorf("Min() = %d; want: %d", min, 1)
	}
}

func TestIntLte(t *testing.T) {
	// Test for less than
	ltInt1 := g.Int(5)
	ltInt2 := g.Int(10)
	isLte := ltInt1.Lte(ltInt2)
	if !isLte {
		t.Errorf("Lte function incorrect for less than. Expected: true, Got: %t", isLte)
	}

	// Test for equal
	eqInt1 := g.Int(10)
	eqInt2 := g.Int(10)
	isLte = eqInt1.Lte(eqInt2)
	if !isLte {
		t.Errorf("Lte function incorrect for equal values. Expected: true, Got: %t", isLte)
	}

	// Test for greater than
	gtInt1 := g.Int(15)
	gtInt2 := g.Int(10)
	isLte = gtInt1.Lte(gtInt2)
	if isLte {
		t.Errorf("Lte function incorrect for greater than. Expected: false, Got: %t", isLte)
	}
}

func TestIntMul(t *testing.T) {
	// Test for positive multiplication
	posInt1 := g.Int(5)
	posInt2 := g.Int(10)
	result := posInt1.Mul(posInt2)
	expected := g.Int(50)
	if result != expected {
		t.Errorf("Mul function incorrect for positive multiplication. Expected: %d, Got: %d", expected, result)
	}

	// Test for negative multiplication
	negInt1 := g.Int(-5)
	negInt2 := g.Int(10)
	result = negInt1.Mul(negInt2)
	expected = g.Int(-50)
	if result != expected {
		t.Errorf("Mul function incorrect for negative multiplication. Expected: %d, Got: %d", expected, result)
	}
}

func TestIntNe(t *testing.T) {
	// Test for inequality
	ineqInt1 := g.Int(5)
	ineqInt2 := g.Int(10)
	isNe := ineqInt1.Ne(ineqInt2)
	if !isNe {
		t.Errorf("Ne function incorrect for inequality. Expected: true, Got: %t", isNe)
	}

	// Test for equality
	eqInt1 := g.Int(10)
	eqInt2 := g.Int(10)
	isNe = eqInt1.Ne(eqInt2)
	if isNe {
		t.Errorf("Ne function incorrect for equality. Expected: false, Got: %t", isNe)
	}
}

func TestIntToBinary(t *testing.T) {
	// Test for positive integer
	posInt := g.Int(10)
	binary := posInt.ToBinary()
	expected := g.String("00001010")
	if binary != expected {
		t.Errorf("ToBinary function incorrect for positive integer. Expected: %s, Got: %s", expected, binary)
	}

	// Test for negative integer
	negInt := g.Int(-10)
	binary = negInt.ToBinary()
	expected = g.String("-0001010") // Two's complement representation
	if binary != expected {
		t.Errorf("ToBinary function incorrect for negative integer. Expected: %s, Got: %s", expected, binary)
	}

	// Test for zero
	zeroInt := g.Int(0)
	binary = zeroInt.ToBinary()
	expected = g.String("00000000")
	if binary != expected {
		t.Errorf("ToBinary function incorrect for zero. Expected: %s, Got: %s", expected, binary)
	}
}

func TestIntAsUInt16(t *testing.T) {
	// Test for positive integer
	posInt := g.Int(100)
	uint16Val := posInt.AsUInt16()
	expected := uint16(100)
	if uint16Val != expected {
		t.Errorf("AsUInt16 function incorrect for positive integer. Expected: %d, Got: %d", expected, uint16Val)
	}

	// Test for negative integer
	negInt := g.Int(-100)
	uint16Val = negInt.AsUInt16()
	expected = 65436 // Conversion to uint16 of negative number results in 0
	if uint16Val != expected {
		t.Errorf("AsUInt16 function incorrect for negative integer. Expected: %d, Got: %d", expected, uint16Val)
	}
}

func TestIntAsUInt32(t *testing.T) {
	// Test for positive integer
	posInt := g.Int(100)
	uint32Val := posInt.AsUInt32()
	expected := uint32(100)
	if uint32Val != expected {
		t.Errorf("AsUInt32 function incorrect for positive integer. Expected: %d, Got: %d", expected, uint32Val)
	}

	// Test for negative integer
	negInt := g.Int(-100)
	uint32Val = negInt.AsUInt32()
	expected = 4294967196 // Conversion to uint32 of negative number results in 0
	if uint32Val != expected {
		t.Errorf("AsUInt32 function incorrect for negative integer. Expected: %d, Got: %d", expected, uint32Val)
	}
}

func TestIntAsUInt8(t *testing.T) {
	// Test for positive integer within range
	posInt := g.Int(100)
	uint8Val := posInt.AsUInt8()
	expected := uint8(100)
	if uint8Val != expected {
		t.Errorf(
			"AsUInt8 function incorrect for positive integer within range. Expected: %d, Got: %d",
			expected,
			uint8Val,
		)
	}

	// Test for positive integer outside range
	posInt = g.Int(300)
	uint8Val = posInt.AsUInt8()
	expected = 44 // Overflow results in 44
	if uint8Val != expected {
		t.Errorf(
			"AsUInt8 function incorrect for positive integer outside range. Expected: %d, Got: %d",
			expected,
			uint8Val,
		)
	}

	// Test for negative integer
	negInt := g.Int(-100)
	uint8Val = negInt.AsUInt8()
	expected = 156 // Conversion to uint8 of negative number results in 156
	if uint8Val != expected {
		t.Errorf("AsUInt8 function incorrect for negative integer. Expected: %d, Got: %d", expected, uint8Val)
	}
}

func TestIntHashingFunctions(t *testing.T) {
	// Test case for SHA1 hashing
	input := g.Int(42)
	expectedSHA1 := "df58248c414f342c81e056b40bee12d17a08bf61"
	sha1Hash := input.Hash().SHA1().Std()
	if sha1Hash != expectedSHA1 {
		t.Errorf("SHA1 hashing failed. Expected: %s, Got: %s", expectedSHA1, sha1Hash)
	}

	// Test case for SHA256 hashing
	expectedSHA256 := "684888c0ebb17f374298b65ee2807526c066094c701bcc7ebbe1c1095f494fc1"
	sha256Hash := input.Hash().SHA256().Std()
	if sha256Hash != expectedSHA256 {
		t.Errorf("SHA256 hashing failed. Expected: %s, Got: %s", expectedSHA256, sha256Hash)
	}

	// Test case for SHA512 hashing
	expectedSHA512 := "7846cdd4c2b9052768b8901640122e5282e0b833a6a58312a7763472d448ee23781c7f08d90793fdfe71ffe74238cf6e4aa778cc9bb8cec03ea7268d4893a502"
	sha512Hash := input.Hash().SHA512().Std()
	if sha512Hash != expectedSHA512 {
		t.Errorf("SHA512 hashing failed. Expected: %s, Got: %s", expectedSHA512, sha512Hash)
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
		i := g.Int(tc.dividend)
		b := g.Int(tc.divisor)

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
