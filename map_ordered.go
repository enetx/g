package g

import (
	"fmt"
	"slices"

	"github.com/enetx/g/cmp"
	"github.com/enetx/g/f"
	"github.com/enetx/g/rand"
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
func NewMapOrd[K, V any](size ...Int) MapOrd[K, V] {
	return make(MapOrd[K, V], 0, Slice[Int](size).Get(0).UnwrapOrDefault())
}

// Transform applies a transformation function to the MapOrd and returns the result.
func (mo MapOrd[K, V]) Transform(fn func(MapOrd[K, V]) MapOrd[K, V]) MapOrd[K, V] { return fn(mo) }

// AsAny converts all key-value pairs in the MapOrd to type `any`.
//
// It returns a new MapOrd[any, any], where both keys and values are of type `any`.
// This is useful when working with dynamic formatting tools like Println or Format,
// which can access map elements dynamically when keys and values are of type `any`.
//
// Example:
//
//	mo := NewMapOrd[string, int]()
//	mo.Set("a", 1)
//	mo.Set("b", 2)
//
//	Println("{1.a} -> {1.b}", mo.AsAny())
//	// Output: "a -> 1"
func (mo MapOrd[K, V]) AsAny() MapOrd[any, any] {
	if mo.Empty() {
		return NewMapOrd[any, any]()
	}

	anymo := make(MapOrd[any, any], len(mo))
	for i, v := range mo {
		anymo[i] = Pair[any, any]{v.Key, v.Value}
	}

	return anymo
}

// Entry returns a MapOrdEntry object for the given key, providing fine-grained
// control over insertion, mutation, and deletion of its value in the ordered Map,
// while preserving the insertion order.
//
// Example:
//
//	mo := g.NewMapOrd[string, int]()
//	// Insert 1 if "foo" is absent, then increment it
//	e := mo.Entry("foo")
//	e.OrSet(1).
//	e.Transform(func(v int) int { return v + 1 })
//
// The entire operation requires only a single key lookup and works without
// additional allocations.
func (mo *MapOrd[K, V]) Entry(key K) MapOrdEntry[K, V] { return MapOrdEntry[K, V]{mo, key} }

// Iter returns an iterator (SeqMapOrd[K, V]) for the ordered Map, allowing for sequential iteration
// over its key-value pairs. It is commonly used in combination with higher-order functions,
// such as 'ForEach', to perform operations on each key-value pair of the ordered Map.
//
// Returns:
//
// A SeqMapOrd[K, V], which can be used for sequential iteration over the key-value pairs of the ordered Map.
//
// Example usage:
//
//	m := g.NewMapOrd[int, int]()
//	m.Set(1, 1)
//	m.Set(2, 2)
//	m.Set(3, 3).
//
//	m.Iter().ForEach(func(k, v int) {
//	    // Process key-value pair
//	})
//
// The 'Iter' method provides a convenient way to traverse the key-value pairs of an ordered Map
// in a functional style, enabling operations like mapping or filtering.
func (mo MapOrd[K, V]) Iter() SeqMapOrd[K, V] {
	return func(yield func(K, V) bool) {
		for _, v := range mo {
			if !yield(v.Key, v.Value) {
				return
			}
		}
	}
}

// IterReverse returns an iterator (SeqMapOrd[K, V]) for the ordered Map that allows for sequential iteration
// over its key-value pairs in reverse order. This method is useful when you need to process the elements
// from the last to the first.
//
// Returns:
//
// A SeqMapOrd[K, V], which can be used for sequential iteration over the key-value pairs of the ordered Map in reverse order.
//
// Example usage:
//
//	m := g.NewMapOrd[int, int]()
//	m.Set(1, 1)
//	m.Set(2, 2)
//	m.Set(3, 3)
//
//	m.IterReverse().ForEach(func(k, v int) {
//	    // Process key-value pair in reverse order
//	    fmt.Println("Key:", k, "Value:", v)
//	})
//
// The 'IterReverse' method complements the 'Iter' method by providing a way to access the elements
// in a reverse sequence, offering additional flexibility in data processing scenarios.
func (mo MapOrd[K, V]) IterReverse() SeqMapOrd[K, V] {
	return func(yield func(K, V) bool) {
		for i := len(mo) - 1; i >= 0; i-- {
			v := mo[i]
			if !yield(v.Key, v.Value) {
				return
			}
		}
	}
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
func MapOrdFromStd[K comparable, V any](m map[K]V) MapOrd[K, V] { return Map[K, V](m).ToMapOrd() }

// SortBy sorts the ordered Map by a custom comparison function.
//
// Parameters:
//
// - fn func(a, b Pair[K, V]) cmp.Ordering: The custom comparison function used for sorting the ordered Map.
//
// Example usage:
//
//	hmapo.SortBy(func(a, b g.Pair[g.String, g.Int]) cmp.Ordering { return a.Key.Cmp(b.Key) })
//	hmapo.SortBy(func(a, b g.Pair[g.String, g.Int]) cmp.Ordering { return a.Value.Cmp(b.Value) })
func (mo MapOrd[K, V]) SortBy(fn func(a, b Pair[K, V]) cmp.Ordering) {
	slices.SortFunc(mo, func(a, b Pair[K, V]) int { return int(fn(a, b)) })
}

// SortByKey sorts the ordered MapOrd[K, V] by the keys using a custom comparison function.
//
// Parameters:
//
// - fn func(a, b K) cmp.Ordering: The custom comparison function used for sorting the keys.
//
// Example usage:
//
//	hmapo.SortByKey(func(a, b g.String) cmp.Ordering { return a.Cmp(b) })
func (mo MapOrd[K, V]) SortByKey(fn func(a, b K) cmp.Ordering) {
	slices.SortFunc(mo, func(a, b Pair[K, V]) int { return int(fn(a.Key, b.Key)) })
}

// SortByValue sorts the ordered MapOrd[K, V] by the values using a custom comparison function.
//
// Parameters:
//
// - fn func(a, b V) cmp.Ordering: The custom comparison function used for sorting the values.
//
// Example usage:
//
//	hmapo.SortByValue(func(a, b g.Int) cmp.Ordering { return a.Cmp(b) })
func (mo MapOrd[K, V]) SortByValue(fn func(a, b V) cmp.Ordering) {
	slices.SortFunc(mo, func(a, b Pair[K, V]) int { return int(fn(a.Value, b.Value)) })
}

// Clone creates a new ordered Map with the same key-value pairs.
func (mo MapOrd[K, V]) Clone() MapOrd[K, V] {
	if mo.Empty() {
		return NewMapOrd[K, V]()
	}

	result := make(MapOrd[K, V], len(mo))
	copy(result, mo)

	return result
}

// Copy copies key-value pairs from the source ordered Map to the current ordered Map.
func (mo *MapOrd[K, V]) Copy(src MapOrd[K, V]) {
	if src.Empty() {
		return
	}

	for _, srcPair := range src {
		if i := mo.index(srcPair.Key); i != -1 {
			(*mo)[i].Value = srcPair.Value
		} else {
			*mo = append(*mo, srcPair)
		}
	}
}

// // ToMap converts the ordered Map to a standard Map.
// func (mo MapOrd[K, V]) ToMap() Map[K, V] {
// 	m := NewMap[K, V](len(mo))
// 	mo.Iter().ForEach(func(k K, v V) { m.Set(k, v) })
//
// 	return m
// }

// Set sets the value for the specified key in the ordered Map,
// and returns the previous value if it existed.
func (mo *MapOrd[K, V]) Set(key K, value V) Option[V] {
	if i := mo.index(key); i != -1 {
		prev := (*mo)[i].Value
		(*mo)[i].Value = value

		return Some(prev)
	}

	mp := Pair[K, V]{Key: key, Value: value}
	*mo = append(*mo, mp)

	return None[V]()
}

// Get returns the value associated with the given key, wrapped in Option[V].
//
// It returns Some(value) if the key exists, or None if it does not.
func (mo MapOrd[K, V]) Get(key K) Option[V] {
	if i := mo.index(key); i != -1 {
		return Some(mo[i].Value)
	}

	return None[V]()
}

// Shuffle randomly reorders the elements of the ordered Map.
// It operates in place and affects the original order of the map's entries.
//
// The function uses the crypto/rand package to generate random indices.
func (mo MapOrd[K, V]) Shuffle() {
	for i := mo.Len() - 1; i > 0; i-- {
		j := rand.N(i + 1)
		mo[i], mo[j] = mo[j], mo[i]
	}
}

// Invert inverts the key-value pairs in the ordered Map, creating a new ordered Map with the
// values as keys and the original keys as values.
func (mo MapOrd[K, V]) Invert() MapOrd[V, K] {
	if mo.Empty() {
		return NewMapOrd[V, K]()
	}

	result := make(MapOrd[V, K], 0, len(mo))
	for _, pair := range mo {
		result = append(result, Pair[V, K]{pair.Value, pair.Key})
	}

	return result
}

func (mo MapOrd[K, V]) index(key K) int {
	var zero K
	if f.IsComparable(zero) {
		for i, mp := range mo {
			if f.Eq[any](mp.Key)(key) {
				return i
			}
		}

		return -1
	}

	for i, mp := range mo {
		if f.Eqd(mp.Key)(key) {
			return i
		}
	}

	return -1
}

// Keys returns an Slice containing all the keys in the ordered Map.
func (mo MapOrd[K, V]) Keys() Slice[K] {
	if mo.Empty() {
		return NewSlice[K]()
	}

	keys := make(Slice[K], len(mo))
	for i, pair := range mo {
		keys[i] = pair.Key
	}

	return keys
}

// Values returns an Slice containing all the values in the ordered Map.
func (mo MapOrd[K, V]) Values() Slice[V] {
	if mo.Empty() {
		return NewSlice[V]()
	}

	values := make(Slice[V], len(mo))
	for i, pair := range mo {
		values[i] = pair.Value
	}

	return values
}

// Delete removes the specified keys from the ordered Map.
func (mo *MapOrd[K, V]) Delete(keys ...K) {
	if len(keys) == 0 || mo.Empty() {
		return
	}

	writePos := 0
	for readPos := 0; readPos < len(*mo); readPos++ {
		shouldDelete := false

		for _, delKey := range keys {
			var zero K
			if f.IsComparable(zero) {
				if f.Eq[any]((*mo)[readPos].Key)(delKey) {
					shouldDelete = true
					break
				}
			} else {
				if f.Eqd((*mo)[readPos].Key)(delKey) {
					shouldDelete = true
					break
				}
			}
		}

		if !shouldDelete {
			if writePos != readPos {
				(*mo)[writePos] = (*mo)[readPos]
			}
			writePos++
		}
	}

	*mo = (*mo)[:writePos]
}

// Eq compares the current ordered Map to another ordered Map and returns true if they are equal.
func (mo MapOrd[K, V]) Eq(other MapOrd[K, V]) bool {
	n := len(mo)

	if n != len(other) {
		return false
	}

	if n == 0 {
		return true
	}

	var zero V
	comparable := f.IsComparable(zero)

	for i, mp := range mo {
		if other.index(mp.Key) != i {
			return false
		}

		value := other[i].Value

		if comparable && !f.Eq[any](value)(mp.Value) || !comparable && !f.Eqd(value)(mp.Value) {
			return false
		}
	}

	return true
}

// String returns a string representation of the ordered Map.
func (mo MapOrd[K, V]) String() string {
	if len(mo) == 0 {
		return "MapOrd{}"
	}

	var b Builder
	b.WriteString("MapOrd{")

	first := true
	for _, pair := range mo {
		if !first {
			b.WriteString(", ")
		}

		first = false
		b.WriteString(Format("{}:{}", pair.Key, pair.Value))
	}

	b.WriteString("}")

	return b.String().Std()
}

// Clear removes all key-value pairs from the ordered Map.
func (mo *MapOrd[K, V]) Clear() { *mo = (*mo)[:0] }

// Contains checks if the ordered Map contains the specified key.
func (mo MapOrd[K, V]) Contains(key K) bool { return mo.index(key) >= 0 }

// Empty checks if the ordered Map is empty.
func (mo MapOrd[K, V]) Empty() bool { return len(mo) == 0 }

// Len returns the number of key-value pairs in the ordered Map.
func (mo MapOrd[K, V]) Len() Int { return Int(len(mo)) }

// Ne compares the current ordered Map to another ordered Map and returns true if they are not
// equal.
func (mo MapOrd[K, V]) Ne(other MapOrd[K, V]) bool { return !mo.Eq(other) }

// NotEmpty checks if the ordered Map is not empty.
func (mo MapOrd[K, V]) NotEmpty() bool { return !mo.Empty() }

// Print writes the key-value pairs of the MapOrd to the standard output (console)
// and returns the MapOrd unchanged.
func (mo MapOrd[K, V]) Print() MapOrd[K, V] { fmt.Print(mo); return mo }

// Println writes the key-value pairs of the MapOrd to the standard output (console) with a newline
// and returns the MapOrd unchanged.
func (mo MapOrd[K, V]) Println() MapOrd[K, V] { fmt.Println(mo); return mo }
