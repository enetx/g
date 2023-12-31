package rand

import (
	"crypto/rand"
	"encoding/binary"
)

type intType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// N generates a random non-negative integer within the range [0, max).
// The generated integer will be less than the provided maximum value.
// If max is less than or equal to 0, the function will treat it as if max is 1.
//
// Usage:
//
//	n := 10
//	randomInt := rand.N(n)
//	fmt.Printf("Random integer between 0 and %d: %d\n", max, randomInt)
//
// Parameters:
//   - n (int): The maximum bound for the random integer to be generated.
//
// Returns:
//   - int: A random non-negative integer within the specified range.
func N[T intType](n T) T {
	// If the provided maximum value is less than or equal to 0,
	// set it to 1 to ensure a valid range.
	if n <= 0 {
		n = 1
	}

	// Declare a uint64 variable to store the random value.
	var randVal uint64

	// Read a random value from the rand.Reader (a cryptographically
	// secure random number generator) into the randVal variable,
	// using binary.LittleEndian for byte ordering.
	_ = binary.Read(rand.Reader, binary.LittleEndian, &randVal)

	// Return the generated random value modulo the maximum value as an integer.
	return T(randVal % uint64(n))
}
