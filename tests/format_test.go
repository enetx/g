package g_test

import (
	"testing"

	"github.com/enetx/g"
)

func TestStringFormatting(t *testing.T) {
	// Test Sprintf
	formatted := g.Sprintf("%s is %d years old", "John", 30)
	expected := g.NewString("John is 30 years old")
	if formatted != expected {
		t.Errorf("Sprintf formatting incorrect. Expected: %s, Got: %s", expected, formatted)
	}

	// Test Sprint
	sprinted := g.Sprint("Hello", "World", 42)
	expected = g.NewString("HelloWorld42")
	if sprinted != expected {
		t.Errorf("Sprint formatting incorrect. Expected: %s, Got: %s", expected, sprinted)
	}
}

func TestStringFormat(t *testing.T) {
	// Test case
	values := g.Map[string, any]{
		"name": "John",
		"age":  30,
		"city": "New York",
	}

	format := "Hello, my name is {name}. I am {age} years old and live in {city}."
	formatted := g.Format(format, values)
	expected := g.NewString("Hello, my name is John. I am 30 years old and live in New York.")

	if formatted != expected {
		t.Errorf("Format function incorrect. Expected: %s, Got: %s", expected, formatted)
	}
}
