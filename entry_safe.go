package g

import "github.com/enetx/g/ref"

// OccupiedSafeEntry represents a view into a concurrent map entry that is known
// to be present.
//
// It is typically obtained from MapSafe.Entry(key) when the key exists.
// All operations on OccupiedSafeEntry are safe for concurrent use.
type OccupiedSafeEntry[K comparable, V any] struct {
	m   *MapSafe[K, V]
	key K
}

// sealed prevents external implementations of the SafeEntry interface.
func (OccupiedSafeEntry[K, V]) sealed() {}

// Key returns the key of this entry.
func (e OccupiedSafeEntry[K, V]) Key() K { return e.key }

// Get returns the current value associated with the key.
//
// If the key is concurrently removed, the zero value of V is returned.
func (e OccupiedSafeEntry[K, V]) Get() V {
	if actual, ok := e.m.data.Load(e.key); ok {
		return *(actual.(*V))
	}
	var zero V
	return zero
}

// Insert replaces the value associated with the key and returns the previous value.
//
// The replacement is performed atomically with respect to other map operations.
func (e OccupiedSafeEntry[K, V]) Insert(value V) V {
	old, _ := e.m.data.Swap(e.key, &value)
	return *(old.(*V))
}

// Remove removes the entry from the map and returns the previously stored value.
//
// If the key is concurrently removed, the zero value of V is returned.
func (e OccupiedSafeEntry[K, V]) Remove() V {
	if actual, loaded := e.m.data.LoadAndDelete(e.key); loaded {
		return *(actual.(*V))
	}
	var zero V
	return zero
}

// OrInsert returns the existing value without modifying the map.
func (e OccupiedSafeEntry[K, V]) OrInsert(value V) V { return e.Get() }

// OrInsertWith returns the existing value without invoking the function.
func (e OccupiedSafeEntry[K, V]) OrInsertWith(fn func() V) V { return e.Get() }

// OrInsertWithKey returns the existing value without invoking the function.
func (e OccupiedSafeEntry[K, V]) OrInsertWithKey(fn func(K) V) V { return e.Get() }

// OrDefault returns the existing value.
func (e OccupiedSafeEntry[K, V]) OrDefault() V { return e.Get() }

// AndModify applies the provided function to the value associated with the key
// and returns the entry.
//
// The modification is performed using a compare-and-swap loop.
// The function receives a pointer to a copy of the value; the updated value
// is written back atomically.
//
// If the key is concurrently removed, AndModify becomes a no-op.
func (e OccupiedSafeEntry[K, V]) AndModify(fn func(*V)) SafeEntry[K, V] {
	for {
		actual, ok := e.m.data.Load(e.key)
		if !ok {
			return e
		}

		oldPtr := actual.(*V)
		newVal := *oldPtr
		fn(&newVal)

		if e.m.data.CompareAndSwap(e.key, oldPtr, &newVal) {
			return e
		}
	}
}

// VacantSafeEntry represents a view into a concurrent map entry that is known
// to be absent at the time of creation.
//
// It is typically obtained from MapSafe.Entry(key) when the key does not exist.
// All operations on VacantSafeEntry are safe for concurrent use.
//
// The modify field stores a pending modification function from AndModify,
// which will be applied if OrInsert loses a race with another goroutine.
type VacantSafeEntry[K comparable, V any] struct {
	m      *MapSafe[K, V]
	key    K
	modify func(*V)
}

// sealed prevents external implementations of the SafeEntry interface.
func (VacantSafeEntry[K, V]) sealed() {}

// Key returns the key that would be used for insertion.
func (e VacantSafeEntry[K, V]) Key() K { return e.key }

// Insert inserts the provided value into the map and returns the stored value.
//
// If another goroutine inserts the same key concurrently, the existing value
// is returned instead.
func (e VacantSafeEntry[K, V]) Insert(value V) V {
	actual, _ := e.m.data.LoadOrStore(e.key, &value)
	return *(actual.(*V))
}

// OrInsert inserts the provided value and returns the stored value.
//
// If another goroutine inserts the same key concurrently (the insert "loses
// the race"), the existing value is used instead. In this case, if a pending
// modification was registered via AndModify, it is applied atomically to the
// existing value before returning.
//
// This ensures that chained calls like Entry(k).AndModify(f).OrInsert(v)
// behave correctly under concurrent access: the modification is never lost.
func (e VacantSafeEntry[K, V]) OrInsert(value V) V {
	actual, loaded := e.m.data.LoadOrStore(e.key, &value)
	if loaded && e.modify != nil {
		OccupiedSafeEntry[K, V]{m: e.m, key: e.key}.AndModify(e.modify)
		return e.m.Get(e.key).Some()
	}

	return *(actual.(*V))
}

// OrInsertWith inserts the value returned by the function and returns the
// stored value.
//
// If another goroutine inserts the key concurrently, the function may not be
// invoked and the existing value is returned.
func (e VacantSafeEntry[K, V]) OrInsertWith(fn func() V) V {
	if actual, ok := e.m.data.Load(e.key); ok {
		return *(actual.(*V))
	}

	actual, _ := e.m.data.LoadOrStore(e.key, ref.Of(fn()))
	return *(actual.(*V))
}

// OrInsertWithKey inserts the value returned by the function and returns the
// stored value.
//
// If another goroutine inserts the key concurrently, the function may not be
// invoked and the existing value is returned.
func (e VacantSafeEntry[K, V]) OrInsertWithKey(fn func(K) V) V {
	if actual, ok := e.m.data.Load(e.key); ok {
		return *(actual.(*V))
	}

	actual, _ := e.m.data.LoadOrStore(e.key, ref.Of(fn(e.key)))
	return *(actual.(*V))
}

// OrDefault inserts the zero value of V into the map and returns the stored value.
func (e VacantSafeEntry[K, V]) OrDefault() V {
	var zero V
	return e.Insert(zero)
}

// AndModify registers a modification function to be applied to the value.
//
// If the key was concurrently inserted by another goroutine since this
// VacantSafeEntry was created, the modification is applied immediately
// via OccupiedSafeEntry.AndModify.
//
// Otherwise, the function is stored and will be applied later by OrInsert
// if it loses the race to insert the key. This ensures that the pattern
// Entry(k).AndModify(f).OrInsert(v) correctly increments existing values
// even under heavy concurrent access.
//
// Returns the appropriate SafeEntry for method chaining.
func (e VacantSafeEntry[K, V]) AndModify(fn func(*V)) SafeEntry[K, V] {
	if _, ok := e.m.data.Load(e.key); ok {
		return OccupiedSafeEntry[K, V]{m: e.m, key: e.key}.AndModify(fn)
	}

	return VacantSafeEntry[K, V]{m: e.m, key: e.key, modify: fn}
}
