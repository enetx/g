package g

import (
	"reflect"

	"github.com/enetx/g/f"
	"github.com/enetx/iter"
)

// isValueComparable reports whether values of type V can be compared with ==
// (V is a comparable type and not the bare interface any). It is the shared
// guard used by Dedup/Unique across the container and iterator types.
func isValueComparable[V any]() bool {
	return f.IsComparable[V]() && reflect.TypeFor[V]().Kind() != reflect.Interface
}

// counterBy tallies how many elements map to each key produced by fn, yielding
// the tally as a MapOrd sequence in first-seen key order. It backs the CounterBy
// method on the value-yielding Seq types.
func counterBy[V any, K comparable](seq iter.Seq[V], fn func(V) K) SeqMapOrd[K, Int] {
	return func(yield func(K, Int) bool) {
		order := NewSlice[K]()
		counts := NewMap[K, Int]()

		seq(func(v V) bool {
			k := fn(v)
			if !counts.Contains(k) {
				order.Push(k)
			}

			counts[k]++

			return true
		})

		for _, k := range order {
			if !yield(k, counts[k]) {
				return
			}
		}
	}
}

// flattenValue recursively descends slices and arrays within item, emitting each
// leaf element assignable to V. It stops and returns false as soon as emit
// returns false. It backs the Flatten method across the Seq types.
func flattenValue[V any](item any, emit func(V) bool) bool {
	rv := reflect.ValueOf(item)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		for i := range rv.Len() {
			if !flattenValue(rv.Index(i).Interface(), emit) {
				return false
			}
		}
	default:
		if v, ok := item.(V); ok {
			if !emit(v) {
				return false
			}
		}
	}

	return true
}
