package main

import (
	. "github.com/enetx/g"
)

func main() {
	// Example 1: FilterMap to parse strings to valid strings only
	SliceOf[String]("1", "2", "abc", "3", "xyz", "4").
		Iter().
		FilterMap(func(s String) Option[String] {
			// Keep only numeric strings
			if s.IsDigit() {
				return Some(s + "_valid")
			}

			return None[String]()
		}).
		Collect().
		Println() // Slice[1_valid, 2_valid, 3_valid, 4_valid]

	// Example 2: FilterMap to get positive doubled values
	SliceOf(1, -2, 3, -4, 5, 0).
		Iter().
		FilterMap(func(n int) Option[int] {
			// Only keep positive numbers and double them
			if n > 0 {
				return Some(n * 2)
			}
			return None[int]()
		}).
		Collect().
		Println() // Slice[2, 6, 10]

	// Example 3: FilterMap to extract valid email domains
	SliceOf[String]("user@example.com", "invalid-email", "admin@test.org", "no-at-sign").
		Iter().
		FilterMap(func(email String) Option[String] {
			// Extract domain if email contains @
			parts := email.Split("@").Collect()
			if parts.Len() == 2 {
				return Some(parts[1])
			}
			return None[String]()
		}).
		Collect().
		Println() // Slice[example.com, test.org]

	// Example 4: FilterMap to calculate safe division
	SliceOf(
		Pair[int, int]{Key: 10, Value: 2},
		Pair[int, int]{Key: 8, Value: 0},
		Pair[int, int]{Key: 12, Value: 3},
		Pair[int, int]{Key: 5, Value: 0},
		Pair[int, int]{Key: 20, Value: 4},
	).
		Iter().
		FilterMap(func(p Pair[int, int]) Option[Pair[int, int]] {
			// Only return result if division is safe (no division by zero)
			if p.Value != 0 {
				return Some(Pair[int, int]{Key: p.Key / p.Value, Value: 0}) // Store result in a field
			}
			return None[Pair[int, int]]()
		}).
		Collect().
		Println()
		// Results for safe divisions

	// Example 5: FilterMap with string processing
	SliceOf[String]("  hello  ", "", "world", "   ", "go").
		Iter().
		FilterMap(func(s String) Option[String] {
			// Trim and uppercase, filter out empty strings
			trimmed := s.Trim()
			if !trimmed.Empty() {
				return Some(trimmed.Upper())
			}
			return None[String]()
		}).
		Collect().
		Println() // Slice[HELLO, WORLD, GO]
}
