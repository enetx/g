package main

import (
	"os"

	. "github.com/enetx/g"
)

func main() {
	// Print writes a formatted string to stdout.
	// Placeholders use {} syntax, not %v.
	Print("Hello, {}!\n", "world")

	// Println appends a newline automatically — no need for \n.
	Println("Hello, {}", "world")

	// Eprint and Eprintln write to stderr instead of stdout.
	Eprint("warning: {} not found\n", "config.toml")
	Eprintln("error: {}", "connection refused")

	// All print functions return Result[int] — bytes written or an error.
	// Typically you can ignore it, but it is there when you need it.
	res := Println("written: {}", 42)
	if res.IsErr() {
		Eprintln("write failed: {}", res.Err())
	}

	// Write targets any io.Writer — useful for buffers, files, or pipes.
	var buf Builder
	Write(&buf, "buffered: {} + {}", "foo", "bar")
	Println("buffer contains: {}", buf.String())

	// Writeln does the same but appends a newline.
	Write(os.Stdout, "direct to stdout: {}\n", "ok")
}
