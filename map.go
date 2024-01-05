package g

import (
	"fmt"
	"maps"
	"reflect"
	"strings"
)

// NewMap creates a new Map of the specified size or an empty Map if no size is provided.
func NewMap[K comparable, V any](size ...int) Map[K, V] {
	if len(size) == 0 {
		return make(Map[K, V], 0)
	}

	return make(Map[K, V], size[0])
}

// MapFromStd creates an Map from a given Go map.
func MapFromStd[K comparable, V any](stdmap map[K]V) Map[K, V] { return stdmap }

// Random returns a new map containing a single randomly selected key-value pair from the original map.
//
// Parameters:
//
// None.
//
// Returns:
//
// - Map[K, V]: A new Map containing a single randomly selected key-value pair.
//
// Example usage:
//
//	myMap := g.Map[string, int]{"one": 1, "two": 2, "three": 3, "four": 4, "five": 5}
//	randomMap := myMap.Random()
//
// The resulting map will contain one randomly selected key-value pair from the original map.
func (m Map[K, V]) Random() Map[K, V] {
	if m.Empty() {
		return m
	}

	key := m.Keys()[0]

	return NewMap[K, V]().Set(key, m.Get(key))
}

// RandomSample returns a new map containing a random sample of key-value pairs from the original map.
//
// Parameters:
//
// - sequence int: The number of unique key-value pairs to include in the random sample.
//
// Returns:
//
// - Map[K, V]: A new Map containing a random sample of unique key-value pairs.
//
// Example usage:
//
//	myMap := g.Map[string, int]{"one": 1, "two": 2, "three": 3, "four": 4, "five": 5}
//	sampledMap := myMap.RandomSample(3)
//
// The resulting map will contain 3 unique key-value pairs randomly selected from the original map.
func (m Map[K, V]) RandomSample(sequence int) Map[K, V] {
	if m.Empty() {
		return m
	}

	keys := m.Keys()

	if sequence >= keys.Len() {
		return m.Clone()
	}

	nmap := NewMap[K, V](sequence)
	keys[0:sequence].Iter().ForEach(func(key K) { nmap.Set(key, m.Get(key)) })

	return nmap
}

// RandomRange returns a new map containing a random range of key-value pairs from the original map.
//
// Parameters:
//
// - from int: The starting index of the range (inclusive).
// - to int: The ending index of the range (exclusive).
//
// Returns:
//
// - Map[K, V]: A new Map containing a random range of key-value pairs from the specified subrange.
//
// Example usage:
//
//	myMap := g.Map[string, int]{"one": 1, "two": 2, "three": 3, "four": 4, "five": 5}
//	subrangeMap := myMap.RandomRange(1, 4)
//
// The resulting map will contain a random range of key-value pairs from index 1 (inclusive) to 4 (exclusive) of the original map.
func (m Map[K, V]) RandomRange(from, to int) Map[K, V] {
	if m.Empty() {
		return m
	}

	if from < 0 {
		from = 0
	}

	if to > m.Len() {
		to = m.Len()
	}

	if from >= to || from >= m.Len() {
		return NewMap[K, V]()
	}

	sequence := Int(from).RandomRange(Int(to)).Std()

	return m.RandomSample(sequence)
}

// Invert inverts the keys and values of the Map, returning a new Map with values as keys and
// keys as values. Note that the inverted Map will have 'any' as the key type, since not all value
// types are guaranteed to be comparable.
func (m Map[K, V]) Invert() Map[any, K] {
	result := NewMap[any, K](m.Len())
	m.ForEach(func(k K, v V) { result.Set(v, k) })

	return result
}

// Keys returns a slice of the Map's keys.
func (m Map[K, V]) Keys() Slice[K] {
	keys := NewSlice[K](0, m.Len())
	m.ForEach(func(k K, _ V) { keys = keys.Append(k) })

	return keys
}

// Values returns a slice of the Map's values.
func (m Map[K, V]) Values() Slice[V] {
	values := NewSlice[V](0, m.Len())
	m.ForEach(func(_ K, v V) { values = values.Append(v) })

	return values
}

// Contains checks if the Map contains the specified key.
func (m Map[K, V]) Contains(key K) bool {
	_, ok := m[key]
	return ok
}

// Clone creates a new Map that is a copy of the original Map.
func (m Map[K, V]) Clone() Map[K, V] { return maps.Clone(m) }

// Copy copies the source Map's key-value pairs to the target Map.
func (m Map[K, V]) Copy(src Map[K, V]) Map[K, V] {
	maps.Copy(m, src)
	return m
}

// Delete removes the specified keys from the Map.
func (m Map[K, V]) Delete(keys ...K) Map[K, V] {
	for _, key := range keys {
		delete(m, key)
	}

	return m
}

// Std converts the Map to a regular Go map.
func (m Map[K, V]) Std() map[K]V { return m }

