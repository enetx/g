package g

import "slices"

// OccupiedOrdEntry represents a view into an ordered map entry that is known
// to be present.
//
// It is typically obtained from MapOrd.Entry(key) when the key already exists.
// OccupiedOrdEntry provides access to the key and the value stored at a specific
// index in the underlying ordered slice, allowing inspection, modification,
// replacement, or removal without additional lookups.
type OccupiedOrdEntry[K comparable, V any] struct {
	mo    *MapOrd[K, V]
	index int
	key   K
}

// sealed prevents external implementations of the OrdEntry interface.
func (OccupiedOrdEntry[K, V]) sealed() {}

// Key returns the key of this occupied entry.
func (e OccupiedOrdEntry[K, V]) Key() K { return e.key }

// Get returns the current value associated with the key.
//
// The value is returned by copy.
func (e OccupiedOrdEntry[K, V]) Get() V { return (*e.mo)[e.index].Value }

// Insert replaces the value at the entry's position with the provided value
// and returns the previously stored value.
//
// The position of the entry in the ordered map is preserved.
func (e OccupiedOrdEntry[K, V]) Insert(value V) V {
	old := (*e.mo)[e.index].Value
	(*e.mo)[e.index].Value = value
	return old
}

// Remove removes the entry from the ordered map and returns the previously
// stored value.
//
// This operation preserves the relative order of the remaining entries.
// After this call, the key is no longer present in the map.
func (e OccupiedOrdEntry[K, V]) Remove() V {
	v := (*e.mo)[e.index].Value
	*e.mo = slices.Delete(*e.mo, e.index, e.index+1)
	return v
}

// OrInsert returns the existing value without modifying the map.
//
// For OccupiedOrdEntry, this is equivalent to Get since the key already exists.
func (e OccupiedOrdEntry[K, V]) OrInsert(value V) V { return e.Get() }

// OrInsertWith returns the existing value without invoking the function.
//
// For OccupiedOrdEntry, the function is never called since the key already exists.
func (e OccupiedOrdEntry[K, V]) OrInsertWith(fn func() V) V { return e.Get() }

// OrInsertWithKey returns the existing value without invoking the function.
//
// For OccupiedOrdEntry, the function is never called since the key already exists.
func (e OccupiedOrdEntry[K, V]) OrInsertWithKey(fn func(K) V) V { return e.Get() }

// OrDefault returns the existing value.
//
// For OccupiedOrdEntry, this is equivalent to Get since the key already exists.
func (e OccupiedOrdEntry[K, V]) OrDefault() V { return e.Get() }

// AndModify applies the provided function to the value stored at the entry's
// position and returns the entry for method chaining.
//
// The function receives a pointer to the actual value stored in the ordered map,
// allowing in-place modification.
//
// Example:
//
//	m.Entry("count").AndModify(func(v *int) { *v++ }).OrInsert(1)
func (e OccupiedOrdEntry[K, V]) AndModify(fn func(*V)) OrdEntry[K, V] {
	fn(&(*e.mo)[e.index].Value)
	return e
}

// VacantOrdEntry represents a view into an ordered map entry that is known
// to be absent.
//
// It is typically obtained from MapOrd.Entry(key) when the key does not exist.
// VacantOrdEntry allows inserting a new key-value pair into the ordered map.
type VacantOrdEntry[K comparable, V any] struct {
	mo  *MapOrd[K, V]
	key K
}

// sealed prevents external implementations of the OrdEntry interface.
func (VacantOrdEntry[K, V]) sealed() {}

// Key returns the key that would be used for insertion.
func (e VacantOrdEntry[K, V]) Key() K { return e.key }

// Insert inserts a new key-value pair into the ordered map and returns the value.
//
// The new entry is appended to the end of the ordered map.
// After this call, the key is present in the map with the given value.
func (e VacantOrdEntry[K, V]) Insert(value V) V {
	*e.mo = append(*e.mo, Pair[K, V]{Key: e.key, Value: value})
	return value
}

// OrInsert inserts the provided value and returns it.
//
// This is the primary method for inserting values via VacantOrdEntry.
func (e VacantOrdEntry[K, V]) OrInsert(value V) V { return e.Insert(value) }

// OrInsertWith inserts the value returned by the function and returns it.
//
// The function is guaranteed to be called exactly once.
// Use this when computing the default value is expensive.
func (e VacantOrdEntry[K, V]) OrInsertWith(fn func() V) V { return e.Insert(fn()) }

// OrInsertWithKey inserts the value returned by the function and returns it.
//
// The function receives the entry key and is guaranteed to be called exactly once.
// Use this when the default value depends on the key.
func (e VacantOrdEntry[K, V]) OrInsertWithKey(fn func(K) V) V {
	return e.Insert(fn(e.key))
}

// OrDefault inserts the zero value of V into the ordered map and returns it.
//
// This is useful for types where the zero value is a valid initial state,
// such as numeric types (0), slices (nil), or structs with zero defaults.
func (e VacantOrdEntry[K, V]) OrDefault() V {
	var zero V
	return e.Insert(zero)
}

// AndModify does nothing for VacantOrdEntry and returns the entry unchanged.
//
// Since there is no existing value to modify, the function is not called.
// This allows fluent chaining like Entry(k).AndModify(f).OrInsert(v)
// to work correctly regardless of whether the key exists.
func (e VacantOrdEntry[K, V]) AndModify(fn func(*V)) OrdEntry[K, V] { return e }
