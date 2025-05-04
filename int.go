package g

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"math/bits"
	"strconv"

	"github.com/enetx/g/cmp"
	"github.com/enetx/g/pkg/constraints"
	"github.com/enetx/g/pkg/rand"
)

// NewInt creates a new Int with the provided int value.
func NewInt[T constraints.Integer | rune | byte](i T) Int { return Int(i) }

// Transform applies a transformation function to the Int and returns the result.
func (i Int) Transform(fn func(Int) Int) Int { return fn(i) }

// Min returns the minimum of Ints.
func (i Int) Min(b ...Int) Int { return cmp.Min(append(b, i)...) }

// Max returns the maximum of Ints.
func (i Int) Max(b ...Int) Int { return cmp.Max(append(b, i)...) }

// RandomRange returns a random Int in the range [from, to].
func (i Int) RandomRange(to Int) Int { return rand.N(to-i+1) + i }

// Bytes returns the Int as a byte slice.
func (i Int) Bytes() Bytes {
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, i.UInt64())

	return buffer[bits.LeadingZeros64(i.UInt64())>>3:]
}

// Abs returns the absolute value of the Int.
func (i Int) Abs() Int { return i.Float().Abs().Int() }

// Add adds two Ints and returns the result.
func (i Int) Add(b Int) Int { return i + b }

// BigInt returns the Int as a *big.Int.
func (i Int) BigInt() *big.Int { return big.NewInt(i.Int64()) }

// Div divides two Ints and returns the result.
func (i Int) Div(b Int) Int { return i / b }

// Eq checks if two Ints are equal.
func (i Int) Eq(b Int) bool { return i == b }

// Gt checks if the Int is greater than the specified Int.
func (i Int) Gt(b Int) bool { return i > b }

// Gte checks if the Int is greater than or equal to the specified Int.
func (i Int) Gte(b Int) bool { return i >= b }

// Float returns the Int as an Float.
func (i Int) Float() Float { return Float(i) }

// String returns the Int as an String.
func (i Int) String() String { return String(strconv.Itoa(int(i))) }

// StringBuf converts the Int to a String using the provided buffer without extra allocations when possible.
// If the buffer is too small, it will be automatically resized to fit the value.
// For values between 0 and 9, a single byte is written directly for maximum performance.
//
// Example:
//
//	buf := NewBytes()
//	name := Int(42).StringBuf(&buf)
//	fmt.Println(name)
//
// Note:
// The returned String shares the underlying memory with the buffer.
// Do not modify the buffer after calling this method unless the String is no longer needed.
func (i Int) StringBuf(buf *Bytes) String {
	if cap(*buf) < 20 {
		*buf = NewBytes(0, 20)
	} else {
		*buf = (*buf)[:0]
	}

	if i >= 0 && i <= 9 {
		*buf = append(*buf, '0'+byte(i))
	} else {
		*buf = strconv.AppendInt(*buf, int64(i), 10)
	}

	return buf.StringUnsafe()
}

// Std returns the Int as an int.
func (i Int) Std() int { return int(i) }

// Cmp compares two Ints and returns an cmp.Ordering.
func (i Int) Cmp(b Int) cmp.Ordering { return cmp.Cmp(i, b) }

// Int16 returns the Int as an int16.
func (i Int) Int16() int16 { return int16(i) }

// Int32 returns the Int as an int32.
func (i Int) Int32() int32 { return int32(i) }

// Int64 returns the Int as an int64.
func (i Int) Int64() int64 { return int64(i) }

// Int8 returns the Int as an int8.
func (i Int) Int8() int8 { return int8(i) }

// IsZero checks if the Int is 0.
func (i Int) IsZero() bool { return i.Eq(0) }

// IsNegative checks if the Int is negative.
func (i Int) IsNegative() bool { return i.Lt(0) }

// IsPositive checks if the Int is positive.
func (i Int) IsPositive() bool { return i.Gte(0) }

// Lt checks if the Int is less than the specified Int.
func (i Int) Lt(b Int) bool { return i < b }

// Lte checks if the Int is less than or equal to the specified Int.
func (i Int) Lte(b Int) bool { return i <= b }

// Mul multiplies two Ints and returns the result.
func (i Int) Mul(b Int) Int { return i * b }

// Ne checks if two Ints are not equal.
func (i Int) Ne(b Int) bool { return i != b }

// Random returns a random Int in the range [0, hi].
func (i Int) Random() Int { return Int(0).RandomRange(i) }

// Rem returns the remainder of the division between the receiver and the input value.
func (i Int) Rem(b Int) Int { return i % b }

// Sub subtracts two Ints and returns the result.
func (i Int) Sub(b Int) Int { return i - b }

// Binary returns the Int as a binary string.
func (i Int) Binary() String { return String(fmt.Sprintf("%08b", i)) }

// Hex returns the Int as a hexadecimal string.
func (i Int) Hex() String { return String(fmt.Sprintf("%x", i)) }

// Octal returns the Int as an octal string.
func (i Int) Octal() String { return String(fmt.Sprintf("%o", i)) }

// UInt returns the Int as a uint.
func (i Int) UInt() uint { return uint(i) }

// UInt16 returns the Int as a uint16.
func (i Int) UInt16() uint16 { return uint16(i) }

// UInt32 returns the Int as a uint32.
func (i Int) UInt32() uint32 { return uint32(i) }

// UInt64 returns the Int as a uint64.
func (i Int) UInt64() uint64 { return uint64(i) }

// UInt8 returns the Int as a uint8.
func (i Int) UInt8() uint8 { return uint8(i) }

// Print writes the value of the Int to the standard output (console)
// and returns the Int unchanged.
func (i Int) Print() Int { fmt.Print(i); return i }

// Println writes the value of the Int to the standard output (console) with a newline
// and returns the Int unchanged.
func (i Int) Println() Int { fmt.Println(i); return i }
