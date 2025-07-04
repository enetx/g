package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	s := String("abc123")
	b := s.BytesUnsafe()
	fmt.Printf("Original string: %q\n", s)
	fmt.Printf("Zero-copy bytes: %q\n", b)

	b2 := Bytes("xyz789")
	s2 := b2.StringUnsafe()
	fmt.Printf("Original bytes: %q\n", b2)
	fmt.Printf("Zero-copy string: %q\n", s2)

	b2[0] = 'X'
	fmt.Printf("Modified bytes: %q\n", b2)
	fmt.Printf("Affected string: %q\n", s2)

	// Original string: "abc123"
	// Zero-copy bytes: "abc123"
	// Original bytes: "xyz789"
	// Zero-copy string: "xyz789"
	// Modified bytes: "Xyz789"
	// Affected string: "Xyz789"
}
