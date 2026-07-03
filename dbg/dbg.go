// Package dbg provides debugging helpers that print an expression together with
// its source file, line and enclosing function (dbg.Dbg, dbg.CallerInfo).
package dbg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// CallerInfo retrieves and formats the file name and line number of the caller's
// location in the source code.
//
// It uses the runtime.Caller function to determine the caller's file and line
// information.
//
// Example usage:
//
//	callerInfo := CallerInfo()
//	fmt.Println("Caller:", callerInfo)
//
// The output will look like this:
//
//	Caller: [filename:line]
//
// Where "filename" is the name of the file where the caller is located, and
// "line" is the line number in that file.
func CallerInfo() string {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Fprintln(os.Stderr, "dbg.CallerInfo: Unable to parse runtime caller")
		return ""
	}

	out := fmt.Sprintf("[%s:%d] [%s]", filepath.Base(file), line, funcName(pc))

	return out
}

// funcName resolves the function name for a program counter, guarding against a
// nil *runtime.Func (FuncForPC may return nil for an unknown PC).
func funcName(pc uintptr) string {
	if fn := runtime.FuncForPC(pc); fn != nil {
		return fn.Name()
	}

	return "?"
}

// Dbg is a debugging utility function that prints the given expression and its
// location in the source code to standard output. It also prints the file name
// and line number where the Dbg function is called.
//
// Parameters:
//
//	exp  - The expression or value to be printed.
//
// Example:
//
//	Dbg(someValue)
//
// The output will look like this:
//
//	[filename:line] variable = value
//
// Where "filename" is the name of the file where Dbg is called, "line" is the
// line number in that file, "variable" is the variable or expression being
// debugged, and "value" is the value of the expression.
func Dbg(exp any) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Fprintln(os.Stderr, "dbg.Dbg: Unable to parse runtime caller")
		return
	}

	base := filepath.Base(file)
	fn := funcName(pc)

	// Default output for when the source line is unavailable: still print the
	// location and value rather than dropping the call entirely.
	out := fmt.Sprintf("[%s:%d] [%s] = %+v", base, line, fn, exp)

	if f, err := os.Open(file); err == nil {
		defer f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)

		for i := 1; scanner.Scan(); i++ {
			if i == line {
				if v, ok := extractExpr(scanner.Text()); ok {
					out = fmt.Sprintf("[%s:%d] [%s] %s = %+v", base, line, fn, v, exp)
				}
				break
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return
		}
	}

	switch exp.(type) {
	case error:
		fmt.Fprintln(os.Stderr, out)
	default:
		fmt.Fprintln(os.Stdout, out)
	}
}

// extractExpr pulls the argument expression out of a `Dbg(...)` source line by
// slicing between the first '(' and the last ')'. Multi-line calls or lines
// without a balanced pair of parentheses cannot be parsed, so it reports
// ok=false and the caller falls back to a source-less rendering instead of
// panicking on an out-of-bounds slice.
func extractExpr(text string) (string, bool) {
	open := strings.Index(text, "(")
	last := strings.LastIndex(text, ")")

	if open < 0 || last < 0 || last <= open+1 {
		return "", false
	}

	return text[open+1 : last], true
}
