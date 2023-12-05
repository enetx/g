// Package deref provides a utility function to dereference a pointer.
package deref

// Of takes a pointer to a value of any type E and returns the value it points to.
// This function can be useful for simplifying code when working with pointers and their values.
//
// Parameters:
//
//	e: A pointer to a value of any type E.
//
// Returns:
//
//	E: The value pointed to by the input pointer.
//
// Example usage:
//
//	pi := 3.141592
//	ptr := &pi
//	value := deref.Of(ptr) // value is 3.141592 (type float64)
func Of[E any](e *E) E { return *e }
