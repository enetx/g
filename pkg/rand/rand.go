package rand

import (
	"crypto/rand"
	"encoding/binary"
)

// Intn generates a random non-negative integer within the range [0, max).
// The generated integer will be less than the provided maximum value.
// If max is less than or equal to 0, the function will treat it as if max is 1.
//
// Usage:
//
//	max := 10
//	randomInt := Intn(max)
//	fmt.Printf("Random integer between 0 and %d: %d\n", max, randomInt)
//
// Parameters:
//   - max (int): The maximum bound for the random integer to be generated.
//
// Returns:
//   - int: A random non-negative integer within the specified range.
func Intn(max int) int {
	// If the provided maximum value is less than or equal to 0,
	// set it to 1 to ensure a valid range.
	if max <= 0 {
		max = 1
	}

	// Declare a uint64 variable to store the random value.
	var randVal uint64

	// Read a random value from the rand.Reader (a cryptographically
	// secure random number generator) into the randVal variable,
	// using binary.LittleEndian for byte ordering.
	_ = binary.Read(rand.Reader, binary.LittleEndian, &randVal)

	// Return the generated random value modulo the maximum value as an integer.
	return int(randVal % uint64(max))
}
