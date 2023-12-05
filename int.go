package g

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"math/bits"
	"strconv"

	"gitlab.com/x0xO/g/pkg/minmax"
	"gitlab.com/x0xO/g/pkg/rand"
)

// NewInt creates a new Int with the provided int value.
func NewInt[T ~int | rune | byte](i T) Int { return Int(i) }

// Min returns the minimum of Ints.
func (i Int) Min(b ...Int) Int { return minmax.Min(i, b...) }

// Max returns the maximum of Ints.
func (i Int) Max(b ...Int) Int { return minmax.Max(i, b...) }

// RandomRange returns a random Int in the range [from, to].
func (Int) RandomRange(from, to Int) Int {
	return Int(rand.Intn(to.Sub(from).Add(1).Std())).Add(from)
}

// Bytes returns the Int as a byte slice.
func (i Int) Bytes() Bytes {
	buffer := make([]byte, 8)
	binary.BigEndian.PutUint64(buffer, i.AsUInt64())

	return buffer[bits.LeadingZeros64(i.AsUInt64())>>3:]
}

// Abs returns the absolute value of the Int.
func (i Int) Abs() Int { return i.ToFloat().Abs().ToInt() }

// Add adds two Ints and returns the result.
func (i Int) Add(b Int) Int { return i + b }

// ToBigInt returns the Int as a *big.Int.
func (i Int) ToBigInt() *big.Int { return big.NewInt(i.AsInt64()) }

// Div divides two Ints and returns the result.
func (i Int) Div(b Int) Int { return i / b }

// Eq checks if two Ints are equal.
func (i Int) Eq(b Int) bool { return i == b }

// Gt checks if the Int is greater than the specified Int.
func (i Int) Gt(b Int) bool { return i > b }

// Gte checks if the Int is greater than or equal to the specified Int.
func (i Int) Gte(b Int) bool { return i >= b }

// ToFloat returns the Int as an Float.
func (i Int) ToFloat() Float { return Float(i) }

// ToString returns the Int as an String.
func (i Int) ToString() String { return String(strconv.Itoa(int(i))) }

// Std returns the Int as an int.
func (i Int) Std() int { return int(i) }

// AsInt16 returns the Int as an int16.
func (i Int) AsInt16() int16 { return int16(i) }

// AsInt32 returns the Int as an int32.
func (i Int) AsInt32() int32 { return int32(i) }

// AsInt64 returns the Int as an int64.
func (i Int) AsInt64() int64 { return int64(i) }

// AsInt8 returns the Int as an int8.
func (i Int) AsInt8() int8 { return int8(i) }

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
func (i Int) Random() Int { return i.RandomRange(0, i) }

// Rem returns the remainder of the division between the receiver and the input value.
func (i Int) Rem(b Int) Int { return i % b }

// Sub subtracts two Ints and returns the result.
func (i Int) Sub(b Int) Int { return i - b }

// ToBinary returns the Int as a binary string.
func (i Int) ToBinary() String { return String(fmt.Sprintf("%08b", i)) }

// ToHex returns the Int as a hexadecimal string.
func (i Int) ToHex() String { return String(fmt.Sprintf("%x", i)) }

// ToOctal returns the Int as an octal string.
func (i Int) ToOctal() String { return String(fmt.Sprintf("%o", i)) }

// AsUInt returns the Int as a uint.
func (i Int) AsUInt() uint { return uint(i) }

// AsUInt16 returns the Int as a uint16.
func (i Int) AsUInt16() uint16 { return uint16(i) }

// AsUInt32 returns the Int as a uint32.
func (i Int) AsUInt32() uint32 { return uint32(i) }

// AsUInt64 returns the Int as a uint64.
func (i Int) AsUInt64() uint64 { return uint64(i) }

// AsUInt8 returns the Int as a uint8.
func (i Int) AsUInt8() uint8 { return uint8(i) }

// Print prints the value of the Int to the standard output (console)
// and returns the Int unchanged.
func (i Int) Print() Int { fmt.Println(i); return i }
