package g_test

import (
	"regexp"
	"testing"

	"github.com/enetx/g/f"
)

func TestIsComparable(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"int is comparable", 42, true},
		{"string is comparable", "test", true},
		{"slice is not comparable", []int{1, 2, 3}, false},
		{"map is not comparable", map[string]int{"key": 1}, false},
		{"struct is comparable", struct{ Name string }{"test"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.IsComparable(tt.value)
			if got != tt.want {
				t.Errorf("IsComparable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsZero(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"zero int", 0, true},
		{"non-zero int", 42, false},
		{"zero string", "", true},
		{"non-zero string", "test", false},
		{"zero float", 0.0, true},
		{"non-zero float", 3.14, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.value.(type) {
			case int:
				got := f.IsZero(v)
				if got != tt.want {
					t.Errorf("IsZero(%v) = %v, want %v", v, got, tt.want)
				}
			case string:
				got := f.IsZero(v)
				if got != tt.want {
					t.Errorf("IsZero(%v) = %v, want %v", v, got, tt.want)
				}
			case float64:
				got := f.IsZero(v)
				if got != tt.want {
					t.Errorf("IsZero(%v) = %v, want %v", v, got, tt.want)
				}
			}
		})
	}
}

func TestIsEven(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  bool
	}{
		{"even number", 4, true},
		{"odd number", 3, false},
		{"zero is even", 0, true},
		{"negative even", -4, true},
		{"negative odd", -3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.IsEven(tt.value)
			if got != tt.want {
				t.Errorf("IsEven(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestIsOdd(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  bool
	}{
		{"odd number", 3, true},
		{"even number", 4, false},
		{"zero is not odd", 0, false},
		{"negative odd", -3, true},
		{"negative even", -4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.IsOdd(tt.value)
			if got != tt.want {
				t.Errorf("IsOdd(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

func TestMatch(t *testing.T) {
	pattern, err := regexp.Compile(`^\d+$`)
	if err != nil {
		t.Fatal("Failed to compile regex pattern")
	}

	matchFunc := f.Match[string](pattern)

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"numbers only", "12345", true},
		{"contains letters", "123abc", false},
		{"empty string", "", false},
		{"single digit", "7", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchFunc(tt.input)
			if got != tt.want {
				t.Errorf("Match(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	containsFunc := f.Contains("test")

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"contains substring", "this is a test", true},
		{"exact match", "test", true},
		{"doesn't contain", "hello world", false},
		{"empty string", "", false},
		{"case sensitive", "TEST", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsFunc(tt.input)
			if got != tt.want {
				t.Errorf("Contains(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestContainsAnyChars(t *testing.T) {
	containsAnyFunc := f.ContainsAnyChars("abc")

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"contains 'a'", "apple", true},
		{"contains 'b'", "banana", true},
		{"contains 'c'", "cat", true},
		{"contains multiple", "abc", true},
		{"doesn't contain any", "xyz", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsAnyFunc(tt.input)
			if got != tt.want {
				t.Errorf("ContainsAnyChars(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestStartsWith(t *testing.T) {
	startsWithFunc := f.StartsWith[string]("hello")

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"starts with prefix", "hello world", true},
		{"exact match", "hello", true},
		{"doesn't start with", "world hello", false},
		{"empty string", "", false},
		{"case sensitive", "Hello world", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := startsWithFunc(tt.input)
			if got != tt.want {
				t.Errorf("StartsWith(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestEndsWith(t *testing.T) {
	suffix := "world"
	endsWith := f.EndsWith(suffix)

	inputWithSuffix := "hello world"
	if !endsWith(inputWithSuffix) {
		t.Errorf("Expected %q to end with %q, but it did not.", inputWithSuffix, suffix)
	}

	inputWithoutSuffix := "world hello"
	if endsWith(inputWithoutSuffix) {
		t.Errorf("Expected %q not to end with %q, but it did.", inputWithoutSuffix, suffix)
	}

	emptySuffix := f.EndsWith("")
	if !emptySuffix(inputWithSuffix) {
		t.Errorf("Expected any string to end with an empty suffix, but %q did not.", inputWithSuffix)
	}

	if !emptySuffix("") {
		t.Errorf("Expected an empty string to end with an empty suffix, but it did not.")
	}

	exactMatch := "world"
	if !endsWith(exactMatch) {
		t.Errorf("Expected %q to end with %q, but it did not.", exactMatch, suffix)
	}
}

func TestEq(t *testing.T) {
	// Test case 1: Test equality of integers
	if !f.Eq(5)(5) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test inequality of strings
	if f.Eq("hello")("world") {
		t.Errorf("Test case 2: Expected false, got true")
	}
}

func TestNe(t *testing.T) {
	// Test case 1: Test inequality of integers
	if !f.Ne(5)(6) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test equality of strings
	if !f.Ne("hello")("world") {
		t.Errorf("Test case 2: Expected true, got false")
	}
}

func TestEqd(t *testing.T) {
	// Test case 1: Test deep equality of slices
	slice1 := []int{1, 2, 3}
	slice2 := []int{1, 2, 3}
	if !f.Eqd(slice1)(slice2) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test deep inequality of maps
	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"a": 1, "c": 3}
	if f.Eqd(map1)(map2) {
		t.Errorf("Test case 2: Expected false, got true")
	}
}

func TestNed(t *testing.T) {
	// Test case 1: Test deep inequality of slices
	slice1 := []int{1, 2, 3}
	slice2 := []int{1, 2, 4}
	if !f.Ned(slice1)(slice2) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test deep equality of structs
	type Person struct {
		Name string
		Age  int
	}

	person1 := Person{Name: "Alice", Age: 30}
	person2 := Person{Name: "Bob", Age: 25}
	if !f.Ned(person1)(person2) {
		t.Errorf("Test case 2: Expected true, got false")
	}
}

func TestGt(t *testing.T) {
	// Test case 1: Test greater than comparison of integers
	if !f.Gt(10)(15) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test greater than comparison of floats
	if !f.Gt(3.14)(3.1416) {
		t.Errorf("Test case 2: Expected true, got false")
	}
}

func TestGte(t *testing.T) {
	// Test case 1: Test greater than or equal to comparison of integers
	if !f.Gte(10)(10) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test greater than or equal to comparison of floats
	if !f.Gte(3.14)(3.14) {
		t.Errorf("Test case 2: Expected true, got false")
	}
}

func TestLt(t *testing.T) {
	// Test case 1: Test less than comparison of integers
	if !f.Lt(5)(3) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test less than comparison of floats
	if !f.Lt(3.14)(2.71) {
		t.Errorf("Test case 2: Expected true, got false")
	}
}

func TestLte(t *testing.T) {
	// Test case 1: Test less than or equal to comparison of integers
	if !f.Lte(5)(5) {
		t.Errorf("Test case 1: Expected true, got false")
	}

	// Test case 2: Test less than or equal to comparison of floats
	if !f.Lte(3.14)(3.14) {
		t.Errorf("Test case 2: Expected true, got false")
	}
}

func TestFilterRxMatch(t *testing.T) {
	regex := regexp.MustCompile(`^\d+$`)
	matchDigits := f.Match[string](regex)

	tests := []struct {
		input    string
		expected bool
	}{
		{"12345", true},     // Matches: Only digits
		{"abc123", false},   // Does not match: Contains letters
		{"", false},         // Does not match: Empty string
		{" 12345 ", false},  // Does not match: Leading/trailing spaces
		{"0000", true},      // Matches: Only digits
		{"123\n456", false}, // Does not match: Contains newline
	}

	for _, test := range tests {
		result := matchDigits(test.input)
		if result != test.expected {
			t.Errorf("Filter RxMatch failed for input %q: expected %v, got %v", test.input, test.expected, result)
		}
	}
}

func TestEqi(t *testing.T) {
	t.Run("nil interface", func(t *testing.T) {
		var a, b any
		if !f.Eqi(a)(b) {
			t.Error("Two nil interfaces should be equal")
		}

		c := 42
		if f.Eqi(a)(c) {
			t.Error("nil interface should not equal non-nil value")
		}
	})

	t.Run("comparable types use fast path", func(t *testing.T) {
		// Integers
		if !f.Eqi(42)(42) {
			t.Error("Same integers should be equal")
		}
		if f.Eqi(42)(43) {
			t.Error("Different integers should not be equal")
		}

		// Strings
		if !f.Eqi("hello")("hello") {
			t.Error("Same strings should be equal")
		}
		if f.Eqi("hello")("world") {
			t.Error("Different strings should not be equal")
		}

		// Floats
		if !f.Eqi(3.14)(3.14) {
			t.Error("Same floats should be equal")
		}
		if f.Eqi(3.14)(2.71) {
			t.Error("Different floats should not be equal")
		}

		// Booleans
		if !f.Eqi(true)(true) {
			t.Error("Same booleans should be equal")
		}
		if f.Eqi(true)(false) {
			t.Error("Different booleans should not be equal")
		}
	})

	t.Run("functions", func(t *testing.T) {
		fn1 := func(x int) int { return x + 1 }
		fn2 := func(x int) int { return x + 2 }

		// Same function
		if !f.Eqi(fn1)(fn1) {
			t.Error("Same function should be equal to itself")
		}

		// Different functions
		if f.Eqi(fn1)(fn2) {
			t.Error("Different functions should not be equal")
		}

		// nil functions
		var nilFn1, nilFn2 func()
		if !f.Eqi(nilFn1)(nilFn2) {
			t.Error("Two nil functions should be equal")
		}

		// nil vs non-nil function
		nonNilFn := func() {}
		if f.Eqi(nilFn1)(nonNilFn) {
			t.Error("nil function should not equal non-nil function")
		}
	})

	t.Run("slices", func(t *testing.T) {
		slice1 := []int{1, 2, 3}
		slice2 := []int{1, 2, 3}
		slice3 := []int{4, 5, 6}

		// Identity comparison - same slice
		if !f.Eqi(slice1)(slice1) {
			t.Error("Same slice should be equal to itself")
		}

		// Different slices with same content
		if f.Eqi(slice1)(slice2) {
			t.Error("Different slices with same content should not be equal (identity comparison)")
		}

		// Different content
		if f.Eqi(slice1)(slice3) {
			t.Error("Different slices with different content should not be equal")
		}

		// nil slices
		var nilSlice1, nilSlice2 []int
		if !f.Eqi(nilSlice1)(nilSlice2) {
			t.Error("Two nil slices should be equal")
		}

		// nil vs non-nil slice
		emptySlice := []int{}
		if f.Eqi(nilSlice1)(emptySlice) {
			t.Error("nil slice should not equal non-nil empty slice")
		}

		// Slices with same underlying array but different length
		base := []int{1, 2, 3, 4, 5}
		sub1 := base[:3]
		sub2 := base[:4]
		if f.Eqi(sub1)(sub2) {
			t.Error("Slices with same base but different length should not be equal")
		}

		// Same sub-slice
		sub3 := base[:3]
		if !f.Eqi(sub1)(sub3) {
			t.Error("Same sub-slices should be equal")
		}
	})

	t.Run("maps", func(t *testing.T) {
		map1 := map[string]int{"a": 1, "b": 2}
		map2 := map[string]int{"a": 1, "b": 2}
		map3 := map[string]int{"c": 3, "d": 4}

		// Identity comparison - same map
		if !f.Eqi(map1)(map1) {
			t.Error("Same map should be equal to itself")
		}

		// Different maps with same content
		if f.Eqi(map1)(map2) {
			t.Error("Different maps with same content should not be equal (identity comparison)")
		}

		// Different content
		if f.Eqi(map1)(map3) {
			t.Error("Different maps with different content should not be equal")
		}

		// nil maps
		var nilMap1, nilMap2 map[string]int
		if !f.Eqi(nilMap1)(nilMap2) {
			t.Error("Two nil maps should be equal")
		}

		// nil vs non-nil map
		emptyMap := make(map[string]int)
		if f.Eqi(nilMap1)(emptyMap) {
			t.Error("nil map should not equal non-nil empty map")
		}
	})

	t.Run("channels", func(t *testing.T) {
		ch1 := make(chan int)
		ch2 := make(chan int)

		// Same channel
		if !f.Eqi(ch1)(ch1) {
			t.Error("Same channel should be equal to itself")
		}

		// Different channels
		if f.Eqi(ch1)(ch2) {
			t.Error("Different channels should not be equal")
		}

		// nil channels
		var nilCh1, nilCh2 chan int
		if !f.Eqi(nilCh1)(nilCh2) {
			t.Error("Two nil channels should be equal")
		}

		// nil vs non-nil channel
		if f.Eqi(nilCh1)(ch1) {
			t.Error("nil channel should not equal non-nil channel")
		}

		// Test with any type containing channels
		var anyCh1 any = make(chan string)
		anyCh2 := anyCh1
		var anyCh3 any = make(chan string)

		if !f.Eqi(anyCh1)(anyCh2) {
			t.Error("Same channel stored as any should be equal")
		}

		if f.Eqi(anyCh1)(anyCh3) {
			t.Error("Different channels stored as any should not be equal")
		}

		// Test nil channel stored as any
		var nilChAny any = (chan bool)(nil)
		var nilChAny2 any = (chan bool)(nil)

		if !f.Eqi(nilChAny)(nilChAny2) {
			t.Error("Two nil channels stored as any should be equal")
		}

		// Different channel types
		var intCh any = make(chan int)
		var strCh any = make(chan string)

		if f.Eqi(intCh)(strCh) {
			t.Error("Channels of different types should not be equal")
		}
	})

	t.Run("structs with non-comparable fields", func(t *testing.T) {
		type StructWithSlice struct {
			Name  string
			Items []int
		}

		s1 := StructWithSlice{Name: "test", Items: []int{1, 2, 3}}
		s2 := StructWithSlice{Name: "test", Items: []int{1, 2, 3}}
		s3 := StructWithSlice{Name: "test", Items: []int{4, 5, 6}}

		// Deep equality for structs with non-comparable fields
		if !f.Eqi(s1)(s2) {
			t.Error("Structs with same content should be deeply equal")
		}

		if f.Eqi(s1)(s3) {
			t.Error("Structs with different content should not be equal")
		}
	})

	t.Run("arrays with non-comparable elements", func(t *testing.T) {
		// Arrays with slices - non-comparable type
		type ArrayWithSlice [2][]int

		slice1 := []int{1, 2, 3}
		slice2 := []int{4, 5, 6}

		a1 := ArrayWithSlice{slice1, slice2}
		a2 := ArrayWithSlice{slice1, slice2}
		a3 := ArrayWithSlice{slice1, []int{4, 5, 6}} // Different slice instance

		// Arrays with slices use DeepEqual, which compares contents
		if !f.Eqi(a1)(a2) {
			t.Error("Arrays with same slices should be deeply equal")
		}

		// a1 and a3 have same content but different slice instances - DeepEqual compares content
		if !f.Eqi(a1)(a3) {
			t.Error("Arrays with same slice content should be deeply equal")
		}

		// Different content
		a4 := ArrayWithSlice{slice1, []int{7, 8, 9}}
		if f.Eqi(a1)(a4) {
			t.Error("Arrays with different slice content should not be equal")
		}
	})

	t.Run("mixed types in any", func(t *testing.T) {
		var a any = 42
		var b any = "42"

		if f.Eqi(a)(b) {
			t.Error("Different types should not be equal even when stored in any")
		}

		var c any = 42
		if !f.Eqi(a)(c) {
			t.Error("Same value and type in any should be equal")
		}
	})

	t.Run("type consistency", func(t *testing.T) {
		// Different types of slices
		intSlice := []int{1, 2, 3}
		var anySlice any = []int{1, 2, 3}

		// Use Eqi with any type
		if f.Eqi[any](intSlice)(anySlice) {
			t.Error("Different slice instances should not be equal")
		}

		// Same slice stored as any
		var sameSlice any = intSlice
		if !f.Eqi[any](intSlice)(sameSlice) {
			t.Error("Same slice instance should be equal")
		}
	})

	t.Run("interfaces", func(t *testing.T) {
		// Test with comparable values stored in interfaces
		var iface1 any = "hello"
		var iface2 any = "hello"
		var iface3 any = "world"

		// Same value but different instances (comparable type)
		if !f.Eqi(iface1)(iface2) {
			t.Error("Interfaces with same comparable value should be equal")
		}

		if f.Eqi(iface1)(iface3) {
			t.Error("Interfaces with different values should not be equal")
		}

		// nil interfaces
		var nilIface1, nilIface2 any
		if !f.Eqi(nilIface1)(nilIface2) {
			t.Error("Two nil interfaces should be equal")
		}
	})

	t.Run("pointers", func(t *testing.T) {
		// Test with pointers
		x := 42
		y := 42
		ptr1 := &x
		ptr2 := &x
		ptr3 := &y

		if !f.Eqi(ptr1)(ptr2) {
			t.Error("Same pointer should be equal")
		}

		if f.Eqi(ptr1)(ptr3) {
			t.Error("Different pointers should not be equal")
		}

		// nil pointers
		var nilPtr1, nilPtr2 *int
		if !f.Eqi(nilPtr1)(nilPtr2) {
			t.Error("Two nil pointers should be equal")
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		// Test with complex types that might trigger different code paths
		type ComplexStruct struct {
			IntField   int
			FloatField float64
		}

		// These are comparable, so they should use the fast path
		cs1 := ComplexStruct{IntField: 1, FloatField: 2.0}
		cs2 := ComplexStruct{IntField: 1, FloatField: 2.0}
		cs3 := ComplexStruct{IntField: 2, FloatField: 3.0}

		if !f.Eqi(cs1)(cs2) {
			t.Error("Same struct values should be equal")
		}

		if f.Eqi(cs1)(cs3) {
			t.Error("Different struct values should not be equal")
		}

		// Test recursive types - pointers are comparable but point to different objects
		type Node struct {
			Value int
			Next  *Node
		}

		node1 := &Node{Value: 1, Next: nil}
		node2 := &Node{Value: 1, Next: nil}

		// These are different pointers, so they should not be equal
		if f.Eqi(node1)(node2) {
			t.Error("Different node pointers should not be equal")
		}

		// Same pointer should be equal
		if !f.Eqi(node1)(node1) {
			t.Error("Same node pointer should be equal to itself")
		}
	})
}

func TestNei(t *testing.T) {
	t.Run("nil interface", func(t *testing.T) {
		var a, b any
		if f.Nei(a)(b) {
			t.Error("Two nil interfaces should be equal, so Nei should return false")
		}

		c := 42
		if !f.Nei(a)(c) {
			t.Error("nil interface should not equal non-nil value, so Nei should return true")
		}
	})

	t.Run("comparable types", func(t *testing.T) {
		// Integers
		if f.Nei(42)(42) {
			t.Error("Same integers should be equal, so Nei should return false")
		}
		if !f.Nei(42)(43) {
			t.Error("Different integers should not be equal, so Nei should return true")
		}

		// Strings
		if f.Nei("hello")("hello") {
			t.Error("Same strings should be equal, so Nei should return false")
		}
		if !f.Nei("hello")("world") {
			t.Error("Different strings should not be equal, so Nei should return true")
		}
	})

	t.Run("functions", func(t *testing.T) {
		fn1 := func(x int) int { return x + 1 }
		fn2 := func(x int) int { return x + 2 }

		// Same function
		if f.Nei(fn1)(fn1) {
			t.Error("Same function should be equal to itself, so Nei should return false")
		}

		// Different functions
		if !f.Nei(fn1)(fn2) {
			t.Error("Different functions should not be equal, so Nei should return true")
		}
	})

	t.Run("slices", func(t *testing.T) {
		slice1 := []int{1, 2, 3}
		slice2 := []int{1, 2, 3}

		// Identity comparison - same slice
		if f.Nei(slice1)(slice1) {
			t.Error("Same slice should be equal to itself, so Nei should return false")
		}

		// Different slices with same content (identity comparison)
		if !f.Nei(slice1)(slice2) {
			t.Error("Different slices are not equal by identity, so Nei should return true")
		}
	})

	t.Run("maps", func(t *testing.T) {
		map1 := map[string]int{"a": 1, "b": 2}
		map2 := map[string]int{"a": 1, "b": 2}

		// Identity comparison - same map
		if f.Nei(map1)(map1) {
			t.Error("Same map should be equal to itself, so Nei should return false")
		}

		// Different maps with same content
		if !f.Nei(map1)(map2) {
			t.Error("Different maps are not equal by identity, so Nei should return true")
		}
	})

	t.Run("channels", func(t *testing.T) {
		ch1 := make(chan int)
		ch2 := make(chan int)

		// Same channel
		if f.Nei(ch1)(ch1) {
			t.Error("Same channel should be equal to itself, so Nei should return false")
		}

		// Different channels
		if !f.Nei(ch1)(ch2) {
			t.Error("Different channels should not be equal, so Nei should return true")
		}
	})

	t.Run("structs with non-comparable fields", func(t *testing.T) {
		type StructWithSlice struct {
			Name  string
			Items []int
		}

		s1 := StructWithSlice{Name: "test", Items: []int{1, 2, 3}}
		s2 := StructWithSlice{Name: "test", Items: []int{1, 2, 3}}
		s3 := StructWithSlice{Name: "test", Items: []int{4, 5, 6}}

		// Deep equality for structs with non-comparable fields
		if f.Nei(s1)(s2) {
			t.Error("Structs with same content should be deeply equal, so Nei should return false")
		}

		if !f.Nei(s1)(s3) {
			t.Error("Structs with different content should not be equal, so Nei should return true")
		}
	})

	t.Run("Nei is opposite of Eqi", func(t *testing.T) {
		// Test various types to ensure Nei(x)(y) == !Eqi(x)(y)
		testCases := []struct {
			name string
			a, b any
		}{
			{"same int", 42, 42},
			{"different int", 42, 43},
			{"same string", "hello", "hello"},
			{"different string", "hello", "world"},
			{"same slice", []int{1, 2}, []int{1, 2}},
			{"nil values", nil, nil},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				eqdResult := f.Eqi(tc.a)(tc.b)
				nedResult := f.Nei(tc.a)(tc.b)

				if eqdResult == nedResult {
					t.Errorf("Nei should be opposite of Eqi for %v and %v", tc.a, tc.b)
				}
			})
		}
	})
}
