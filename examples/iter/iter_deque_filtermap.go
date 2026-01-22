package main

import (
	"strconv"

	. "github.com/enetx/g"
)

func main() {
	// Example 1: FilterMap to parse and validate numbers - first filter valid strings
	inputs := DequeOf("1", "abc", "42", "", "7", "xyz")
	validStrings := inputs.Iter().FilterMap(func(s string) Option[string] {
		if _, err := strconv.Atoi(s); err == nil {
			return Some(s) // Keep valid numeric strings
		}
		return None[string]()
	}).Collect()

	validStrings.Println() // Deque[1, 42, 7]

	// Example 2: FilterMap to process and filter scores
	scores := DequeOf(85, 45, 92, 78, 33, 98)
	highGrades := scores.Iter().FilterMap(func(score int) Option[int] {
		if score >= 80 {
			return Some(score * 2) // Double high scores
		}
		return None[int]() // Filter out scores below 80
	}).Collect()

	highGrades.Println() // Show doubled high scores

	// Example 3: FilterMap to extract file extensions
	filenames := DequeOf("document.pdf", "image", "script.py", "data.json", "readme")
	extensions := filenames.Iter().FilterMap(func(filename string) Option[string] {
		if parts := String(filename).Split(".").Collect(); parts.Len() > 1 {
			ext := parts[parts.Len()-1]
			return Some("." + string(ext))
		}
		return None[string]() // Filter out files without extensions
	}).Collect()

	extensions.Println() // Deque[.pdf, .py, .json]

	// Example 4: FilterMap to process large numbers
	numbers := DequeOf(10, 8, 15, 7, 20)

	// Create safe division results
	safeResults := numbers.Iter().FilterMap(func(num int) Option[int] {
		if num > 10 { // Only process numbers > 10
			return Some(num * 2) // Double them
		}
		return None[int]()
	}).Collect()

	safeResults.Println() // Deque[30, 40]

	// Example 5: FilterMap to clean and validate user input
	userInputs := DequeOf("  john@email.com  ", "", "invalid-email", "  jane@test.org", "   ", "admin@site.net  ")
	validEmails := userInputs.Iter().FilterMap(func(input string) Option[string] {
		cleaned := String(input).Trim()
		// Simple email validation (contains @ and .)
		if !cleaned.IsEmpty() && cleaned.Contains("@") && cleaned.Contains(".") {
			return Some(string(cleaned.Lower()))
		}
		return None[string]()
	}).Collect()

	validEmails.Println() // Deque[john@email.com, jane@test.org, admin@site.net]

	// Example 6: FilterMap for string processing
	words := DequeOf("hello", "", "world", "   ", "go", "rust")
	validWords := words.Iter().FilterMap(func(s string) Option[string] {
		trimmed := String(s).Trim()
		if !trimmed.IsEmpty() && len(trimmed) > 2 {
			return Some(string(trimmed.Upper()))
		}
		return None[string]()
	}).Collect()

	validWords.Println() // Deque[HELLO, WORLD, RUST]
}
