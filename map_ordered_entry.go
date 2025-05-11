package g

// Get returns Some(value) if present, otherwise None.
func (e MapOrdEntry[K, V]) Get() Option[V] { return e.mo.Get(e.key) }

// OrSet inserts val if the key is vacant and returns the entry.
func (e MapOrdEntry[K, V]) OrSet(value V) MapOrdEntry[K, V] {
	if !e.mo.Contains(e.key) {
		e.mo.Set(e.key, value)
	}

	return e
}

// OrSetBy inserts the value produced by fn if the key is vacant and returns the entry.
func (e MapOrdEntry[K, V]) OrSetBy(fn func() V) MapOrdEntry[K, V] {
	if !e.mo.Contains(e.key) {
		e.mo.Set(e.key, fn())
	}

	return e
}

// OrDefault inserts V's zero value if the key is vacant and returns the entry.
func (e MapOrdEntry[K, V]) OrDefault() MapOrdEntry[K, V] {
	if !e.mo.Contains(e.key) {
		var zero V
		e.mo.Set(e.key, zero)
	}

	return e
}

// Transform applies fn to the existing value (if any) and returns the entry.
func (e MapOrdEntry[K, V]) Transform(fn func(*V)) MapOrdEntry[K, V] {
	if i := e.mo.index(e.key); i != -1 {
		value := (*e.mo)[i].Value
		fn(&value)
		(*e.mo)[i].Value = value
	}

	return e
}

// Set unconditionally sets the value for the key and returns the entry.
func (e MapOrdEntry[K, V]) Set(value V) MapOrdEntry[K, V] {
	e.mo.Set(e.key, value)
	return e
}

// Delete removes the key from the ordered Map and returns the entry.
func (e MapOrdEntry[K, V]) Delete() MapOrdEntry[K, V] {
	e.mo.Delete(e.key)
	return e
}
