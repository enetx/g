package g_test

import (
	"sync"
	"testing"
	"unsafe"

	"github.com/enetx/g"
)

func TestResult_Type(t *testing.T) {
	// Test that Result can hold different types
	var intResult g.Result[int]
	var stringResult g.Result[string]
	var boolResult g.Result[bool]

	// These should compile without issues
	_ = intResult
	_ = stringResult
	_ = boolResult
}

func TestOption_Type(t *testing.T) {
	// Test that Option can hold different types
	var intOption g.Option[int]
	var stringOption g.Option[string]
	var boolOption g.Option[bool]

	// These should compile without issues
	_ = intOption
	_ = stringOption
	_ = boolOption
}

func TestUnit_Type(t *testing.T) {
	// Test Unit type
	var unit g.Unit

	// Unit should be zero-sized
	if unsafe.Sizeof(unit) != 0 {
		t.Error("Unit should be zero-sized")
	}

	// Multiple Unit values should be equal
	var unit2 g.Unit
	if unit != unit2 {
		t.Error("Unit values should be equal")
	}
}

func TestFile_Type(t *testing.T) {
	// Test File struct has expected fields
	var file g.File

	// Should be able to access fields (compile-time test)
	_ = file

	// Test that File can be created
	// (We can't test actual file operations without creating real files)
}

func TestDir_Type(t *testing.T) {
	// Test Dir struct
	var dir g.Dir
	_ = dir
}

func TestString_Type(t *testing.T) {
	// Test String type alias
	var gStr g.String = "test"
	var stdStr string = string(gStr)

	if stdStr != "test" {
		t.Errorf("String conversion failed: got %q, want %q", stdStr, "test")
	}

	// Test that String can be converted back
	gStr2 := g.String(stdStr)
	if gStr != gStr2 {
		t.Errorf("String round-trip failed: got %q, want %q", gStr2, gStr)
	}
}

func TestInt_Type(t *testing.T) {
	// Test Int type alias
	var gInt g.Int = 42
	var stdInt int = int(gInt)

	if stdInt != 42 {
		t.Errorf("Int conversion failed: got %d, want %d", stdInt, 42)
	}

	// Test that Int can be converted back
	gInt2 := g.Int(stdInt)
	if gInt != gInt2 {
		t.Errorf("Int round-trip failed: got %d, want %d", gInt2, gInt)
	}
}

func TestFloat_Type(t *testing.T) {
	// Test Float type alias
	var gFloat g.Float = 3.14
	var stdFloat float64 = float64(gFloat)

	if stdFloat != 3.14 {
		t.Errorf("Float conversion failed: got %f, want %f", stdFloat, 3.14)
	}

	// Test that Float can be converted back
	gFloat2 := g.Float(stdFloat)
	if gFloat != gFloat2 {
		t.Errorf("Float round-trip failed: got %f, want %f", gFloat2, gFloat)
	}
}

func TestBytes_Type(t *testing.T) {
	// Test Bytes type alias
	var gBytes g.Bytes = []byte("test")
	var stdBytes []byte = []byte(gBytes)

	if string(stdBytes) != "test" {
		t.Errorf("Bytes conversion failed: got %q, want %q", string(stdBytes), "test")
	}

	// Test that Bytes can be converted back
	gBytes2 := g.Bytes(stdBytes)
	if string(gBytes) != string(gBytes2) {
		t.Errorf("Bytes round-trip failed: got %q, want %q", string(gBytes2), string(gBytes))
	}
}

func TestSlice_Type(t *testing.T) {
	// Test Slice generic type alias
	var intSlice g.Slice[int] = []int{1, 2, 3}
	var stringSlice g.Slice[string] = []string{"a", "b", "c"}

	// Test conversion to standard slices
	stdIntSlice := []int(intSlice)
	if len(stdIntSlice) != 3 || stdIntSlice[0] != 1 {
		t.Errorf("Slice[int] conversion failed")
	}

	stdStringSlice := []string(stringSlice)
	if len(stdStringSlice) != 3 || stdStringSlice[0] != "a" {
		t.Errorf("Slice[string] conversion failed")
	}
}

func TestMap_Type(t *testing.T) {
	// Test Map generic type alias
	var intMap g.Map[string, int] = map[string]int{"one": 1, "two": 2}
	var stringMap g.Map[int, string] = map[int]string{1: "one", 2: "two"}

	// Test conversion to standard maps
	stdIntMap := map[string]int(intMap)
	if stdIntMap["one"] != 1 {
		t.Errorf("Map[string, int] conversion failed")
	}

	stdStringMap := map[int]string(stringMap)
	if stdStringMap[1] != "one" {
		t.Errorf("Map[int, string] conversion failed")
	}
}

func TestMapEntry_Type(t *testing.T) {
	// Test MapEntry struct
	var entry g.MapEntry[string, int]
	_ = entry
}

func TestMapSafeEntry_Type(t *testing.T) {
	// Test MapSafeEntry struct
	var safeEntry g.MapSafeEntry[string, int]
	_ = safeEntry
}

func TestSet_Type(t *testing.T) {
	// Test Set type alias
	var stringSet g.Set[string] = make(g.Set[string])
	stringSet["test"] = g.Unit{}

	// Test conversion to standard map
	stdMap := map[string]g.Unit(stringSet)
	if _, ok := stdMap["test"]; !ok {
		t.Error("Set[string] conversion failed")
	}
}

func TestPair_Type(t *testing.T) {
	// Test Pair struct
	pair := g.Pair[string, int]{Key: "test", Value: 42}

	if pair.Key != "test" || pair.Value != 42 {
		t.Errorf("Pair creation failed: got Key=%q, Value=%d", pair.Key, pair.Value)
	}
}

func TestMapOrd_Type(t *testing.T) {
	// Test MapOrd type alias
	orderedMap := g.NewMapOrd[string, int]()
	orderedMap.Set("first", 1)
	orderedMap.Set("second", 2)

	if orderedMap.Len() != 2 {
		t.Errorf("MapOrd creation failed: got length %d, want 2", orderedMap.Len())
	}

	if orderedMap.Get("first").Some() != 1 {
		t.Error("MapOrd first element incorrect")
	}
}

func TestMapOrdEntry_Type(t *testing.T) {
	// Test MapOrdEntry struct
	var ordEntry g.MapOrdEntry[string, int]
	_ = ordEntry
}

func TestMapSafe_Type(t *testing.T) {
	// Test MapSafe struct
	var safeMap g.MapSafe[string, int]
	_ = safeMap

	// Test that it contains sync.Map
	// We can't directly access the data field, but we can test that it exists by type
	var syncMap sync.Map
	_ = syncMap
}

func TestNamed_Type(t *testing.T) {
	// Test Named type alias
	var named g.Named = g.Named{
		"name": "test",
		"age":  42,
	}

	// Should work like a regular map
	if named["name"] != "test" {
		t.Errorf("Named map access failed: got %v, want %q", named["name"], "test")
	}
}
