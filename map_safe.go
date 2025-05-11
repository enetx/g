package g

import (
	"fmt"

	"github.com/enetx/g/f"
)

// NewMapSafe creates a new instance of MapSafe with an optional initial size.
func NewMapSafe[K comparable, V any]() *MapSafe[K, V] { return &MapSafe[K, V]{} }

// Iter provides a thread-safe iterator over the MapSafe's key-value pairs.
func (ms *MapSafe[K, V]) Iter() SeqMap[K, V] {
	return func(yield func(K, V) bool) {
		ms.data.Range(func(key, value any) bool {
			k := key.(K)
			v := value.(V)
			return yield(k, v)
		})
	}
}

// IntoIter returns a consuming iterator (SeqMap[K, V]) over the MapSafe's key-value pairs.
// The iterator transfers ownership by removing the elements from the underlying map
// as they are iterated over. After iteration, the map will be empty.
//
// Returns:
//
// A SeqMap[K, V] that yields key-value pairs and removes them from the MapSafe.
//
// Example:
//
//	m := g.NewMapSafe[string, int]()
//	m.Set("a", 1)
//	m.Set("b", 2)
//	m.IntoIter().ForEach(func(k string, v int) {
//		fmt.Println(k, v)
//	})
//	m.Len() // Output: 0
func (ms *MapSafe[K, V]) IntoIter() SeqMap[K, V] {
	return func(yield func(K, V) bool) {
		var keys []K
		ms.data.Range(func(key, _ any) bool {
			keys = append(keys, key.(K))
			return true
		})

		for _, k := range keys {
			if val := ms.GetAndDelete(k); val.IsSome() {
				if !yield(k, val.Some()) {
					return
				}
			}
		}
	}
}

// Keys returns a slice of the MapSafe's keys.
func (ms *MapSafe[K, V]) Keys() Slice[K] { return ms.Iter().Keys().Collect() }

// Values returns a slice of the MapSafe's values.
func (ms *MapSafe[K, V]) Values() Slice[V] { return ms.Iter().Values().Collect() }

// Invert inverts the keys and values of the MapSafe, returning a new MapSafe with values as keys and
// keys as values. Note that the inverted Map will have 'any' as the key type, since not all value
// types are guaranteed to be comparable.
func (ms *MapSafe[K, V]) Invert() *MapSafe[any, K] {
	res := NewMapSafe[any, K]()
	ms.data.Range(func(key, value any) bool {
		res.data.Store(value, key.(K))
		return true
	})

	return res
}

// Contains checks if the MapSafe contains the specified key.
func (ms *MapSafe[K, V]) Contains(key K) bool {
	_, ok := ms.data.Load(key)
	return ok
}

// Clone creates a new MapSafe that is a copy of the original MapSafe.
func (ms *MapSafe[K, V]) Clone() *MapSafe[K, V] {
	res := NewMapSafe[K, V]()
	ms.data.Range(func(key, value any) bool {
		res.data.Store(key, value)
		return true
	})

	return res
}

// Copy copies the source MapSafe's key-value pairs to the target MapSafe.
func (ms *MapSafe[K, V]) Copy(src *MapSafe[K, V]) {
	src.data.Range(func(key, value any) bool {
		ms.data.Store(key, value)
		return true
	})
}

// Delete removes the specified keys from the MapSafe.
func (ms *MapSafe[K, V]) Delete(keys ...K) {
	for _, k := range keys {
		ms.data.Delete(k)
	}
}

// Eq checks if two MapSafes are equal.
func (ms *MapSafe[K, V]) Eq(other *MapSafe[K, V]) bool {
	n := ms.Len()
	if n != other.Len() {
		return false
	}

	if n == 0 {
		return true
	}

	res := true

	key := ms.Iter().Take(1).Keys().Collect()[0]
	comparable := f.IsComparable(ms.Get(key).Some())

	ms.data.Range(func(key, value any) bool {
		ov, ok := other.data.Load(key)
		if !ok || (comparable && !f.Eq(value)(ov)) || (!comparable && !f.Eqd(value)(ov)) {
			res = false
			return false
		}

		return true
	})

	return res
}

// Get retrieves the value associated with the given key.
func (ms *MapSafe[K, V]) Get(key K) Option[V] {
	if v, ok := ms.data.Load(key); ok {
		return Some(v.(V))
	}

	return None[V]()
}

// Set sets the value for the given key in the MapSafe.
func (ms *MapSafe[K, V]) Set(key K, value V) { ms.data.Store(key, value) }

// Len returns the number of key-value pairs in the MapSafe.
func (ms *MapSafe[K, V]) Len() Int {
	count := 0
	ms.data.Range(func(_, _ any) bool {
		count++
		return true
	})

	return Int(count)
}

// Ne checks if two MapSafes are not equal.
func (ms *MapSafe[K, V]) Ne(other *MapSafe[K, V]) bool { return !ms.Eq(other) }

// NotEmpty checks if the MapSafe is not empty.
func (ms *MapSafe[K, V]) NotEmpty() bool { return !ms.Empty() }

// GetOrSet retrieves the value for a key, or sets it to a default value if the key does not exist.
func (ms *MapSafe[K, V]) GetOrSet(key K, value V) (V, bool) {
	actual, loaded := ms.data.LoadOrStore(key, value)
	return actual.(V), loaded
}

// GetAndSet atomically sets a new value for the given key and returns the previous value, if any.
//
// Returns:
//   - Some(previous) if the key was present before the update.
//   - None if the key did not exist.
func (ms *MapSafe[K, V]) GetAndSet(key K, value V) Option[V] {
	if previous, loaded := ms.data.Swap(key, value); loaded {
		return Some(previous.(V))
	}

	return None[V]()
}

// GetAndDelete atomically retrieves and removes the value for the given key.
//
// Returns:
//   - Some(value) if the key existed and was removed.
//   - None if the key was not present.
func (ms *MapSafe[K, V]) GetAndDelete(key K) Option[V] {
	if value, loaded := ms.data.LoadAndDelete(key); loaded {
		return Some(value.(V))
	}

	return None[V]()
}

// Clear removes all key-value pairs from the MapSafe.
func (ms *MapSafe[K, V]) Clear() { ms.data.Clear() }

// Empty checks if the MapSafe is empty.
func (ms *MapSafe[K, V]) Empty() bool { return ms.Len() == 0 }

// String returns a string representation of the MapSafe.
func (ms *MapSafe[K, V]) String() string {
	builder := NewBuilder()
	ms.data.Range(func(key, value any) bool {
		builder.Write(Format("{}:{}, ", key, value))
		return true
	})

	return builder.String().StripSuffix(", ").Format("MapSafe\\{{}\\}").Std()
}

// Print writes the key-value pairs of the MapSafe to the standard output (console)
// and returns the MapSafe unchanged.
func (ms *MapSafe[K, V]) Print() *MapSafe[K, V] {
	fmt.Print(ms)
	return ms
}

// Println writes the key-value pairs of the MapSafe to the standard output (console) with a newline
// and returns the MapSafe unchanged.
func (ms *MapSafe[K, V]) Println() *MapSafe[K, V] {
	fmt.Println(ms)
	return ms
}
