package g

import (
	"cmp"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"math/bits"
	"strconv"

	"github.com/enetx/g/pkg/minmax"
)

// NewFloat creates a new Float with the provided value.
func NewFloat[T float64 | float32 | ~int](float T) Float { return Float(float) }

// Bytes returns the Float as a byte slice.
func (f Float) Bytes() Bytes {
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, f.ToUInt64())

	return buffer[bits.LeadingZeros64(f.ToUInt64())>>3:]
}

// Min returns the minimum of two Floats.
func (f Float) Min(b ...Float) Float { return minmax.Min(f, b...) }

// Max returns the maximum of two Floats.
func (f Float) Max(b ...Float) Float { return minmax.Max(f, b...) }

// Abs returns the absolute value of the Float.
func (f Float) Abs() Float { return Float(math.Abs(f.Std())) }

// Add adds two Floats and returns the result.
func (f Float) Add(b Float) Float { return f + b }

// ToBigFloat returns the Float as a *big.Float.
func (f Float) ToBigFloat() *big.Float { return big.NewFloat(f.Std()) }

// Compare compares two Floats and returns an Int.
func (f Float) Compare(b Float) Int { return Int(cmp.Compare(f, b)) }

// Div divides two Floats and returns the result.
func (f Float) Div(b Float) Float { return f / b }

// IsZero checks if the Float is 0.
func (f Float) IsZero() bool { return f.Eq(0) }

// Eq checks if two Floats are equal.
func (f Float) Eq(b Float) bool { return f.Compare(b).Eq(0) }

// Std returns the Float as a float64.
func (f Float) Std() float64 { return float64(f) }

// Gt checks if the Float is greater than the specified Float.
func (f Float) Gt(b Float) bool { return f.Compare(b).Gt(0) }

// ToInt returns the Float as an Int.
func (f Float) ToInt() Int { return Int(f) }

// ToString returns the Float as an String.
func (f Float) ToString() String { return String(strconv.FormatFloat(f.Std(), 'f', -1, 64)) }

// Lt checks if the Float is less than the specified Float.
func (f Float) Lt(b Float) bool { return f.Compare(b).Lt(0) }

// Mul multiplies two Floats and returns the result.
func (f Float) Mul(b Float) Float { return f * b }

// Ne checks if two Floats are not equal.
func (f Float) Ne(b Float) bool { return !f.Eq(b) }

// Round rounds the Float to the nearest integer and returns the result as an Int.
// func (f Float) Round() Int { return Int(math.Round(f.Std())) }
func (f Float) Round() Int {
	if f >= 0 {
		return Int(f + 0.5)
	}

	return Int(f - 0.5)
}

// RoundDecimal rounds the Float value to the specified number of decimal places.
//
// The function takes the number of decimal places (precision) as an argument and returns a new
// Float value rounded to that number of decimals. This is achieved by multiplying the Float
// value by a power of 10 equal to the desired precision, rounding the result, and then dividing
// the rounded result by the same power of 10.
//
// Parameters:
//
// - precision (int): The number of decimal places to round the Float value to.
//
// Returns:
//
// - Float: A new Float value rounded to the specified number of decimal places.
//
// Example usage:
//
//	f := g.Float(3.14159)
//	rounded := f.RoundDecimal(2) // rounded will be 3.14
func (f Float) RoundDecimal(precision int) Float {
	if precision < 0 {
		return f
	}

	mult := 1
	for i := 0; i < precision; i++ {
		mult *= 10
	}

	result := f * Float(mult)
	if result >= 0 {
		result += 0.5
	} else {
		result -= 0.5
	}

	return Float(int(result)) / Float(mult)
}

// Sub subtracts two Floats and returns the result.
func (f Float) Sub(b Float) Float { return f - b }

// ToUInt64 returns the Float as a uint64.
func (f Float) ToUInt64() uint64 { return math.Float64bits(f.Std()) }

// AsFloat32 returns the Float as a float32.
func (f Float) AsFloat32() float32 { return float32(f) }

// Print prints the value of the Float to the standard output (console)
// and returns the Float unchanged.
func (f Float) Print() Float { fmt.Println(f); return f }
