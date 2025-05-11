package g

// MapOrdEntry provides a view into a single key of an ordered Map (MapOrd),
// enabling fluent insertion, mutation, and deletion while preserving entry order.
//
// Example:
//
//	mo := g.NewMapOrd[string,int]()
//	// Insert 1 if "foo" absent, then double it
//	mo.Entry("foo").OrSet(1).AndModify(func(v *int){ *v *= 2 })
//
//	// Lazily set and then retrieve as Option
//	opt := mo.Entry("bar").OrSetBy(func() int { return expensive() }).Get()
//	if opt.IsSome() {
//	    fmt.Println(opt.Unwrap())
//	}
//
// All methods (except Get and Occupied) return the MapOrdEntry for chaining.
type MapOrdEntry[K, V any] struct {
	mo  *MapOrd[K, V]
	key K
}

// Entry returns a MapOrdEntry for the given key, allowing fine-grained control
// over the value in the ordered Map. Use with a pointer receiver to mutate.
func (mo *MapOrd[K, V]) Entry(key K) MapOrdEntry[K, V] { return MapOrdEntry[K, V]{mo, key} }

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

// AndModify applies fn to the existing value (if any) and returns the entry.
func (e MapOrdEntry[K, V]) AndModify(fn func(*V)) MapOrdEntry[K, V] {
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
