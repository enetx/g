package g

// Get returns Some(value) if present, otherwise None.
func (e MapOrdEntry[K, V]) Get() Option[V] {
	return e.mo.Get(e.key)
}

// OrSet inserts value if the key is vacant.
// Returns Some(existing_value) if key was present, None otherwise.
func (e MapOrdEntry[K, V]) OrSet(value V) Option[V] {
	if i, exists := e.mo.lookup[e.key]; exists {
		return Some(e.mo.pairs[i].Value)
	}

	e.mo.Set(e.key, value)

	return None[V]()
}

// OrSetBy inserts the value produced by fn if the key is vacant.
// Returns Some(existing_value) if key was present, None otherwise.
func (e MapOrdEntry[K, V]) OrSetBy(fn func() V) Option[V] {
	if i, exists := e.mo.lookup[e.key]; exists {
		return Some(e.mo.pairs[i].Value)
	}

	e.mo.Set(e.key, fn())

	return None[V]()
}

// OrDefault inserts V's zero value if the key is vacant.
// Returns Some(existing_value) if key was present, None otherwise.
func (e MapOrdEntry[K, V]) OrDefault() Option[V] {
	var zero V
	return e.OrSet(zero)
}

// Transform applies fn to the value if it exists.
// Returns Some(updated_value) if key was present, None otherwise.
func (e MapOrdEntry[K, V]) Transform(fn func(V) V) Option[V] {
	if i, exists := e.mo.lookup[e.key]; exists {
		value := fn(e.mo.pairs[i].Value)
		e.mo.pairs[i].Value = value

		return Some(value)
	}

	return None[V]()
}

// Set sets the value for the specified key in the ordered map.
// Returns Some(previous_value) if the key existed, or None if it was newly inserted.
func (e MapOrdEntry[K, V]) Set(value V) Option[V] {
	return e.mo.Set(e.key, value)
}

// Delete removes the key from the Map.
// Returns Some(removed_value) if present, None otherwise.
func (e MapOrdEntry[K, V]) Delete() Option[V] {
	if i, exists := e.mo.lookup[e.key]; exists {
		value := e.mo.pairs[i].Value

		delete(e.mo.lookup, e.key)
		e.mo.pairs = append(e.mo.pairs[:i], e.mo.pairs[i+1:]...)

		for j := i; j < len(e.mo.pairs); j++ {
			e.mo.lookup[e.mo.pairs[j].Key] = j
		}

		return Some(value)
	}

	return None[V]()
}
