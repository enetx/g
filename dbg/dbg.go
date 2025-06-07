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

	out := fmt.Sprintf("[%s:%d] [%s]", filepath.Base(file), line, runtime.FuncForPC(pc).Name())

	return out
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

	f, err := os.Open(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dbg.Dbg: Unable to open expected file")
		return
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var out string
	for i := 1; scanner.Scan(); i++ {
		if i == line {
			v := scanner.Text()[strings.Index(scanner.Text(), "(")+1 : strings.LastIndex(scanner.Text(), ")")]
			out = fmt.Sprintf("[%s:%d] [%s] %s = %+v", filepath.Base(file), line, runtime.FuncForPC(pc).Name(), v, exp)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return
	}

	switch exp.(type) {
	case error:
		fmt.Fprintln(os.Stderr, out)
	default:
		fmt.Fprintln(os.Stdout, out)
	}
}
