package g

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// NewMapOrd creates a new ordered Map with the specified size (if provided).
// An ordered Map is an Map that maintains the order of its key-value pairs based on the
// insertion order. If no size is provided, the default size will be used.
//
// Parameters:
//
// - size ...int: (Optional) The initial size of the ordered Map. If not provided, a default size
// will be used.
//
// Returns:
//
// - MapOrd[K, V]: Ordered Map with the specified initial size (or default
// size if not provided).
//
// Example usage:
//
//	mapOrd := g.NewMapOrd[string, int](10)
//
// Creates a new ordered Map with an initial size of 10.
func NewMapOrd[K, V any](size ...int) MapOrd[K, V] {
	if len(size) == 0 {
		return make(MapOrd[K, V], 0)
	}

	return make(MapOrd[K, V], 0, size[0])
}

// Iter returns an iterator (*liftIterMO) for the ordered Map, allowing for sequential iteration
// over its key-value pairs. It is commonly used in combination with higher-order functions,
// such as 'ForEach', to perform operations on each key-value pair of the ordered Map.
//
// Returns:
//
// A pointer to a liftIterMO, which can be used for sequential iteration over the key-value pairs of the ordered Map.
//
// Example usage:
//
//	iter := g.NewMapOrd[int, int]()
//	iter.
//		Set(1, 1).
//		Set(2, 2).
//		Set(3, 3).
//		Iter()
//
//	iter.ForEach(func(k, v int) {
//	    // Process key-value pair
//	})
//
// The 'Iter' method provides a convenient way to traverse the key-value pairs of an ordered Map
// in a functional style, enabling operations like mapping or filtering.
func (mo MapOrd[K, V]) Iter() seqMapOrd[K, V] { return liftMO(mo) }

// MapOrdFromMap converts a standard Map to an ordered Map.
// The resulting ordered Map will maintain the order of its key-value pairs based on the order of
// insertion.
// This function is useful when you want to create an ordered Map from an existing Map.
//
// Parameters:
//
// - m Map[K, V]: The input Map to be converted to an ordered Map.
//
// Returns:
//
// - MapOrd[K, V]: New ordered Map containing the same key-value pairs as the
// input Map.
//
// Example usage:
//
//	mapOrd := g.MapOrdFromMap[string, int](hmap)
//
// Converts the standard Map 'hmap' to an ordered Map.
func MapOrdFromMap[K comparable, V any](m Map[K, V]) MapOrd[K, V] {
	mo := NewMapOrd[K, V](m.Len())
	m.Iter().ForEach(func(k K, v V) { mo.Set(k, v) })

	return mo
}

// MapOrdFromStd converts a standard Go map to an ordered Map.
// The resulting ordered Map will maintain the order of its key-value pairs based on the order of
// insertion.
// This function is useful when you want to create an ordered Map from an existing Go map.
//
// Parameters:
//
// - m map[K]V: The input Go map to be converted to an ordered Map.
//
// Returns:
//
// - MapOrd[K, V]: New ordered Map containing the same key-value pairs as the
// input Go map.
//
// Example usage:
//
//	mapOrd := g.MapOrdFromStd[string, int](goMap)
//
// Converts the standard Go map 'map[K]V' to an ordered Map.
func MapOrdFromStd[K comparable, V any](m map[K]V) MapOrd[K, V] { return MapOrdFromMap(MapFromStd(m)) }

// SortBy sorts the ordered Map by a custom comparison function.
// The comparison function should return true if the element at index i should be sorted before the
// element at index j. This function is useful when you want to sort the key-value pairs in an
// ordered Map based on a custom comparison logic.
//
// Parameters:
//
// - f func(i, j int) bool: The custom comparison function used for sorting the ordered Map.
//
// Returns:
//
// - MapOrd[K, V]: Same ordered Map, sorted according to the custom comparison
// function.
//
// Example usage:
//
//	hmapo.SortBy(func(a, b g.Pair[string, int]) bool { return a.Key < b.Key })
//	hmapo.SortBy(func(a, b g.Pair[string, int]) bool { return a.Value < b.Value })
func (mo MapOrd[K, V]) SortBy(fn func(a, b Pair[K, V]) bool) MapOrd[K, V] {
	sort.Slice(mo, func(i, j int) bool {
		return fn(mo[i], mo[j])
	})

	return mo
}

// Clone creates a new ordered Map with the same key-value pairs.
func (mo MapOrd[K, V]) Clone() MapOrd[K, V] {
	result := NewMapOrd[K, V](mo.Len())
	mo.Iter().ForEach(func(k K, v V) { result.Set(k, v) })

	return result
}

// Copy copies key-value pairs from the source ordered Map to the current ordered Map.
func (mo *MapOrd[K, V]) Copy(src MapOrd[K, V]) MapOrd[K, V] {
	src.Iter().ForEach(func(k K, v V) { mo.Set(k, v) })
	return *mo
}

// ToMap converts the ordered Map to a standard Map.
// func (mo MapOrd[K, V]) ToMap() Map[K, V] {
// 	m := NewMap[K, V](mo.Len())
// 	mo.Iter().ForEach(func(k K, v V) { m.Set(k, v) })

// 	return m
// }

