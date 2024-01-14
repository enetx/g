package g

import (
	"bufio"
	"context"
	"io"
)

// Collect reads all values from the iterator until an error occurs and returns them as a slice.
func (iter *baseIterF) Collect() Slice[String] {
	values := make([]String, 0)

	for {
		next := iter.Next()
		if next.IsErr() {
			return values
		}

		values = append(values, next.Ok())
	}
}

// Skip returns a new iterator skipping the first n elements.
func (iter *baseIterF) Skip(n uint) *skipIterF {
	return skipF(iter, n)
}

// Inspect applies a function to each element in the iterator without modifying the elements.
func (iter *baseIterF) Inspect(fn func(String)) *inspectIterF {
	return inspectF(iter, fn)
}

// StepBy creates an iterator that steps over elements in the base iterator with a specified step.
func (iter *baseIterF) StepBy(n int) *stepByIterF {
	return stepByF(iter, n)
}

// ForEach applies a function to each element in the iterator.
func (iter *baseIterF) ForEach(fn func(String)) {
	for {
		next := iter.Next()
		if next.IsErr() {
			return
		}

		fn(next.Ok())
	}
}

// Range applies a function to each element in the iterator until the function returns false.
func (iter *baseIterF) Range(fn func(String) bool) {
	for {
		next := iter.Next()
		if next.IsErr() || !fn(next.Ok()) {
			return
		}
	}
}

// Find returns the first element that satisfies a given condition.
func (iter *baseIterF) Find(fn func(String) bool) Result[String] {
	for {
		next := iter.Next()
		if next.IsErr() {
			return Err[String](io.EOF)
		}

		if fn(next.Ok()) {
			return next
		}
	}
}

// Unique creates an iterator that filters out duplicate elements.
func (iter *baseIterF) Unique() *uniqueIterF {
	return uniqueF(iter)
}

// Dedup creates an iterator that filters out consecutive duplicate elements.
func (iter *baseIterF) Dedup() *dedupIterF {
	return dedupF(iter)
}

// Map creates an iterator that applies a function to each element in the base iterator.
func (iter *baseIterF) Map(fn func(String) String) *mapIterF {
	return mapiterF(iter, fn)
}

// Exclude creates an iterator that excludes elements satisfying a given condition.
func (iter *baseIterF) Exclude(fn func(String) bool) *filterIterF {
	return excludeF(iter, fn)
}

