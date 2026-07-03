package main

import (
	. "github.com/enetx/g"
)

func main() {
	// Example conversions showing round-trip behavior
	Println("=== BigEndian Examples ===")

	// Positive number
	val1 := Int(1000)
	bytes1 := val1.BytesBE()
	back1 := bytes1.IntBE()
	Println("Val: {} -> Bytes BE: {} -> Back: {}", val1, bytes1, back1)

	// Negative number
	val2 := Int(-1000)
	bytes2 := val2.BytesBE()
	back2 := bytes2.IntBE()
	Println("Val: {} -> Bytes BE: {} -> Back: {}", val2, bytes2, back2)

	// Small positive number
	val3 := Int(255)
	bytes3 := val3.BytesBE()
	back3 := bytes3.IntBE()
	Println("Val: {} -> Bytes BE: {} -> Back: {}", val3, bytes3, back3)

	// Large number
	val4 := Int(0x123456789ABCDEF0)
	bytes4 := val4.BytesBE()
	back4 := bytes4.IntBE()
	Println("Val: {} -> Bytes BE: {} -> Back: {}", val4, bytes4, back4)

	Println("\n=== LittleEndian Examples ===")

	// Positive number
	val5 := Int(1000)
	bytes5 := val5.BytesLE()
	back5 := bytes5.IntLE()
	Println("Val: {} -> Bytes LE: {} -> Back: {}", val5, bytes5, back5)

	// Negative number
	val6 := Int(-1000)
	bytes6 := val6.BytesLE()
	back6 := bytes6.IntLE()
	Println("Val: {} -> Bytes LE: {} -> Back: {}", val6, bytes6, back6)

	// Small positive number
	val7 := Int(255)
	bytes7 := val7.BytesLE()
	back7 := bytes7.IntLE()
	Println("Val: {} -> Bytes LE: {} -> Back: {}", val7, bytes7, back7)

	// Large number
	val8 := Int(0x123456789ABCDEF0)
	bytes8 := val8.BytesLE()
	back8 := bytes8.IntLE()
	Println("Val: {} -> Bytes LE: {} -> Back: {}", val8, bytes8, back8)

	Println("\n=== Edge Cases ===")

	// Zero
	val9 := Int(0)
	bytes9BE := val9.BytesBE()
	bytes9LE := val9.BytesLE()
	Println("Zero -> BE: {}, LE: {}", bytes9BE, bytes9LE)

	// -1 (all bits set)
	val10 := Int(-1)
	bytes10BE := val10.BytesBE()
	bytes10LE := val10.BytesLE()
	Println("-1 -> BE: {}, LE: {}", bytes10BE, bytes10LE)

	// Max int64
	val11 := Int(9223372036854775807) // 0x7FFFFFFFFFFFFFFF
	bytes11BE := val11.BytesBE()
	bytes11LE := val11.BytesLE()
	Println("MaxInt64 -> BE: {}, LE: {}", bytes11BE, bytes11LE)

	// Min int64
	val12 := Int(-9223372036854775808) // 0x8000000000000000
	bytes12BE := val12.BytesBE()
	bytes12LE := val12.BytesLE()
	Println("MinInt64 -> BE: {}, LE: {}", bytes12BE, bytes12LE)
}
