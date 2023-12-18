package g

import (
	"bufio"
	"os"
)

type (
	// Result is a generic struct for representing a result value along with an error.
	Result[T any] struct {
		value *T    // Pointer to the value.
		err   error // Associated error.
	}

	// Option is a generic struct for representing an optional value.
	Option[T any] struct {
		value *T // Pointer to the value.
	}

	// fiter is a struct for iterating through an file.
	fiter struct {
		scanner *bufio.Scanner // Scanner for reading from the file.
		file    *File          // Associated File.
	}

	// File is a struct that represents a file along with an iterator for reading lines.
	File struct {
		file  *os.File // Underlying os.File.
		fiter *fiter   // Iterator for reading lines.
		name  String   // File name.
		guard bool     // Guard indicates whether the file is protected against concurrent access.
	}

	// Dir is a struct representing a directory path.
	Dir struct {
		path String // Directory path.
	}

	// String is an alias for the string type.
	String string

	// Int is an alias for the int type.
	Int int

	// Float is an alias for the float64 type.
	Float float64

	// Bytes is an alias for the []byte type.
	Bytes []byte

	// Slice is a generic alias for a slice.
	Slice[T any] []T

	// Map is a generic alias for a map.
	Map[K comparable, V any] map[K]V

	// Set is a generic alias for a set implemented using a map.
	Set[T comparable] map[T]struct{}

	// mapPair is a struct representing a key-value pair for MapOrd.
	mapPair[K comparable, V any] struct {
		Key   K // Key of the pair.
		Value V // Value associated with the key.
	}

	// MapOrd is a generic alias for a slice of ordered key-value pairs.
	MapOrd[K comparable, V any] Slice[mapPair[K, V]]
)