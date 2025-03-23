package g

import (
	"maps"

	"github.com/enetx/g/f"
)

// NewMap creates a new Map of the specified size or an empty Map if no size is provided.
func NewMap[K comparable, V any](size ...Int) Map[K, V] {
	if len(size) == 0 {
		return make(Map[K, V], 0)
	}

	return make(Map[K, V], size[0])
}

// Transform applies a transformation function to the Map and returns the result.
func (m Map[K, V]) Transform(fn func(Map[K, V]) Map[K, V]) Map[K, V] { return fn(m) }

// Iter returns an iterator (SeqMap[K, V]) for the Map, allowing for sequential iteration
// over its key-value pairs. It is commonly used in combination with higher-order functions,
// such as 'ForEach', to perform operations on each key-value pair of the Map.
//
// Returns:
//
// - SeqMap[K, V], which can be used for sequential iteration over the key-value pairs of the Map.
//
// Example usage:
//
//	myMap := g.Map[string, int]{"one": 1, "two": 2, "three": 3}
//	iterator := myMap.Iter()
//	iterator.ForEach(func(key string, value int) {
//		// Perform some operation on each key-value pair
//		fmt.Printf("%s: %d\n", key, value)
//	})
//
// The 'Iter' method provides a convenient way to traverse the key-value pairs of a Map
// in a functional style, enabling operations like mapping or filtering.
func (m Map[K, V]) Iter() SeqMap[K, V] { return seqMap(m) }

// Invert inverts the keys and values of the Map, returning a new Map with values as keys and
// keys as values. Note that the inverted Map will have 'any' as the key type, since not all value
// types are guaranteed to be comparable.
func (m Map[K, V]) Invert() Map[any, K] {
	result := NewMap[any, K](m.Len())
	for k, v := range m {
		result.Set(v, k)
	}

	return result
}

// Keys returns a slice of the Map's keys.
func (m Map[K, V]) Keys() Slice[K] { return m.Iter().Keys().Collect() }

// Values returns a slice of the Map's values.
func (m Map[K, V]) Values() Slice[V] { return m.Iter().Values().Collect() }

// Contains checks if the Map contains the specified key.
func (m Map[K, V]) Contains(key K) bool {
	_, ok := m[key]
	return ok
}

// Clone creates a new Map that is a copy of the original Map.
func (m Map[K, V]) Clone() Map[K, V] { return maps.Clone(m) }

// Copy copies the source Map's key-value pairs to the target Map.
func (m Map[K, V]) Copy(src Map[K, V]) { maps.Copy(m, src) }

// Delete removes the specified keys from the Map.
func (m Map[K, V]) Delete(keys ...K) {
	for _, key := range keys {
		delete(m, key)
	}
}

// Std converts the Map to a regular Go map.
func (m Map[K, V]) Std() map[K]V { return m }

// ToMapOrd converts a standard Map to an ordered Map.
func (m Map[K, V]) ToMapOrd() MapOrd[K, V] {
	mo := NewMapOrd[K, V](m.Len())

	for k, v := range m {
		mo.Set(k, v)
	}

	return mo
}

// ToMapSafe converts a standard Map to a thread-safe Map.
func (m Map[K, V]) ToMapSafe() *MapSafe[K, V] {
	ms := NewMapSafe[K, V]()
	ms.data = m

	return ms
}

// Eq checks if two Maps are equal.
func (m Map[K, V]) Eq(other Map[K, V]) bool {
	n := len(m)

	if n != len(other) {
		return false
	}

	if n == 0 {
		return true
	}

	key := m.Iter().Take(1).Keys().Collect()[0]
	comparable := f.IsComparable(key) && f.IsComparable(m[key])

	for k, v1 := range m {
		v2, ok := other[k]
		if !ok || (comparable && !f.Eq[any](v1)(v2)) || (!comparable && !f.Eqd(v1)(v2)) {
			return false
		}
	}

	return true
}

// String returns a string representation of the Map.
func (m Map[K, V]) String() string {
	builder := NewBuilder()

	for k, v := range m {
		builder.Write(Sprintf("{}:{}, ", k, v))
	}

	return builder.String().StripSuffix(", ").Sprintf("Map\\{{}\\}").Std()
}

// GetOrSet returns the value for a key. If the key exists in the Map, it returns the associated value.
// If the key does not exist, it sets the key to the provided default value and returns that value.
// This function is useful when you want to both retrieve and potentially set a default value for keys
// that may or may not be present in the Map.
//
// Parameters:
//
// - key K: The key for which to retrieve the value.
//
// - defaultValue V: The default value to return if the key does not exist in the Map.
// If the key is not found, this default value will also be set for the key in the Map.
//
// Returns:
//
// - V: The value associated with the key if it exists in the Map, or the default value if the key is not found.
//
// Eaxmple usage:
//
//	// Create a new ordered Map called "gos" with string keys and integer pointers as values
//	gos := g.NewMap[string, *int]()
//
//	// Use GetOrSet to set the value for the key "root" to 3 if it doesn't exist,
//	// and then print whether the value is equal to 3.
//	gos.GetOrSet("root", ref.Of(3))
//	fmt.Println(*gos.Get("root") == 3) // Should print "true"
//
//	// Use GetOrSet to retrieve the value for the key "root" (which is 3), multiply it by 2,
//	// and then print whether the value is equal to 6.
//	*gos.GetOrSet("root", ref.Of(10)) *= 2
//	fmt.Println(*gos.Get("root") == 6) // Should print "true"
//
// In this example, you first create an ordered Map "gos" with string keys and integer pointers as values.
// Then, you use GetOrSet to set and retrieve values for the key "root" with default values of 3 and perform
// multiplication operations, demonstrating the behavior of GetOrSet.
func (m Map[K, V]) GetOrSet(key K, defaultValue V) V {
	if value, ok := m[key]; ok {
		return value
	}

	m[key] = defaultValue

	return defaultValue
}

// Clear removes all key-value pairs from the Map.
func (m Map[K, V]) Clear() Map[K, V] { clear(m); return m }

// Empty checks if the Map is empty.
func (m Map[K, V]) Empty() bool { return len(m) == 0 }

// Get retrieves the value associated with the given key.
func (m Map[K, V]) Get(k K) Option[V] {
	if v, ok := m[k]; ok {
		return Some(v)
	}

	return None[V]()
}

// Len returns the number of key-value pairs in the Map.
func (m Map[K, V]) Len() Int { return Int(len(m)) }

// Ne checks if two Maps are not equal.
func (m Map[K, V]) Ne(other Map[K, V]) bool { return !m.Eq(other) }

// NotEmpty checks if the Map is not empty.
func (m Map[K, V]) NotEmpty() bool { return !m.Empty() }

// Set sets the value for the given key in the Map.
func (m Map[K, V]) Set(key K, value V) { m[key] = value }

// Print writes the key-value pairs of the Map to the standard output (console)
// and returns the Map unchanged.
func (m Map[K, V]) Print() Map[K, V] { Print(m); return m }

// Println writes the key-value pairs of the Map to the standard output (console) with a newline
// and returns the Map unchanged.
func (m Map[K, V]) Println() Map[K, V] { Println(m); return m }
