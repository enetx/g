package g

import (
	"maps"
)

// NewMapSafe creates a new instance of MapSafe with an optional initial size.
func NewMapSafe[K comparable, V any](size ...Int) *MapSafe[K, V] {
	return &MapSafe[K, V]{data: NewMap[K, V](size...)}
}

// Iter provides a thread-safe iterator over the MapSafe's key-value pairs.
func (ms *MapSafe[K, V]) Iter() SeqMap[K, V] {
	ms.mu.RLock()
	keys := maps.Keys(ms.data)
	ms.mu.RUnlock()

	return func(yield func(K, V) bool) {
		for k := range keys {
			if v := ms.Get(k); v.IsSome() {
				if !yield(k, v.Some()) {
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
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	result := NewMapSafe[any, K](ms.Len())
	for k, v := range ms.data {
		result.data.Set(v, k)
	}

	return result
}

// Contains checks if the MapSafe contains the specified key.
func (ms *MapSafe[K, V]) Contains(key K) bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	return ms.data.Contains(key)
}

// Clone creates a new MapSafe that is a copy of the original MapSafe.
func (ms *MapSafe[K, V]) Clone() *MapSafe[K, V] {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	result := NewMapSafe[K, V]()
	result.data = maps.Clone(ms.data)

	return result
}

// Copy copies the source MapSafe's key-value pairs to the target MapSafe.
func (ms *MapSafe[K, V]) Copy(src *MapSafe[K, V]) {
	src.mu.RLock()
	defer src.mu.RUnlock()

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.data.Copy(src.data)
}

// Delete removes the specified keys from the MapSafe.
func (ms *MapSafe[K, V]) Delete(keys ...K) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.data.Delete(keys...)
}

// Eq checks if two MapSafes are equal.
func (ms *MapSafe[K, V]) Eq(other *MapSafe[K, V]) bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	other.mu.RLock()
	defer other.mu.RUnlock()

	return ms.data.Eq(other.data)
}

// Get retrieves the value associated with the given key.
func (ms *MapSafe[K, V]) Get(key K) Option[V] {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	return ms.data.Get(key)
}

// Set sets the value for the given key in the MapSafe.
func (ms *MapSafe[K, V]) Set(key K, value V) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.data.Set(key, value)
}

// Len returns the number of key-value pairs in the MapSafe.
func (ms *MapSafe[K, V]) Len() Int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	return ms.data.Len()
}

// Ne checks if two MapSafes are not equal.
func (ms *MapSafe[K, V]) Ne(other *MapSafe[K, V]) bool { return !ms.Eq(other) }

// NotEmpty checks if the MapSafe is not empty.
func (ms *MapSafe[K, V]) NotEmpty() bool { return !ms.Empty() }

// GetOrSet retrieves the value for a key, or sets it to a default value if the key does not exist.
func (ms *MapSafe[K, V]) GetOrSet(key K, defaultValue V) V {
	ms.mu.RLock()
	value, ok := ms.data[key]
	ms.mu.RUnlock()

	if ok {
		return value
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.data[key] = defaultValue

	return defaultValue
}

// Clear removes all key-value pairs from the MapSafe.
func (ms *MapSafe[K, V]) Clear() {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	clear(ms.data)
}

// Empty checks if the MapSafe is empty.
func (ms *MapSafe[K, V]) Empty() bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	return ms.data.Empty()
}

// String returns a string representation of the MapSafe.
func (ms *MapSafe[K, V]) String() string {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	builder := NewBuilder()

	for k, v := range ms.data {
		builder.Write(Sprintf("{}:{}, ", k, v))
	}

	return builder.String().StripSuffix(", ").Format("MapSafe\\{{}\\}").Std()
}

// Print writes the key-value pairs of the MapSafe to the standard output (console)
// and returns the MapSafe unchanged.
func (ms *MapSafe[K, V]) Print() *MapSafe[K, V] {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	Print(ms)
	return ms
}

// Println writes the key-value pairs of the MapSafe to the standard output (console) with a newline
// and returns the MapSafe unchanged.
func (ms *MapSafe[K, V]) Println() *MapSafe[K, V] {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	Println(ms)
	return ms
}