// ToChannel converts the iterator to a channel for concurrent consumption.
func (iter *baseIterF) ToChannel(ctxs ...context.Context) chan String {
	ch := make(chan String)

	ctx := context.Background()
	if len(ctxs) != 0 {
		ctx = ctxs[0]
	}

	go func() {
		defer close(ch)

		for {
			next := iter.Next()
			if next.IsErr() {
				return
			}

			select {
			case <-ctx.Done():
				return
			default:
				ch <- next.Ok()
			}
		}
	}()

	return ch
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// lift
type liftIterF struct {
	baseIterF
	r         *bufio.Reader
	f         *File
	exhausted bool
}

func liftF(f *File, r io.Reader) *liftIterF {
	iter := &liftIterF{f: f, r: bufio.NewReader(r)}
	iter.baseIterF = baseIterF{iter}
	return iter
}

func (iter *liftIterF) Next() Result[String] {
	if iter.exhausted {
		return Err[String](io.EOF)
	}

	content, err := iter.r.ReadString('\n')
	if err == io.EOF {
		iter.exhausted = true
		iter.f.Close()

		return Ok(String(content))
	}

	if err != nil {
		iter.exhausted = true
		iter.f.Close()

		return Err[String](err)
	}

	return Ok(String(content).TrimRight("\r\n"))
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// map
type mapIterF struct {
	baseIterF
	iter      iteratorF
	fn        func(String) String
	exhausted bool
}

func mapiterF(iter iteratorF, fn func(String) String) *mapIterF {
	iterator := &mapIterF{iter: iter, fn: fn}
	iterator.baseIterF = baseIterF{iterator}

	return iterator
}

func (iter *mapIterF) Next() Result[String] {
	if iter.exhausted {
		return Err[String](io.EOF)
	}

	next := iter.iter.Next()

	if next.IsErr() {
		iter.exhausted = true
		return next
	}

	return Ok(iter.fn(next.Ok()))
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// filter
type filterIterF struct {
	baseIterF
	iter      iteratorF
	fn        func(String) bool
	exhausted bool
}

func filterF(iter iteratorF, fn func(String) bool) *filterIterF {
	iterator := &filterIterF{iter: iter, fn: fn}
	iterator.baseIterF = baseIterF{iterator}

	return iterator
}

func (iter *filterIterF) Next() Result[String] {
	if iter.exhausted {
		return Err[String](io.EOF)
	}

	for {
		next := iter.iter.Next()
		if next.IsErr() {
			iter.exhausted = true
			return next
		}

		if iter.fn(next.Ok()) {
			return next
		}
	}
}

func excludeF(iter iteratorF, fn func(String) bool) *filterIterF {
	inverse := func(s String) bool { return !fn(s) }
	iterator := &filterIterF{iter: iter, fn: inverse}
	iterator.baseIterF = baseIterF{iterator}

	return iterator
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// unique
type uniqueIterF struct {
	baseIterF
	iter      iteratorF
	seen      map[String]struct{}
	exhausted bool
}

func uniqueF(iter iteratorF) *uniqueIterF {
	iterator := &uniqueIterF{iter: iter}
	iterator.baseIterF = baseIterF{iterator}
	iterator.seen = make(map[String]struct{})

	return iterator
}

func (iter *uniqueIterF) Next() Result[String] {
	if iter.exhausted {
		return Err[String](io.EOF)
	}

	for {
		next := iter.iter.Next()
		if next.IsErr() {
			iter.exhausted = true
			return next
		}

		val := next.Ok()
		if _, ok := iter.seen[val]; !ok {
			iter.seen[val] = struct{}{}
			return next
		}
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// dedup
type dedupIterF struct {
	baseIterF
	iter      iteratorF
	current   String
	exhausted bool
}

func dedupF(iter iteratorF) *dedupIterF {
	iterator := &dedupIterF{iter: iter}
	iterator.baseIterF = baseIterF{iterator}

	return iterator
}

func (iter *dedupIterF) Next() Result[String] {
	if iter.exhausted {
		return Err[String](io.EOF)
	}

	for {
		next := iter.iter.Next()
		if next.IsErr() {
			iter.exhausted = true
			return next
		}

		if iter.current.Ne(next.Ok()) {
			iter.current = next.Ok()
			return next
		}
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// stepby
type stepByIterF struct {
	baseIterF
	iter      iteratorF
	n         int
	counter   uint
	exhausted bool
}

func stepByF(iter iteratorF, n int) *stepByIterF {
	iterator := &stepByIterF{iter: iter, n: n}
	iterator.baseIterF = baseIterF{iterator}

	return iterator
}

func (iter *stepByIterF) Next() Result[String] {
	if iter.exhausted {
		return Err[String](io.EOF)
	}

	for {
		next := iter.iter.Next()
		if next.IsErr() {
			iter.exhausted = true
			return next
		}

		iter.counter++
		if (iter.counter-1)%uint(iter.n) == 0 {
			return next
		}
	}
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// inspect
type inspectIterF struct {
	baseIterF
	iter      iteratorF
	fn        func(String)
	exhausted bool
}

func inspectF(iter iteratorF, fn func(String)) *inspectIterF {
	iterator := &inspectIterF{iter: iter, fn: fn}
	iterator.baseIterF = baseIterF{iterator}

	return iterator
}

func (iter *inspectIterF) Next() Result[String] {
	if iter.exhausted {
		return Err[String](io.EOF)
	}

	next := iter.iter.Next()
	if next.IsErr() {
		iter.exhausted = true
		return next
	}

	iter.fn(next.Ok())

	return next
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// skip
type skipIterF struct {
	baseIterF
	iter      iteratorF
	count     uint
	skipped   bool
	exhausted bool
}

func skipF(iter iteratorF, count uint) *skipIterF {
	iterator := &skipIterF{iter: iter, count: count}
	iterator.baseIterF = baseIterF{iterator}

	return iterator
}

func (iter *skipIterF) Next() Result[String] {
	if iter.exhausted {
		return Err[String](io.EOF)
	}

	if !iter.skipped {
		iter.skipped = true

		for i := uint(0); i < iter.count; i++ {
			if iter.delegateNext().IsErr() {
				return Err[String](io.EOF)
			}
		}
	}

	return iter.delegateNext()
}

func (iter *skipIterF) delegateNext() Result[String] {
	next := iter.iter.Next()
	if next.IsErr() {
		iter.exhausted = true
	}

	return next
}
