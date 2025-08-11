package g_test

import (
	"testing"

	"github.com/enetx/g/ref"
)

func TestOf_Int(t *testing.T) {
	value := 42
	ptr := ref.Of(value)

	// Check that we got a valid pointer
	if ptr == nil {
		t.Error("Of() returned nil pointer")
	}

	// Check that the pointer points to the correct value
	if *ptr != value {
		t.Errorf("Of(%v) = pointer to %v, want pointer to %v", value, *ptr, value)
	}

	// Check that modifying the original doesn't affect the pointed value
	// (because Of creates a copy)
	originalValue := value
	value = 100
	if *ptr != originalValue {
		t.Errorf("Pointer value changed when original was modified: got %v, want %v", *ptr, originalValue)
	}
}

func TestOf_String(t *testing.T) {
	value := "hello world"
	ptr := ref.Of(value)

	if ptr == nil {
		t.Error("Of() returned nil pointer")
	}

	if *ptr != value {
		t.Errorf("Of(%q) = pointer to %q, want pointer to %q", value, *ptr, value)
	}
}

func TestOf_Struct(t *testing.T) {
	type TestStruct struct {
		Name string
		Age  int
	}

	value := TestStruct{Name: "Alice", Age: 30}
	ptr := ref.Of(value)

	if ptr == nil {
		t.Error("Of() returned nil pointer")
	}

	if *ptr != value {
		t.Errorf("Of(%+v) = pointer to %+v, want pointer to %+v", value, *ptr, value)
	}

	// Test field access
	if ptr.Name != "Alice" {
		t.Errorf("ptr.Name = %q, want %q", ptr.Name, "Alice")
	}

	if ptr.Age != 30 {
		t.Errorf("ptr.Age = %v, want %v", ptr.Age, 30)
	}
}

func TestOf_Slice(t *testing.T) {
	value := []int{1, 2, 3, 4, 5}
	ptr := ref.Of(value)

	if ptr == nil {
		t.Error("Of() returned nil pointer")
	}

	// Check slice equality
	if len(*ptr) != len(value) {
		t.Errorf("Slice lengths don't match: got %v, want %v", len(*ptr), len(value))
	}

	for i, v := range value {
		if (*ptr)[i] != v {
			t.Errorf("Slice element at index %d: got %v, want %v", i, (*ptr)[i], v)
		}
	}
}

func TestOf_Map(t *testing.T) {
	value := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	ptr := ref.Of(value)

	if ptr == nil {
		t.Error("Of() returned nil pointer")
	}

	// Check map equality
	if len(*ptr) != len(value) {
		t.Errorf("Map lengths don't match: got %v, want %v", len(*ptr), len(value))
	}

	for k, v := range value {
		if (*ptr)[k] != v {
			t.Errorf("Map value for key %q: got %v, want %v", k, (*ptr)[k], v)
		}
	}
}

func TestOf_Pointer(t *testing.T) {
	value := 42
	valuePtr := &value
	ptr := ref.Of(valuePtr)

	if ptr == nil {
		t.Error("Of() returned nil pointer")
	}

	// Check that we have a pointer to a pointer
	if *ptr != valuePtr {
		t.Error("Of() didn't correctly handle pointer input")
	}

	// Check that we can dereference twice to get original value
	if **ptr != value {
		t.Errorf("Double dereference: got %v, want %v", **ptr, value)
	}
}

func TestOf_Interface(t *testing.T) {
	var value any = "test string"
	ptr := ref.Of(value)

	if ptr == nil {
		t.Error("Of() returned nil pointer")
	}

	if *ptr != value {
		t.Errorf("Of(%v) = pointer to %v, want pointer to %v", value, *ptr, value)
	}

	// Type assertion should work
	if str, ok := (*ptr).(string); !ok || str != "test string" {
		t.Errorf("Type assertion failed: got %v, want %q", str, "test string")
	}
}

func TestOf_ZeroValues(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{"zero int", 0},
		{"empty string", ""},
		{"nil pointer", (*int)(nil)},
		{"empty slice", []int{}},
		{"nil slice", []int(nil)},
		{"empty map", map[string]int{}},
		{"nil map", map[string]int(nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.value.(type) {
			case int:
				ptr := ref.Of(v)
				if ptr == nil || *ptr != v {
					t.Errorf("Of(%v) failed for zero value", v)
				}
			case string:
				ptr := ref.Of(v)
				if ptr == nil || *ptr != v {
					t.Errorf("Of(%q) failed for zero value", v)
				}
			case *int:
				ptr := ref.Of(v)
				if ptr == nil || *ptr != v {
					t.Errorf("Of(%v) failed for zero value", v)
				}
			case []int:
				ptr := ref.Of(v)
				if ptr == nil {
					t.Errorf("Of(%v) returned nil", v)
				}
				// For nil slice, both should be nil
				if v == nil && *ptr != nil {
					t.Errorf("Of(nil slice) should point to nil, got %v", *ptr)
				}
			case map[string]int:
				ptr := ref.Of(v)
				if ptr == nil {
					t.Errorf("Of(%v) returned nil", v)
				}
				// For nil map, both should be nil
				if v == nil && *ptr != nil {
					t.Errorf("Of(nil map) should point to nil, got %v", *ptr)
				}
			}
		})
	}
}
