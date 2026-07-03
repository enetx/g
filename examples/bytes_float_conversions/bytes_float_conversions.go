package main

import (
	"fmt"
	"math"

	. "github.com/enetx/g"
)

func main() {
	// Example conversions showing round-trip behavior for Float64
	fmt.Println("=== Float BigEndian Examples ===")

	// Positive number
	val1 := Float(3.14159)
	bytes1 := val1.BytesBE()
	back1 := bytes1.FloatBE()
	Println("Val: {} -> Bytes BE: {} -> Back: {}", val1, bytes1, back1)

	// Negative number
	val2 := Float(-123.456)
	bytes2 := val2.BytesBE()
	back2 := bytes2.FloatBE()
	Println("Val: {} -> Bytes BE: {} -> Back: {}", val2, bytes2, back2)

	// Very small number
	val3 := Float(1e-10)
	bytes3 := val3.BytesBE()
	back3 := bytes3.FloatBE()
	Println("Val: {} -> Bytes BE: {} -> Back: {}", val3, bytes3, back3)

	// Very large number
	val4 := Float(1.7976931348623157e+308) // Near max float64
	bytes4 := val4.BytesBE()
	back4 := bytes4.FloatBE()
	Println("Val: {} -> Bytes BE: {} -> Back: {}", val4, bytes4, back4)

	fmt.Println("\n=== Float LittleEndian Examples ===")

	// Positive number
	val5 := Float(3.14159)
	bytes5 := val5.BytesLE()
	back5 := bytes5.FloatLE()
	Println("Val: {} -> Bytes LE: {} -> Back: {}", val5, bytes5, back5)

	// Negative number
	val6 := Float(-123.456)
	bytes6 := val6.BytesLE()
	back6 := bytes6.FloatLE()
	Println("Val: {} -> Bytes LE: {} -> Back: {}", val6, bytes6, back6)

	// Very small number
	val7 := Float(1e-10)
	bytes7 := val7.BytesLE()
	back7 := bytes7.FloatLE()
	Println("Val: {} -> Bytes LE: {} -> Back: {}", val7, bytes7, back7)

	// Very large number
	val8 := Float(1.7976931348623157e+308)
	bytes8 := val8.BytesLE()
	back8 := bytes8.FloatLE()
	Println("Val: {} -> Bytes LE: {} -> Back: {}", val8, bytes8, back8)

	fmt.Println("\n=== Special Values ===")

	// Zero
	val10 := Float(0.0)
	bytes10BE := val10.BytesBE()
	bytes10LE := val10.BytesLE()
	Println("Zero -> BE: {}, LE: {}", bytes10BE, bytes10LE)

	// Negative zero
	val11 := Float(math.Copysign(0, -1))
	bytes11BE := val11.BytesBE()
	bytes11LE := val11.BytesLE()
	Println("Negative Zero -> BE: {}, LE: {}", bytes11BE, bytes11LE)

	// Infinity
	val12 := Float(math.Inf(1))
	bytes12BE := val12.BytesBE()
	bytes12LE := val12.BytesLE()
	Println("Positive Infinity -> BE: {}, LE: {}", bytes12BE, bytes12LE)

	// Negative infinity
	val13 := Float(math.Inf(-1))
	bytes13BE := val13.BytesBE()
	bytes13LE := val13.BytesLE()
	Println("Negative Infinity -> BE: {}, LE: {}", bytes13BE, bytes13LE)

	// NaN
	val14 := Float(math.NaN())
	bytes14BE := val14.BytesBE()
	bytes14LE := val14.BytesLE()
	back14BE := bytes14BE.FloatBE()
	back14LE := bytes14LE.FloatLE()
	Println("NaN -> BE: {}, LE: {}\n", bytes14BE, bytes14LE)
	Println("NaN recovered -> BE: {} (isNaN: {}), LE: {} (isNaN: {})",
		back14BE, math.IsNaN(float64(back14BE)), back14LE, math.IsNaN(float64(back14LE)))

	fmt.Println("\n=== Precision Comparison Float64 vs Float32 ===")

	// Show precision difference
	preciseVal := Float(3.141592653589793238462643383279)

	// Float64 representation
	bytes64BE := preciseVal.BytesBE()
	back64BE := bytes64BE.FloatBE()

	Println("Original: {}", preciseVal)
	Println("Float64:  {} (bytes: {})", back64BE, bytes64BE)

	Println("\n=== Byte Order Comparison ===")

	// Show byte order difference for the same value
	testVal := Float(1.23456789)
	bytesBE := testVal.BytesBE()
	bytesLE := testVal.BytesLE()

	Println("Value: {}", testVal)
	Println("BigEndian:    {}", bytesBE)
	Println("LittleEndian: {}", bytesLE)
	Println("Bytes reversed: {}", bytesBE[0] == bytesLE[7] && bytesBE[7] == bytesLE[0])
}
