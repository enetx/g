// Package ref provides a utility function for creating a pointer to a value.
// It is designed to simplify the process of obtaining a pointer to a value of any type.
package ref

// Of creates a pointer to the provided value of type 'E'.
// The primary purpose of this function is to simplify the creation of pointers to values
// without needing to use temporary variables.
//
// Copy hazard: the argument is passed by value, so Of returns a pointer to a fresh copy,
// not to the caller's original variable. Do not use Of with types that must not be copied
// after first use (e.g. sync.Mutex, sync.WaitGroup, strings.Builder, or any struct embedding
// such a type, including g's *Mutex/*RwLock/*Builder) — the copy may carry stale internal
// state and `go vet` will not flag it here. Take the address of the original value instead.
//
// Parameters:
//
//	e: The value of type 'E' to create a pointer for.
//
// Returns:
//
//	*E: A pointer to a copy of the provided value 'e'.
//
// Example usage:
//
//	intValue := 42
//	intPtr := ref.Of(intValue)
//	fmt.Println(*intPtr)
func Of[E any](e E) *E { return &e }
