package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/enetx/g"
)

var errNotFound = errors.New("not found")

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

func main() {
	// ErrIs - check for specific error
	exampleErrIs()

	// ErrAs - cast error to type
	exampleErrAs()

	// ErrSource - get wrapped error
	exampleErrSource()

	// Wrap - wrap error with context
	exampleWrap()
}

func exampleErrIs() {
	fmt.Println("=== ErrIs ===")

	// Check for sentinel error
	res := g.Err[string](errNotFound)
	if res.ErrIs(errNotFound) {
		fmt.Println("error: not found")
	}

	// Works with wrapped errors
	wrapped := g.Err[string](fmt.Errorf("failed to load config: %w", errNotFound))
	if wrapped.ErrIs(errNotFound) {
		fmt.Println("wrapped error also matched")
	}

	// Check for os.ErrNotExist
	fileRes := g.ResultOf(os.Open("nonexistent.txt"))
	if fileRes.ErrIs(os.ErrNotExist) {
		fmt.Println("file does not exist")
	}

	// Ok result - always false
	ok := g.Ok(42)
	fmt.Printf("Ok.ErrIs(errNotFound) = %v\n", ok.ErrIs(errNotFound))

	fmt.Println()
}

func exampleErrAs() {
	fmt.Println("=== ErrAs ===")

	// Cast to custom error type
	valErr := &ValidationError{Field: "email", Message: "invalid format"}
	res := g.Err[string](valErr)

	var target *ValidationError
	if res.ErrAs(&target) {
		fmt.Printf("field: %s, message: %s\n", target.Field, target.Message)
	}

	// Works with wrapped errors
	wrapped := g.Err[string](fmt.Errorf("request failed: %w", valErr))
	var wrappedTarget *ValidationError
	if wrapped.ErrAs(&wrappedTarget) {
		fmt.Printf("from wrapped: field: %s\n", wrappedTarget.Field)
	}

	// Cast to *os.PathError
	fileRes := g.ResultOf(os.Open("nonexistent.txt"))
	var pathErr *os.PathError
	if fileRes.ErrAs(&pathErr) {
		fmt.Printf("path: %s, op: %s\n", pathErr.Path, pathErr.Op)
	}

	fmt.Println()
}

func exampleErrSource() {
	fmt.Println("=== ErrSource ===")

	// Get inner error
	inner := errors.New("connection refused")
	outer := fmt.Errorf("failed to connect: %w", inner)
	res := g.Err[string](outer)

	if source := res.ErrSource(); source.IsSome() {
		fmt.Printf("source: %s\n", source.Some())
	}

	// Error without wrapper - None
	plain := g.Err[string](errors.New("simple error"))
	if plain.ErrSource().IsNone() {
		fmt.Println("plain error has no source")
	}

	// Ok result - None
	ok := g.Ok(42)
	if ok.ErrSource().IsNone() {
		fmt.Println("Ok result - no source")
	}

	// Unwrap chain
	fmt.Println("\nError chain:")
	level1 := errors.New("root cause")
	level2 := fmt.Errorf("level 2: %w", level1)
	level3 := fmt.Errorf("level 3: %w", level2)

	r := g.Err[int](level3)
	for {
		fmt.Printf("  -> %s\n", r.Err())
		source := r.ErrSource()
		if !source.IsSome() {
			fmt.Println("  (root)")
			break
		}

		r = g.Err[int](source.Some())
	}
}

var (
	errLoadConfig = errors.New("failed to load config")
	errStep1      = errors.New("step 1")
	errStep2      = errors.New("step 2")
	errStep3      = errors.New("step 3")
	errIgnored    = errors.New("ignored")
)

func exampleWrap() {
	fmt.Println()
	fmt.Println("=== Wrap ===")

	// Wrap error with context error
	res := g.ResultOf(os.Open("nonexistent.txt")).Wrap(errLoadConfig)
	fmt.Printf("wrapped: %s\n", res.Err())

	// Both errors are preserved in chain
	if res.ErrIs(os.ErrNotExist) {
		fmt.Println("original error accessible via ErrIs")
	}
	if res.ErrIs(errLoadConfig) {
		fmt.Println("wrapper error accessible via ErrIs")
	}

	// Chain multiple wraps
	chained := g.Err[int](errNotFound).
		Wrap(errStep1).
		Wrap(errStep2).
		Wrap(errStep3)
	fmt.Printf("chained: %s\n", chained.Err())

	// All errors in chain are accessible
	fmt.Printf("ErrIs(errNotFound): %v\n", chained.ErrIs(errNotFound))
	fmt.Printf("ErrIs(errStep1): %v\n", chained.ErrIs(errStep1))
	fmt.Printf("ErrIs(errStep3): %v\n", chained.ErrIs(errStep3))

	// Ok result - unchanged
	ok := g.Ok(42).Wrap(errIgnored)
	fmt.Printf("Ok.Wrap() = %v\n", ok.Unwrap())
}
