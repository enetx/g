package g

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"slices"

	"github.com/enetx/g/cmp"
	"github.com/enetx/g/f"
	"github.com/enetx/iter"
)

// Pair is a struct representing a key-value Pair for MapOrd.
type Pair[K, V any] = iter.Pair[K, V]

// MapOrd is an ordered map that maintains insertion order using a slice for pairs
// and a map for fast index lookups.
type MapOrd[K comparable, V any] []Pair[K, V] // ordered key-value pairs

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
func NewMapOrd[K comparable, V any](size ...Int) MapOrd[K, V] {
	if len(size) > 0 {
		return make(MapOrd[K, V], 0, size[0])
	}

	return make(MapOrd[K, V], 0)
}

// Transform applies a transformation function to the MapOrd and returns the result.
func (mo MapOrd[K, V]) Transform(fn func(MapOrd[K, V]) MapOrd[K, V]) MapOrd[K, V] { return fn(mo) }

// Entry returns an OrdEntry for the given key.
func (mo *MapOrd[K, V]) Entry(key K) OrdEntry[K, V] {
	if i := mo.index(key); i != -1 {
		return OccupiedOrdEntry[K, V]{mo: mo, key: key, index: i}
	}

	return VacantOrdEntry[K, V]{mo: mo, key: key}
}

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
func MapOrdFromStd[K comparable, V any](m map[K]V) MapOrd[K, V] { return Map[K, V](m).Ordered() }

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

// IsSortedBy checks if the ordered Map is sorted according to a custom comparison function.
//
// Parameters:
//
// - fn func(a, b Pair[K, V]) cmp.Ordering: The custom comparison function used for checking sort order.
//
// Returns:
//
// - bool: true if the map is sorted according to the comparison function, false otherwise.
//
// Example usage:
//
//	sorted := hmapo.IsSortedBy(func(a, b g.Pair[g.String, g.Int]) cmp.Ordering { return a.Key.Cmp(b.Key) })
func (mo MapOrd[K, V]) IsSortedBy(fn func(a, b Pair[K, V]) cmp.Ordering) bool {
	if len(mo) <= 1 {
		return true
	}

	for i := 1; i < len(mo); i++ {
		if fn(mo[i-1], mo[i]).IsGt() {
			return false
		}
	}

	return true
}

// IsSortedByKey checks if the ordered MapOrd[K, V] is sorted by the keys using a custom comparison function.
//
// Parameters:
//
// - fn func(a, b K) cmp.Ordering: The custom comparison function used for checking key sort order.
//
// Returns:
//
// - bool: true if the map is sorted by keys according to the comparison function, false otherwise.
//
// Example usage:
//
//	sorted := hmapo.IsSortedByKey(func(a, b g.String) cmp.Ordering { return a.Cmp(b) })
func (mo MapOrd[K, V]) IsSortedByKey(fn func(a, b K) cmp.Ordering) bool {
	if len(mo) <= 1 {
		return true
	}

	for i := 1; i < len(mo); i++ {
		if fn(mo[i-1].Key, mo[i].Key).IsGt() {
			return false
		}
	}

	return true
}

// IsSortedByValue checks if the ordered MapOrd[K, V] is sorted by the values using a custom comparison function.
//
// Parameters:
//
// - fn func(a, b V) cmp.Ordering: The custom comparison function used for checking value sort order.
//
// Returns:
//
// - bool: true if the map is sorted by values according to the comparison function, false otherwise.
//
// Example usage:
//
//	sorted := hmapo.IsSortedByValue(func(a, b g.Int) cmp.Ordering { return a.Cmp(b) })
func (mo MapOrd[K, V]) IsSortedByValue(fn func(a, b V) cmp.Ordering) bool {
	if len(mo) <= 1 {
		return true
	}

	for i := 1; i < len(mo); i++ {
		if fn(mo[i-1].Value, mo[i].Value).IsGt() {
			return false
		}
	}

	return true
}

// Clone creates a new ordered Map with the same key-value pairs.
func (mo MapOrd[K, V]) Clone() MapOrd[K, V] {
	nmo := NewMapOrd[K, V](mo.Len())
	nmo.Copy(mo)

	return nmo
}

// Copy copies key-value pairs from the source ordered Map to the current ordered Map.
func (mo *MapOrd[K, V]) Copy(src MapOrd[K, V]) {
	idx := mo.indexMap()

	for _, p := range src {
		if i, ok := idx[p.Key]; ok {
			(*mo)[i].Value = p.Value
		} else {
			*mo = append(*mo, p)
			idx[p.Key] = len(*mo) - 1
		}
	}
}

