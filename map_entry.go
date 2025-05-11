package g

// Get returns Some(value) if the key exists, otherwise None.
func (e MapEntry[K, V]) Get() Option[V] { return e.m.Get(e.key) }

// OrSet inserts value if the key is vacant and returns the Entry for chaining.
func (e MapEntry[K, V]) OrSet(value V) MapEntry[K, V] {
	if _, ok := e.m[e.key]; !ok {
		e.m[e.key] = value
	}

	return e
}

// OrSetBy inserts the value produced by `fn` if the key is vacant and returns the Entry.
func (e MapEntry[K, V]) OrSetBy(fn func() V) MapEntry[K, V] {
	if _, ok := e.m[e.key]; !ok {
		e.m[e.key] = fn()
	}

	return e
}

// OrDefault inserts V's zero value if the key is vacant and returns the Entry.
func (e MapEntry[K, V]) OrDefault() MapEntry[K, V] {
	if _, ok := e.m[e.key]; !ok {
		var zero V
		e.m[e.key] = zero
	}

	return e
}

// Transform applies `fn` to the existing value if present, and returns the Entry.
// If the key is vacant, it does nothing.
func (e MapEntry[K, V]) Transform(fn func(*V)) MapEntry[K, V] {
	if value, ok := e.m[e.key]; ok {
		fn(&value)
		e.m[e.key] = value
	}

	return e
}

// Set unconditionally sets the value for the key and returns the Entry.
func (e MapEntry[K, V]) Set(value V) MapEntry[K, V] {
	e.m[e.key] = value
	return e
}

// Delete removes the key from the Map and returns the Entry.
func (e MapEntry[K, V]) Delete() MapEntry[K, V] {
	delete(e.m, e.key)
	return e
}
