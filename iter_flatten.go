package g

type flattenIter[T any] struct {
	baseIter[T]
	currentIter iterator[T]
	innerIter   *flattenIter[T]
}

func flatten[T any](iter iterator[T]) *flattenIter[T] {
	flattenIter := &flattenIter[T]{currentIter: iter}
	flattenIter.baseIter = baseIter[T]{flattenIter}

	return flattenIter
}

func (iter *flattenIter[T]) Next() Option[T] {
	for {
		if iter.innerIter != nil {
			if next := iter.innerIter.Next(); next.IsSome() {
				return next
			}
			iter.innerIter = nil
		}

		next := iter.currentIter.Next()
		if next.IsNone() {
			return None[T]()
		}

		inner, ok := any(next.Some()).(Slice[T])
		if !ok {
			return next
		}

		iter.innerIter = flatten(inner.Iter())
	}
}