// Set sets the value for the specified key in the ordered Map.
func (mo *MapOrd[K, V]) Set(key K, value V) *MapOrd[K, V] {
	if i := mo.index(key); i != -1 {
		(*mo)[i].Value = value
		return mo
	}

	mp := Pair[K, V]{key, value}
	*mo = append(*mo, mp)

	return mo
}

// Get retrieves the value for the specified key, along with a boolean indicating whether the key
// was found in the ordered Map. This function is useful when you want to access the value
// associated with a key in the ordered Map, and also check if the key exists in the map.
//
// Parameters:
//
// - key K: The key to search for in the ordered Map.
//
// Returns:
//
// - V: The value associated with the specified key if found, or the zero value for the value type
// if the key is not found.
//
// - bool: A boolean value indicating whether the key was found in the ordered Map.
//
// Example usage:
//
//	value, found := mo.Get("some_key")
//
// Retrieves the value associated with the key "some_key" and checks if the key exists in the
// ordered Map.
func (mo MapOrd[K, V]) Get(key K) Option[V] {
	if i := mo.index(key); i != -1 {
		return Some(mo[i].Value)
	}

	return None[V]()
}

// GetOrSet returns the value for a key. If the key does not exist, it returns the default value
// instead and also sets the default value for the key in the ordered Map. This function is useful
// when you want to access the value associated with a key in the ordered Map, and if the key does
// not exist, you want to return a specified default value and set that default value for the key.
//
// Parameters:
//
// - key K: The key to search for in the ordered Map.
//
// - defaultValue V: The default value to return if the key is not found in the ordered Map.
// If the key is not found, this default value will also be set for the key in the ordered Map.
//
// Returns:
//
// - V: The value associated with the specified key if found, or the provided default value if the key is not found.
//
// Example usage:
//
//	value := mo.GetOrSet("some_key", "default_value")
//
// Retrieves the value associated with the key "some_key" or returns "default_value" if the key is
// not found, and sets "default_value" as the value for "some_key" in the ordered Map if it's not
// present.
func (mo *MapOrd[K, V]) GetOrSet(key K, defaultValue V) V {
	if value := mo.Get(key); value.IsSome() {
		return value.Some()
	}

	mo.Set(key, defaultValue)

	return defaultValue
}

// Invert inverts the key-value pairs in the ordered Map, creating a new ordered Map with the
// values as keys and the original keys as values.
func (mo MapOrd[K, V]) Invert() MapOrd[V, K] {
	result := NewMapOrd[V, K](mo.Len())
	mo.Iter().ForEach(func(k K, v V) { result.Set(v, k) })

	return result
}

func (mo MapOrd[K, V]) index(key K) int {
	for i, mp := range mo {
		if reflect.DeepEqual(mp.Key, key) {
			return i
		}
	}

	return -1
}

// Keys returns an Slice containing all the keys in the ordered Map.
func (mo MapOrd[K, V]) Keys() Slice[K] { return mo.Iter().Keys().Collect() }

// Values returns an Slice containing all the values in the ordered Map.
func (mo MapOrd[K, V]) Values() Slice[V] { return mo.Iter().Values().Collect() }

// Delete removes the specified keys from the ordered Map.
func (mo *MapOrd[K, V]) Delete(keys ...K) MapOrd[K, V] {
	for _, key := range keys {
		if i := mo.index(key); i != -1 {
			*mo = append((*mo)[:i], (*mo)[i+1:]...)
		}
	}

	return *mo
}

// Eq compares the current ordered Map to another ordered Map and returns true if they are equal.
func (mo MapOrd[K, V]) Eq(other MapOrd[K, V]) bool {
	if mo.Len() != other.Len() {
		return false
	}

	for _, mp := range mo {
		value := other.Get(mp.Key)
		if value.IsNone() || !reflect.DeepEqual(value.Some(), mp.Value) {
			return false
		}
	}

	return true
}

// String returns a string representation of the ordered Map.
func (mo MapOrd[K, V]) String() string {
	var builder strings.Builder

	mo.Iter().ForEach(func(k K, v V) { builder.WriteString(fmt.Sprintf("%v:%v, ", k, v)) })

	return String(builder.String()).TrimRight(", ").Format("MapOrd{%s}").Std()
}

// Clear removes all key-value pairs from the ordered Map.
func (mo *MapOrd[K, V]) Clear() MapOrd[K, V] { return mo.Delete(mo.Keys()...) }

// Contains checks if the ordered Map contains the specified key.
func (mo MapOrd[K, V]) Contains(key K) bool { return mo.index(key) >= 0 }

// Empty checks if the ordered Map is empty.
func (mo MapOrd[K, V]) Empty() bool { return mo.Len() == 0 }

// Len returns the number of key-value pairs in the ordered Map.
func (mo MapOrd[K, V]) Len() int { return len(mo) }

// Ne compares the current ordered Map to another ordered Map and returns true if they are not
// equal.
func (mo MapOrd[K, V]) Ne(other MapOrd[K, V]) bool { return !mo.Eq(other) }

// NotEmpty checks if the ordered Map is not empty.
func (mo MapOrd[K, V]) NotEmpty() bool { return !mo.Empty() }

// Print prints the key-value pairs of the MapOrd to the standard output (console)
// and returns the MapOrd pointer unchanged.
func (mo MapOrd[K, V]) Print() MapOrd[K, V] { fmt.Println(mo); return mo }
