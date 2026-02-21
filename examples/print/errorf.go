package main

import (
	"errors"

	. "github.com/enetx/g"
)

var ErrNotFound = errors.New("not found")

func findUser(id int) error {
	// {1.wrap} formats the error into the message AND wraps it so
	// errors.Is / errors.As can inspect the chain.
	return Errorf("findUser({}): {.wrap}", id, ErrNotFound)
}

func loadProfile(id int) error {
	err := findUser(id)
	// Wrap the upstream error while adding more context.
	return Errorf("loadProfile: {.wrap}", err)
}

func main() {
	// --- Plain Errorf — no wrapping ---
	// Without {.wrap} the error chain is not preserved.
	plain := Errorf("something failed: {}", ErrNotFound)
	Println("message : {}", plain)
	Println("errors.Is: {}", errors.Is(plain, ErrNotFound)) // false

	Println("---")

	// --- Auto-index wrap: {.wrap} ---
	// The placeholder formats the value and simultaneously wraps the error.
	autoWrap := Errorf("step failed: {.wrap}", ErrNotFound)
	Println("message : {}", autoWrap)
	Println("errors.Is: {}", errors.Is(autoWrap, ErrNotFound)) // true

	Println("---")

	// --- Positional wrap: {1.wrap}, {2.wrap} ---
	// Useful when the error is not the only argument.
	e1 := errors.New("disk full")
	positional := Errorf("write {1}: {2.wrap}", "/var/log/app.log", e1)
	Println("message : {}", positional)
	Println("errors.Is: {}", errors.Is(positional, e1)) // true

	Println("---")

	// --- Named wrap: {cause.wrap} ---
	named := Errorf("open {file}: {cause.wrap}",
		Named{"file": "/etc/passwd", "cause": ErrNotFound})

	Println("message : {}", named)
	Println("errors.Is: {}", errors.Is(named, ErrNotFound)) // true

	Println("---")

	// --- Multiple wraps ---
	// Every {.wrap} / {N.wrap} in the template is collected.
	// errors.Is works for all of them (Go 1.20+ multi-unwrap).
	e2 := errors.New("timeout")
	multi := Errorf("{1.wrap} + {2.wrap}", e1, e2)
	Println("message  : {}", multi)
	Println("errors.Is e1: {}", errors.Is(multi, e1)) // true
	Println("errors.Is e2: {}", errors.Is(multi, e2)) // true

	Println("---")

	// --- Chained calls — wrap survives up the stack ---
	err := loadProfile(42)
	Println("top-level error   : {}", err)
	Println("errors.Is sentinel: {}", errors.Is(err, ErrNotFound)) // true
}
