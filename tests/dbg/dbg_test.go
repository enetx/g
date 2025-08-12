package g_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/enetx/g/dbg"
)

func TestCallerInfo(t *testing.T) {
	info := dbg.CallerInfo()

	// Should contain the format [filename:line] [function_name]
	if !strings.Contains(info, "dbg_test.go") {
		t.Errorf("CallerInfo() should contain current file name, got: %s", info)
	}

	if !strings.Contains(info, "TestCallerInfo") {
		t.Errorf("CallerInfo() should contain current function name, got: %s", info)
	}

	// Check format pattern [filename:line] [function]
	if !strings.HasPrefix(info, "[") || !strings.Contains(info, ":") {
		t.Errorf("CallerInfo() should have format [filename:line] [function], got: %s", info)
	}
}

func TestDbg(t *testing.T) {
	// Since Dbg prints to stdout/stderr, we need to capture output
	// This is a basic test to ensure the function doesn't panic

	// Backup original stdout
	oldStdout := os.Stdout

	// Create pipe to capture output
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	// Set stdout to write end of pipe
	os.Stdout = w

	// Call Dbg function
	testValue := 42
	dbg.Dbg(testValue)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	output := make([]byte, 1000)
	n, err := r.Read(output)
	if err != nil && n == 0 {
		t.Fatal("Failed to read captured output")
	}

	outputStr := string(output[:n])

	// Verify output contains expected elements
	if !strings.Contains(outputStr, "dbg_test.go") {
		t.Errorf("Dbg output should contain filename, got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "42") {
		t.Errorf("Dbg output should contain the value 42, got: %s", outputStr)
	}

	if !strings.Contains(outputStr, "testValue") {
		t.Errorf("Dbg output should contain variable name, got: %s", outputStr)
	}
}

func TestDbgWithError(t *testing.T) {
	// Test that error types are sent to stderr
	// This is more complex to test properly, so we'll just ensure no panic

	err := errors.New("test error")

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dbg with error should not panic: %v", r)
		}
	}()

	dbg.Dbg(err)
}

func TestDbgWithNilValue(t *testing.T) {
	// Test that nil values don't cause issues
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dbg with nil should not panic: %v", r)
		}
	}()

	var nilPtr *int
	dbg.Dbg(nilPtr)
}

func TestDbgWithComplexType(t *testing.T) {
	// Test with complex data structures
	testStruct := struct {
		Name string
		Age  int
	}{"Test", 25}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dbg with struct should not panic: %v", r)
		}
	}()

	dbg.Dbg(testStruct)
}

func TestDbgWithSlice(t *testing.T) {
	// Test with slice to hit more branches
	testSlice := []int{1, 2, 3, 4, 5}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dbg with slice should not panic: %v", r)
		}
	}()

	dbg.Dbg(testSlice)
}

func TestDbgWithMap(t *testing.T) {
	// Test with map
	testMap := map[string]int{"a": 1, "b": 2}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dbg with map should not panic: %v", r)
		}
	}()

	dbg.Dbg(testMap)
}

func TestDbgMultipleValues(t *testing.T) {
	// Test multiple calls to ensure different code paths are hit
	values := []any{
		"string value",
		123,
		3.14,
		true,
		[]byte("bytes"),
		struct{ X int }{42},
	}

	for _, val := range values {
		func(v any) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Dbg with value %v should not panic: %v", v, r)
				}
			}()
			dbg.Dbg(v)
		}(val)
	}
}