// Map applies a function to each key-value pair in the Map and returns a new Map with the
// results. The provided function 'fn' should take a key and a value as input parameters and return
// a new key-value pair.
//
// Parameters:
//
// - fn func(K, V) (K, V): A function that takes a key and a value as input parameters and returns
// a new key-value pair.
//
// Returns:
//
// - Map[K, V]: A new Map containing the key-value pairs resulting from applying the provided
// function to each key-value pair in the original Map.
//
// Example usage:
//
//	mappedMap := originalMap.Map(func(key K, value V) (K, V) {
//		return key, value * 2
//	})
func (m Map[K, V]) Map(fn func(K, V) (K, V)) Map[K, V] {
	result := NewMap[K, V](m.Len())
	m.ForEach(func(k K, v V) { result.Set(fn(k, v)) })

	return result
}

// Filter filters the Map based on a given function and returns a new Map containing the matching
// key-value pairs. The provided function 'fn' should take a key and a value as input parameters
// and return a boolean value.
// If the function returns true, the key-value pair will be included in the resulting Map.
//
// Parameters:
//
// - fn func(K, V) bool: A function that takes a key and a value as input parameters and returns a
// boolean value.
//
// Returns:
//
// - Map[K, V]: A new Map containing the key-value pairs for which the provided function returned
// true.
//
// Example usage:
//
//	filteredMap := originalMap.Filter(func(key K, value V) bool {
//		return value >= 10
//	})
func (m Map[K, V]) Filter(fn func(K, V) bool) Map[K, V] {
	result := NewMap[K, V]()

	m.ForEach(func(k K, v V) {
		if fn(k, v) {
			result.Set(k, v)
		}
	})

	return result
}

// ForEach applies a function to each key-value pair in the Map.
// The provided function 'fn' should take a key and a value as input parameters and perform an
// operation.
// This function is useful for side effects, as it does not return a new Map.
//
// Parameters:
//
// - fn func(K, V): A function that takes a key and a value as input parameters and performs an
// operation.
//
// Example usage:
//
//	originalMap.ForEach(func(key K, value V) {
//		fmt.Printf("Key: %v, Value: %v\n", key, value)
//	})
func (m Map[K, V]) ForEach(fn func(K, V)) {
	for key, val := range m {
		fn(key, val)
	}
}

// Range applies a given function to each key-value pair in the Map until the function returns false.
//
// The provided function 'fn' should take a key and a value as input parameters and return a boolean.
// If the function returns false for any key-value pair, the iteration stops.
//
// Parameters:
//
// - fn func(K, V) bool: A function that takes a key and a value as input parameters and returns a boolean.
// If it returns false, the iteration will stop.
//
// Example usage:
//
//	originalMap.Range(func(key K, value V) bool {
//	    fmt.Printf("Key: %v, Value: %v\n", key, value)
//	    return key != stopKey // Stop iteration condition
//	})
func (m Map[K, V]) Range(fn func(K, V) bool) {
	for key, val := range m {
		if !fn(key, val) {
			break
		}
	}
}

// Eq checks if two Maps are equal.
func (m Map[K, V]) Eq(other Map[K, V]) bool {
	if m.Len() != other.Len() {
		return false
	}

	for key, value := range m {
		if value2, ok := other[key]; !ok || !reflect.DeepEqual(value, value2) {
			return false
		}
	}

	return true
}

// String returns a string representation of the Map.
func (m Map[K, V]) String() string {
	var builder strings.Builder

	m.ForEach(func(k K, v V) { builder.WriteString(fmt.Sprintf("%v:%v, ", k, v)) })

	return String(builder.String()).TrimRight(", ").Format("Map{%s}").Std()
}

// GetOrDefault returns the value for a key. If the key does not exist, returns the default value
// instead. This function is useful when you want to provide a fallback value for keys that may not
// be present in the Map.
//
// Parameters:
//
// - key K: The key for which to retrieve the value.
//
// - defaultValue V: The default value to return if the key does not exist in the Map.
//
// Returns:
//
// - V: The value associated with the key if it exists in the Map, or the default value if the key
// is not found.
//
// Example usage:
//
//	value := m.GetOrDefault("someKey", "defaultValue")
func (m Map[K, V]) GetOrDefault(key K, defaultValue V) V {
	if value, ok := m[key]; ok {
		return value
	}

	return defaultValue
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

	m.Set(key, defaultValue)

	return defaultValue
}

// Clear removes all key-value pairs from the Map.
func (m Map[K, V]) Clear() Map[K, V] { clear(m); return m }

// Empty checks if the Map is empty.
func (m Map[K, V]) Empty() bool { return m.Len() == 0 }

// Get retrieves the value associated with the given key.
func (m Map[K, V]) Get(k K) V { return m[k] }

// Len returns the number of key-value pairs in the Map.
func (m Map[K, V]) Len() int { return len(m) }

// Ne checks if two Maps are not equal.
func (m Map[K, V]) Ne(other Map[K, V]) bool { return !m.Eq(other) }

// NotEmpty checks if the Map is not empty.
func (m Map[K, V]) NotEmpty() bool { return !m.Empty() }

// Set sets the value for the given key in the Map.
func (m Map[K, V]) Set(k K, v V) Map[K, V] { m[k] = v; return m }

// Print prints the key-value pairs of the Map to the standard output (console)
// and returns the Map unchanged.
func (m Map[K, V]) Print() Map[K, V] { fmt.Println(m); return m }
