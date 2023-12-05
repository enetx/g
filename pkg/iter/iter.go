// Package iter provides a utility function for creating a range of integers.
// It is designed to simplify iteration over a range of integers using for-range loops.
package iter

// N creates a slice of length 'n' containing empty struct{} elements.
// The primary purpose of this function is to simplify the creation of for-range loops
// with a specific number of iterations.
//
// Parameters:
//
//	n: The length of the range to create, represented by the length of the returned slice.
//
// Returns:
//
//	[]struct{}: A slice of length 'n' containing empty struct{} elements.
//
// Example usage:
//
//	for i := range iter.N(5) {
//	    fmt.Println(i)
//	}
func N(n int) []struct{} { return make([]struct{}, n) }