// Map converts the ordered Map to a standard Map.
func (mo MapOrd[K, V]) Map() Map[K, V] {
	m := NewMap[K, V](mo.Len())
	mo.Iter().ForEach(func(k K, v V) { m[k] = v })

	return m
}

// Safe converts a ordered Map to a thread-safe Map.
func (mo MapOrd[K, V]) Safe() *MapSafe[K, V] {
	ms := NewMapSafe[K, V]()
	mo.Iter().ForEach(func(k K, v V) { ms.Insert(k, v) })

	return ms
}

// Insert sets the value for the specified key in the ordered Map,
// and returns the previous value if it existed.
func (mo *MapOrd[K, V]) Insert(key K, value V) Option[V] {
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

func (mo MapOrd[K, V]) index(key K) int {
	for i, mp := range mo {
		if mp.Key == key {
			return i
		}
	}

	return -1
}

// Keys returns an Slice containing all the keys in the ordered Map.
func (mo MapOrd[K, V]) Keys() Slice[K] { return mo.Iter().Keys().Collect() }

// Values returns an Slice containing all the values in the ordered Map.
func (mo MapOrd[K, V]) Values() Slice[V] { return mo.Iter().Values().Collect() }

// Remove removes the specified key from the ordered Map and returns the removed value.
func (mo *MapOrd[K, V]) Remove(key K) Option[V] {
	if mo.IsEmpty() {
		return None[V]()
	}

	for i, p := range *mo {
		if p.Key == key {
			*mo = append((*mo)[:i], (*mo)[i+1:]...)
			return Some(p.Value)
		}
	}

	return None[V]()
}

// Eq compares the current ordered Map to another ordered Map and returns true if they are equal.
func (mo MapOrd[K, V]) Eq(other MapOrd[K, V]) bool {
	if len(mo) != len(other) {
		return false
	}
	if len(mo) == 0 {
		return true
	}

	idx := other.indexMap()

	comparable := f.IsComparable[V]()
	for i, mp := range mo {
		j, ok := idx[mp.Key]

		if !ok || j != i {
			return false
		}

		if comparable {
			if any(other[j].Value) != any(mp.Value) {
				return false
			}
		} else {
			if !reflect.DeepEqual(other[j].Value, mp.Value) {
				return false
			}
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
	b.Grow(Int(len(mo)) * 16)
	b.WriteString("MapOrd{")

	first := true
	for _, pair := range mo {
		if !first {
			b.WriteString(", ")
		}

		first = false
		fmt.Fprint(&b, pair.Key)
		b.WriteByte(':')
		fmt.Fprint(&b, pair.Value)
	}

	b.WriteString("}")

	return b.String().Std()
}

// Clear removes all key-value pairs from the ordered Map.
func (mo *MapOrd[K, V]) Clear() { *mo = (*mo)[:0] }

// Contains checks if the ordered Map contains the specified key.
func (mo MapOrd[K, V]) Contains(key K) bool { return mo.index(key) != -1 }

// Empty checks if the ordered Map is empty.
func (mo MapOrd[K, V]) IsEmpty() bool { return len(mo) == 0 }

// Len returns the number of key-value pairs in the ordered Map.
func (mo MapOrd[K, V]) Len() Int { return Int(len(mo)) }

// Ne compares the current ordered Map to another ordered Map and returns true if they are not equal.
func (mo MapOrd[K, V]) Ne(other MapOrd[K, V]) bool { return !mo.Eq(other) }

// Print writes the key-value pairs of the MapOrd to the standard output (console)
// and returns the MapOrd unchanged.
func (mo MapOrd[K, V]) Print() MapOrd[K, V] { fmt.Print(mo); return mo }

// Println writes the key-value pairs of the MapOrd to the standard output (console) with a newline
// and returns the MapOrd unchanged.
func (mo MapOrd[K, V]) Println() MapOrd[K, V] { fmt.Println(mo); return mo }

// indexMap builds a map from keys to their corresponding indices in the MapOrd.
//
// This function is used to create a temporary indexMap that maps each key in the
// ordered map to its position (insertion order) within the slice. It is useful
// for optimizing lookup operations such as Set, Delete, Copy, or Eq.
//
// Time complexity: O(n), where n is the number of key-value pairs in the MapOrd.
func (mo MapOrd[K, V]) indexMap() map[K]int {
	idx := make(map[K]int, len(mo))

	for i, p := range mo {
		idx[p.Key] = i
	}

	return idx
}
