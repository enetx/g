package g

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"

	"gitlab.com/x0xO/g/pkg/iter"
	"gitlab.com/x0xO/g/pkg/rand"
)

// NewSlice creates a new Slice of the given generic type T with the specified length and
// capacity.
// The size variadic parameter can have zero, one, or two integer values.
// If no values are provided, an empty Slice with a length and capacity of 0 is returned.
// If one value is provided, it sets both the length and capacity of the Slice.
// If two values are provided, the first value sets the length and the second value sets the
// capacity.
//
// Parameters:
//
// - size ...int: A variadic parameter specifying the length and/or capacity of the Slice
//
// Returns:
//
// - Slice[T]: A new Slice of the specified generic type T with the given length and capacity
//
// Example usage:
//
//	s1 := g.NewSlice[int]()        // Creates an empty Slice of type int
//	s2 := g.NewSlice[int](5)       // Creates an Slice with length and capacity of 5
//	s3 := g.NewSlice[int](3, 10)   // Creates an Slice with length of 3 and capacity of 10
func NewSlice[T any](size ...int) Slice[T] {
	length, capacity := 0, 0

	switch {
	case len(size) > 1:
		length, capacity = size[0], size[1]
	case len(size) == 1:
		length, capacity = size[0], size[0]
	}

	return make(Slice[T], length, capacity)
}

// SliceOf creates a new generic slice containing the provided elements.
func SliceOf[T any](slice ...T) Slice[T] { return slice }

// Compact removes consecutive duplicate elements from a sorted slice efficiently.
// It assumes the slice is already sorted. The function operates faster than the Unique method.
// Type T must support element comparison via reflect.DeepEqual.
// The function takes a pointer to a slice (*Slice[T]) and returns the modified slice.
// If the original slice contains fewer than two elements, the function returns it unchanged.
//
// Example usage:
//
//	slice := g.Slice[int]{2, 2, 3, 4, 4, 4, 5, 5, 6, 7, 7, 8, 8, 8}
//	slice.Compact()
//
//	// slice now contains: [2 3 4 5 6 7 8]
func (sl *Slice[T]) Compact() Slice[T] {
	if sl.Len() < 2 {
		return *sl
	}

	i := 1

	for k := 1; k < sl.Len(); k++ {
		if !reflect.DeepEqual((*sl)[k], (*sl)[k-1]) {
			if i != k {
				(*sl)[i] = (*sl)[k]
			}

			i++
		}
	}

	*sl = (*sl)[:i]

	return *sl
}

// Counter returns an unordered Map with the counts of each unique element in the slice.
// This function is useful when you want to count the occurrences of each unique element in an
// Slice.
//
// Returns:
//
// - Map[any, int]: An unordered Map with keys representing the unique elements in the Slice
// and values representing the counts of those elements.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 1, 2, 1}
//	counts := slice.Counter()
//	// The counts unordered Map will contain:
//	// 1 -> 3 (since 1 appears three times)
//	// 2 -> 2 (since 2 appears two times)
//	// 3 -> 1 (since 3 appears once)
func (sl Slice[T]) Counter() Map[any, int] {
	result := NewMap[any, int](sl.Len())
	sl.ForEach(func(t T) { result[t]++ })

	return result
}

// Enumerate returns a map with the index of each element as the key.
// This function is useful when you want to create an Map where the keys are the indices of the
// elements in an Slice, and the values are the corresponding elements.
//
// Returns:
//
// - Map[int, T]: An Map with keys representing the indices of the elements in the Slice and
// values representing the corresponding elements.
//
// Example usage:
//
//	slice := g.Slice[int]{10, 20, 30}
//	indexedMap := slice.Enumerate()
//	// The indexedMap Map will contain:
//	// 0 -> 10 (since 10 is at index 0)
//	// 1 -> 20 (since 20 is at index 1)
//	// 2 -> 30 (since 30 is at index 2)
func (sl Slice[T]) Enumerate() Map[int, T] {
	result := NewMap[int, T](sl.Len())
	for k, v := range sl {
		result.Set(k, v)
	}

	return result
}

// Fill fills the slice with the specified value.
// This function is useful when you want to create an Slice with all elements having the same
// value.
// This method can be used in place, as it modifies the original slice.
//
// Parameters:
//
// - val T: The value to fill the Slice with.
//
// Returns:
//
// - Slice[T]: A reference to the original Slice filled with the specified value.
//
// Example usage:
//
//	slice := g.Slice[int]{0, 0, 0}
//	slice.Fill(5)
//
// The modified slice will now contain: 5, 5, 5.
func (sl Slice[T]) Fill(val T) Slice[T] {
	for i := range iter.N(sl.Len()) {
		sl.Set(i, val)
	}

	return sl
}

// ToMapHashed returns a map with the hashed version of each element as the key.
func (sl Slice[T]) ToMapHashed() Map[String, T] {
	result := NewMap[String, T](sl.Len())

	sl.ForEach(func(t T) {
		switch val := any(t).(type) {
		case Int:
			result.Set(val.Hash().MD5(), t)
		case int:
			result.Set(Int(val).Hash().MD5(), t)
		case String:
			result.Set(val.Hash().MD5(), t)
		case string:
			result.Set(String(val).Hash().MD5(), t)
		case Bytes:
			result.Set(val.Hash().MD5().ToString(), t)
		case []byte:
			result.Set(Bytes(val).Hash().MD5().ToString(), t)
		case Float:
			result.Set(val.Hash().MD5(), t)
		case float64:
			result.Set(Float(val).Hash().MD5(), t)
		}
	})

	return result
}

// Chunks splits the Slice into smaller chunks of the specified size.
// The function iterates through the Slice, creating new Slice[T] chunks of the specified size.
// If the size is less than or equal to 0 or the Slice is empty,
// it returns an empty slice of Slice[T].
// If the size is greater than or equal to the length of the Slice,
// it returns a slice of Slice[T] containing the original Slice.
//
// Parameters:
//
// - size int: The size of each chunk.
//
// Returns:
//
// - []Slice[T]: A slice of Slice[T] containing the chunks of the original Slice.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5, 6}
//	batches := slice.Chunks(2)
//
// The resulting chunks will be: [{1, 2}, {3, 4}, {5, 6}].
func (sl Slice[T]) Chunks(size int) []Slice[T] {
	if size <= 0 || sl.Empty() {
		return nil
	}

	chunks := (sl.Len() + size - 1) / size // Round up to ensure all items are included
	result := make([]Slice[T], 0, chunks)

	for i := 0; i < sl.Len(); i += size {
		end := i + size
		if end > sl.Len() {
			end = sl.Len()
		}

		result = append(result, sl.extract(i, end))
	}

	return result
}

// All returns true if all elements in the slice satisfy the provided condition.
// This function is useful when you want to check if all elements in an Slice meet a certain
// criteria.
//
// Parameters:
//
// - fn func(T) bool: A function that returns a boolean indicating whether the element satisfies
// the condition.
//
// Returns:
//
// - bool: True if all elements in the Slice satisfy the condition, false otherwise.
//
// Example usage:
//
//	slice := g.Slice[int]{2, 4, 6, 8, 10}
//	isEven := func(num int) bool { return num%2 == 0 }
//	allEven := slice.All(isEven)
//
// The resulting allEven will be true since all elements in the slice are even.
func (sl Slice[T]) All(fn func(T) bool) bool {
	for _, val := range sl {
		if !fn(val) {
			return false
		}
	}

	return true
}

// Any returns true if any element in the slice satisfies the provided condition.
// This function is useful when you want to check if at least one element in an Slice meets a
// certain criteria.
//
// Parameters:
//
// - fn func(T) bool: A function that returns a boolean indicating whether the element satisfies
// the condition.
//
// Returns:
//
// - bool: True if at least one element in the Slice satisfies the condition, false otherwise.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 3, 5, 7, 9}
//	isEven := func(num int) bool { return num%2 == 0 }
//	anyEven := slice.Any(isEven)
//
// The resulting anyEven will be false since none of the elements in the slice are even.
func (sl Slice[T]) Any(fn func(T) bool) bool {
	for _, val := range sl {
		if fn(val) {
			return true
		}
	}

	return false
}

// Index returns the index of the first occurrence of the specified value in the slice, or -1 if
// not found.
func (sl Slice[T]) Index(val T) int {
	for i, v := range sl {
		if reflect.DeepEqual(v, val) {
			return i
		}
	}

	return -1
}

// RandomSample returns a new slice containing a random sample of elements from the original slice.
// The sampling is done without replacement, meaning that each element can only appear once in the result.
//
// Parameters:
//
// - sequence int: The number of unique elements to include in the random sample.
//
// Returns:
//
// - Slice[T]: A new Slice containing the random sample of unique elements.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9}
//	sample := slice.RandomSample(3)
//
// The resulting sample will contain 3 unique elements randomly selected from the original slice.
func (sl Slice[T]) RandomSample(sequence int) Slice[T] {
	if sequence >= sl.Len() {
		return sl
	}

	return sl.Clone().Shuffle()[0:sequence]
}

// RandomRange returns a new slice containing a random sample of elements from a subrange of the original slice.
// The sampling is done without replacement, meaning that each element can only appear once in the result.
func (sl Slice[T]) RandomRange(from, to int) Slice[T] {
	if from < 0 {
		from = 0
	}

	if to < 0 || to > sl.Len() {
		to = sl.Len()
	}

	if from > to {
		from = to
	}

	return sl.RandomSample(Int(from).RandomRange(Int(to)).Std())
}

// Insert inserts values at the specified index in the slice and returns the resulting slice.
// The original slice remains unchanged.
//
// Parameters:
//
// - i int: The index at which to insert the new values.
//
// - values ...T: A variadic list of values to insert at the specified index.
//
// Returns:
//
// - Slice[T]: A new Slice containing the original elements and the inserted values.
//
// Example usage:
//
//	slice := g.Slice[string]{"a", "b", "c", "d"}
//	newSlice := slice.Insert(2, "e", "f")
//
// The resulting newSlice will be: ["a", "b", "e", "f", "c", "d"].
func (sl Slice[T]) Insert(i int, values ...T) Slice[T] { return sl.Replace(i, i, values...) }

// InsertInPlace inserts values at the specified index in the slice and modifies the original
// slice.
//
// Parameters:
//
// - i int: The index at which to insert the new values.
//
// - values ...T: A variadic list of values to insert at the specified index.
//
// Example usage:
//
//	slice := g.Slice[string]{"a", "b", "c", "d"}
//	slice.InsertInPlace(2, "e", "f")
//
// The resulting slice will be: ["a", "b", "e", "f", "c", "d"].
func (sl *Slice[T]) InsertInPlace(i int, values ...T) Slice[T] {
	sl.ReplaceInPlace(i, i, values...)
	return *sl
}

// Replace replaces the elements of sl[i:j] with the given values, and returns
// a new slice with the modifications. The original slice remains unchanged.
// Replace panics if sl[i:j] is not a valid slice of sl.
//
// Parameters:
//
// - i int: The starting index of the slice to be replaced.
//
// - j int: The ending index of the slice to be replaced.
//
// - values ...T: A variadic list of values to replace the existing slice.
//
// Returns:
//
// - Slice[T]: A new Slice containing the original elements with the specified elements replaced.
//
// Example usage:
//
//	slice := g.Slice[string]{"a", "b", "c", "d"}
//	newSlice := slice.Replace(1, 3, "e", "f")
//
// The original slice remains ["a", "b", "c", "d"], and the newSlice will be: ["a", "e", "f", "d"].
func (sl Slice[T]) Replace(i, j int, values ...T) Slice[T] {
	i = sl.normalizeIndex(i)
	j = sl.normalizeIndex(j)

	if i > j {
		return NewSlice[T]()
	}

	total := sl[:i].Len() + len(values) + sl[j:].Len()
	slice := NewSlice[T](total)

	copy(slice, sl[:i])
	copy(slice[i:], values)
	copy(slice[i+len(values):], sl[j:])

	return slice
}

// ReplaceInPlace replaces the elements of sl[i:j] with the given values,
// and modifies the original slice in place. ReplaceInPlace panics if sl[i:j]
// is not a valid slice of sl.
//
// Parameters:
//
// - i int: The starting index of the slice to be replaced.
//
// - j int: The ending index of the slice to be replaced.
//
// - values ...T: A variadic list of values to replace the existing slice.
//
// Example usage:
//
//	slice := g.Slice[string]{"a", "b", "c", "d"}
//	slice.ReplaceInPlace(1, 3, "e", "f")
//
// After the ReplaceInPlace operation, the resulting slice will be: ["a", "e", "f", "d"].
func (sl *Slice[T]) ReplaceInPlace(i, j int, values ...T) Slice[T] {
	i = sl.normalizeIndex(i)
	j = sl.normalizeIndex(j)

	if i > j {
		*sl = (*sl)[:0]
		return *sl
	}

	if i == j {
		if len(values) > 0 {
			*sl = (*sl)[:i].Append(append(values, (*sl)[i:]...)...)
		}

		return *sl
	}

	diff := len(values) - (j - i)

	if diff > 0 {
		*sl = (*sl).Append(NewSlice[T](diff)...)
	}

	copy((*sl)[i+len(values):], (*sl)[j:])
	copy((*sl)[i:], values)

	if diff < 0 {
		*sl = (*sl)[:sl.Len()+diff]
	}

	return *sl
}

// Unique returns a new slice containing unique elements from the current slice.
// The order of elements in the returned slice is not guaranteed to be the same as in the original slice.
//
// Returns:
// - Slice[T]: A new Slice containing unique elements from the current slice.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 2, 4, 5, 3}
//	unique := slice.Unique()
//
// The resulting unique slice will be: [1, 2, 3, 4, 5].
func (sl Slice[T]) Unique() Slice[T] {
	seen := NewMap[any, struct{}](sl.Len())
	sl.ForEach(func(t T) { seen.Set(t, struct{}{}) })

	unique := NewSlice[T](0, seen.Len())
	seen.ForEach(func(k any, _ struct{}) { unique = unique.Append(k.(T)) })

	return unique
}

// ForEach applies a given function to each element in the slice.
//
// The function takes one parameter of type T (the same type as the elements of the slice).
// The function is applied to each element in the order they appear in the slice.
//
// Parameters:
//
// - fn (func(T)): The function to be applied to each element of the slice.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3}
//	slice.ForEach(func(val int) {
//	    fmt.Println(val * 2)
//	})
//	// Output:
//	// 2
//	// 4
//	// 6
func (sl Slice[T]) ForEach(fn func(T)) {
	for _, val := range sl {
		fn(val)
	}
}

// ForEachBack applies a given function to each element in the slice in reverse order.
//
// The function takes one parameter of type T (the same type as the elements of the slice).
// The function is applied to each element in the reverse order they appear in the slice, starting from the last element.
//
// Parameters:
//
// - fn (func(T)): The function to be applied to each element of the slice in reverse order.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3}
//	slice.ForEachBack(func(val int) {
//	    fmt.Println(val * 2)
//	})
//	// Output:
//	// 6
//	// 4
//	// 2
func (sl Slice[T]) ForEachBack(fn func(T)) {
	for i := sl.LastIndex(); i >= 0; i-- {
		fn(sl[i])
	}
}

// ForEachParallel applies a given function to each element in the slice concurrently.
//
// If the length of the slice is below a certain threshold (max), it performs the operation sequentially.
// Otherwise, it divides the slice into halves and processes each half concurrently using goroutines.
//
// Parameters:
// - fn (func(T)): The function to be applied to each element of the slice.
//
// Note:
// The provided function 'fn' should be safe for concurrent execution to prevent race conditions.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, ...}
//	slice.ForEachParallel(func(val int) {
//	    fmt.Println(val * 2)
//	})
//	// Output (order may vary due to concurrent execution):
//	// ...
func (sl Slice[T]) ForEachParallel(fn func(T)) {
	const max = 1 << 11
	if sl.Len() < max {
		sl.ForEach(fn)
		return
	}

	half := sl.Len() / 2
	left := sl.extract(0, half)
	right := sl.extract(half, sl.Len())

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		left.ForEachParallel(fn)

		wg.Done()
	}()

	right.ForEachParallel(fn)

	wg.Wait()
}

// Range applies a given function to each element in the slice until the function returns false.
//
// The function takes one parameter of type T (the same type as the elements of the slice).
// The function is applied to each element in the order they appear in the slice until the provided function
// returns false. Once the function returns false for an element, the iteration stops.
//
// Parameters:
//
// - fn (func(T) bool): The function to be applied to each element of the slice.
// It should return a boolean value. If it returns false, the iteration will stop.
//
// Example usage:
//
//   slice := g.Slice[int]{1, 2, 3, 4, 5}
//   slice.Range(func(val int) bool {
//       fmt.Println(val)
//       return val != 3
//   })
//   // Output:
//   // 1
//   // 2
//   // 3

func (sl Slice[T]) Range(fn func(T) bool) {
	for _, val := range sl {
		if !fn(val) {
			break
		}
	}
}

// RangeBack applies a given function to each element in the slice in reverse order until the function returns false.
//
// The function takes one parameter of type T (the same type as the elements of the slice).
// The function is applied to each element in reverse order they appear in the slice until the provided function
// returns false. Once the function returns false for an element, the iteration stops.
//
// Parameters:
//
// - fn (func(T) bool): The function to be applied to each element of the slice.
// It should return a boolean value. If it returns false, the iteration will stop.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5}
//	slice.RangeBack(func(val int) bool {
//	    fmt.Println(val)
//	    return val != 3
//	})
//	// Output:
//	// 5
//	// 4
//	// 3
func (sl Slice[T]) RangeBack(fn func(T) bool) {
	for i := sl.LastIndex(); i >= 0; i-- {
		if !fn(sl[i]) {
			break
		}
	}
}

// Map returns a new slice by applying a given function to each element in the current slice.
//
// The function takes one parameter of type T (the same type as the elements of the slice)
// and returns a value of type T. The returned value is added to a new slice,
// which is then returned as the result.
//
// Parameters:
//
// - fn (func(T) T): The function to be applied to each element of the slice.
//
// Returns:
//
// - Slice[T]: A new slice containing the results of applying the function to each element
// of the current slice.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3}
//	doubled := slice.Map(func(val int) int {
//	    return val * 2
//	})
//	fmt.Println(doubled)
//
// Output: [2 4 6].
func (sl Slice[T]) Map(fn func(T) T) Slice[T] { return SliceMap(sl, fn) }

// MapInPlace applies a given function to each element in the current slice,
// modifying the elements in place.
//
// The function takes one parameter of type T (the same type as the elements of the slice)
// and returns a value of type T. The returned value replaces the original element in the slice.
//
// Parameters:
//
// - fn (func(T) T): The function to be applied to each element of the slice.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3}
//	slice.MapInPlace(func(val int) int {
//	    return val * 2
//	})
//	fmt.Println(slice)
//
// Output: [2 4 6].
func (sl *Slice[T]) MapInPlace(fn func(T) T) Slice[T] {
	for i := range iter.N(sl.Len()) {
		sl.Set(i, fn(sl.Get(i)))
	}

	return *sl
}

// Filter returns a new slice containing elements that satisfy a given condition.
//
// The function takes one parameter of type T (the same type as the elements of the slice)
// and returns a boolean value. If the returned value is true, the element is added
// to a new slice, which is then returned as the result.
//
// Parameters:
//
// - fn (func(T) bool): The function to be applied to each element of the slice
// to determine if it should be included in the result.
//
// Returns:
//
// - Slice[T]: A new slice containing the elements that satisfy the given condition.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5}
//	even := slice.Filter(func(val int) bool {
//	    return val%2 == 0
//	})
//	fmt.Println(even)
//
// Output: [2 4].
func (sl Slice[T]) Filter(fn func(T) bool) Slice[T] {
	result := NewSlice[T](0, sl.Len())

	sl.ForEach(func(t T) {
		if fn(t) {
			result = result.Append(t)
		}
	})

	return result.Clip()
}

// FilterInPlace removes elements from the current slice that do not satisfy a given condition.
//
// The function takes one parameter of type T (the same type as the elements of the slice)
// and returns a boolean value. If the returned value is false, the element is removed
// from the slice.
//
// Parameters:
//
// - fn (func(T) bool): The function to be applied to each element of the slice
// to determine if it should be kept in the slice.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5}
//	slice.FilterInPlace(func(val int) bool {
//	    return val%2 == 0
//	})
//	fmt.Println(slice)
//
// Output: [2 4].
func (sl *Slice[T]) FilterInPlace(fn func(T) bool) Slice[T] {
	j := 0

	for i := range iter.N(sl.Len()) {
		if fn(sl.Get(i)) {
			sl.Set(j, sl.Get(i))
			j++
		}
	}

	*sl = (*sl)[:j]

	return *sl
}

// Reduce reduces the slice to a single value using a given function and an initial value.
//
// The function takes two parameters of type T (the same type as the elements of the slice):
// an accumulator and a value from the slice. The accumulator is initialized with the provided
// initial value, and the function is called for each element in the slice. The returned value
// from the function becomes the new accumulator value for the next iteration. After processing
// all the elements in the slice, the final accumulator value is returned as the result.
//
// Parameters:
//
// - fn (func(acc, val T) T): The function to be applied to each element of the slice
// and the accumulator. This function should return a new value for the accumulator.
//
// - initial (T): The initial value for the accumulator.
//
// Returns:
//
// - T: The final accumulator value after processing all the elements in the slice.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5}
//	sum := slice.Reduce(func(acc, val int) int {
//	    return acc + val
//	}, 0)
//	fmt.Println(sum)
//
// Output: 15.
func (sl Slice[T]) Reduce(fn func(acc, val T) T, initial T) T {
	acc := initial

	sl.ForEach(func(t T) { acc = fn(acc, t) })

	return acc
}

// MapParallel applies a given function to each element in the slice in parallel and returns a new
// slice.
//
// The function iterates over the elements of the slice and applies the provided function
// to each element. If the length of the slice is less than a predefined threshold (max),
// it falls back to the sequential Map function. Otherwise, the slice is divided into two
// halves and the function is applied to each half in parallel using goroutines. The
// resulting slices are then combined to form the final output slice.
//
// Note: The order of the elements in the output slice may not be the same as the input
// slice due to parallel processing.
//
// Parameters:
//
// - fn (func(T) T): The function to be applied to each element of the slice.
//
// Returns:
//
// - Slice[T]: A new slice with the results of applying the given function to each element
// of the original slice.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5}
//	squared := slice.MapParallel(func(val int) int {
//	    return val * val
//	})
//	fmt.Println(squared)
//
// Output: {1 4 9 16 25}.
func (sl Slice[T]) MapParallel(fn func(T) T) Slice[T] {
	const max = 1 << 11
	if sl.Len() < max {
		return sl.Map(fn)
	}

	half := sl.Len() / 2
	left := sl.extract(0, half)
	right := sl.extract(half, sl.Len())

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		left = left.MapParallel(fn)

		wg.Done()
	}()

	right = right.MapParallel(fn)

	wg.Wait()

	return NewSlice[T](0, sl.Len()).Append(left...).Append(right...)
}

// FilterParallel returns a new slice containing elements that satisfy a given condition, computed
// in parallel.
//
// The function iterates over the elements of the slice and applies the provided predicate
// function to each element. If the length of the slice is less than a predefined threshold (max),
// it falls back to the sequential Filter function. Otherwise, the slice is divided into two
// halves and the predicate function is applied to each half in parallel using goroutines. The
// resulting slices are then combined to form the final output slice.
//
// Note: The order of the elements in the output slice may not be the same as the input
// slice due to parallel processing.
//
// Parameters:
//
// - fn (func(T) bool): The predicate function to be applied to each element of the slice.
//
// Returns:
//
// - Slice[T]: A new slice containing the elements that satisfy the given condition.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5}
//	even := slice.FilterParallel(func(val int) bool {
//	    return val % 2 == 0
//	})
//	fmt.Println(even)
//
// Output: {2 4}.
func (sl Slice[T]) FilterParallel(fn func(T) bool) Slice[T] {
	const max = 1 << 11
	if sl.Len() < max {
		return sl.Filter(fn)
	}

	half := sl.Len() / 2
	left := sl.extract(0, half)
	right := sl.extract(half, sl.Len())

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		left = left.FilterParallel(fn)

		wg.Done()
	}()

	right = right.FilterParallel(fn)

	wg.Wait()

	return NewSlice[T](0, left.Len()+right.Len()).Append(left...).Append(right...)
}

// ReduceParallel reduces the slice to a single value using a given function and an initial value,
// computed in parallel.
//
// The function iterates over the elements of the slice and applies the provided reducer function
// to each element in a pairwise manner. If the length of the slice is less than a predefined
// threshold (max),
// it falls back to the sequential Reduce function. Otherwise, the slice is divided into two
// halves and the reducer function is applied to each half in parallel using goroutines. The
// resulting values are combined using the reducer function to produce the final output value.
//
// Note: Due to parallel processing, the order in which the reducer function is applied to the
// elements may not be the same as the input slice.
//
// Parameters:
//
// - fn (func(T, T) T): The reducer function to be applied to each element of the slice.
//
// - initial (T): The initial value to be used as the starting point for the reduction.
//
// Returns:
//
// - T: A single value obtained by applying the reducer function to the elements of the slice.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5}
//	sum := slice.ReduceParallel(func(acc, val int) int {
//	    return acc + val
//	}, 0)
//	fmt.Println(sum)
//
// Output: 15.
func (sl Slice[T]) ReduceParallel(fn func(T, T) T, initial T) T {
	const max = 1 << 11
	if sl.Len() < max {
		return sl.Reduce(fn, initial)
	}

	half := sl.Len() / 2
	left := sl.extract(0, half)
	right := sl.extract(half, sl.Len())

	result := NewSlice[T](2)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		result.Set(0, left.ReduceParallel(fn, initial))

		wg.Done()
	}()

	result.Set(1, right.ReduceParallel(fn, initial))

	wg.Wait()

	return result.Reduce(fn, initial)
}

// AddUnique appends unique elements from the provided arguments to the current slice.
//
// The function iterates over the provided elements and checks if they are already present
// in the slice. If an element is not already present, it is appended to the slice. The
// resulting slice is returned, containing the unique elements from both the original
// slice and the provided elements.
//
// Parameters:
//
// - elems (...T): A variadic list of elements to be appended to the slice.
//
// Returns:
//
// - Slice[T]: A new slice containing the unique elements from both the original slice
// and the provided elements.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5}
//	slice = slice.AddUnique(3, 4, 5, 6, 7)
//	fmt.Println(slice)
//
// Output: [1 2 3 4 5 6 7].
func (sl Slice[T]) AddUnique(elems ...T) Slice[T] {
	for _, elem := range elems {
		if !sl.Contains(elem) {
			sl = sl.Append(elem)
		}
	}

	return sl
}

// AddUniqueInPlace appends unique elements from the provided arguments to the current slice.
//
// The function iterates over the provided elements and checks if they are already present
// in the slice. If an element is not already present, it is appended to the slice.
//
// Parameters:
//
// - elems (...T): A variadic list of elements to be appended to the slice.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5}
//	slice.AddUniqueInPlace(3, 4, 5, 6, 7)
//	fmt.Println(slice)
//
// Output: [1 2 3 4 5 6 7].
func (sl *Slice[T]) AddUniqueInPlace(elems ...T) Slice[T] {
	for _, elem := range elems {
		if !sl.Contains(elem) {
			*sl = sl.Append(elem)
		}
	}

	return *sl
}

// Get returns the element at the given index, handling negative indices as counting from the end
// of the slice.
func (sl Slice[T]) Get(index int) T {
	index = sl.normalizeIndex(index)
	return sl[index]
}

// Count returns the count of the given element in the slice.
func (sl Slice[T]) Count(elem T) int {
	if sl.Empty() {
		return 0
	}

	var counter int

	sl.ForEach(func(t T) {
		if reflect.DeepEqual(t, elem) {
			counter++
		}
	})

	return counter
}

// Max returns the maximum element in the slice, assuming elements are comparable.
func (sl Slice[T]) Max() T {
	if sl.Empty() {
		return *new(T)
	}

	max := sl.Get(0)

	var greater func(a, b any) bool

	switch any(max).(type) {
	case Int:
		greater = func(a, b any) bool { return a.(Int).Gt(b.(Int)) }
	case int:
		greater = func(a, b any) bool { return a.(int) > b.(int) }
	case String:
		greater = func(a, b any) bool { return a.(String).Gt(b.(String)) }
	case string:
		greater = func(a, b any) bool { return a.(string) > b.(string) }
	case Float:
		greater = func(a, b any) bool { return a.(Float).Gt(b.(Float)) }
	case float64:
		greater = func(a, b any) bool { return Float(a.(float64)).Gt(Float(b.(float64))) }
	}

	sl.ForEach(func(t T) {
		if greater(t, max) {
			max = t
		}
	})

	return max
}

// Min returns the minimum element in the slice, assuming elements are comparable.
func (sl Slice[T]) Min() T {
	if sl.Empty() {
		return *new(T)
	}

	min := sl.Get(0)

	var less func(a, b any) bool

	switch any(min).(type) {
	case Int:
		less = func(a, b any) bool { return a.(Int).Lt(b.(Int)) }
	case int:
		less = func(a, b any) bool { return a.(int) < b.(int) }
	case String:
		less = func(a, b any) bool { return a.(String).Lt(b.(String)) }
	case string:
		less = func(a, b any) bool { return a.(string) < b.(string) }
	case Float:
		less = func(a, b any) bool { return a.(Float).Lt(b.(Float)) }
	case float64:
		less = func(a, b any) bool { return Float(a.(float64)).Lt(Float(b.(float64))) }
	}

	sl.ForEach(func(t T) {
		if less(t, min) {
			min = t
		}
	})

	return min
}

// Shuffle shuffles the elements in the slice randomly. This method can be used in place, as it
// modifies the original slice.
//
// The function uses the crypto/rand package to generate random indices.
//
// Returns:
//
// - Slice[T]: The modified slice with the elements shuffled randomly.
//
// Example usage:
//
// slice := g.Slice[int]{1, 2, 3, 4, 5}
// shuffled := slice.Shuffle()
// fmt.Println(shuffled)
//
// Output: A randomly shuffled version of the original slice, e.g., [4 1 5 2 3].
func (sl Slice[T]) Shuffle() Slice[T] {
	n := sl.Len()

	for i := n - 1; i > 0; i-- {
		j := rand.N(i + 1)
		sl.Swap(i, j)
	}

	return sl
}

// Reverse reverses the order of the elements in the slice. This method can be used in place, as it
// modifies the original slice.
//
// Returns:
//
// - Slice[T]: The modified slice with the elements reversed.
//
// Example usage:
//
// slice := g.Slice[int]{1, 2, 3, 4, 5}
// reversed := slice.Reverse()
// fmt.Println(reversed)
//
// Output: [5 4 3 2 1].
func (sl Slice[T]) Reverse() Slice[T] {
	for i, j := 0, sl.Len()-1; i < j; i, j = i+1, j-1 {
		sl.Swap(i, j)
	}

	return sl
}

// Sort sorts the elements in the slice in increasing order. It modifies the original
// slice in place. For proper functionality, the type T used in the slice must support
// comparison via the Less method.
func (sl Slice[T]) Sort() Slice[T] {
	sort.Sort(sl)
	return sl
}

// SortBy sorts the elements in the slice using the provided comparison function. This method can
// be used in place, as it modifies the original slice. It requires the elements to be of a type
// that is comparable.
//
// The function takes a custom comparison function as an argument and sorts the elements
// of the slice using the provided logic. The comparison function should return true if
// the element at index i should come before the element at index j, and false otherwise.
//
// Parameters:
//
// - f func(i, j int) bool: A comparison function that takes two indices i and j and returns a bool.
//
// Returns:
//
// - Slice[T]: The sorted Slice.
//
// Example usage:
//
// sl := NewSlice[int](1, 5, 3, 2, 4)
// sl.SortBy(func(i, j int) bool { return sl[i] < sl[j] }) // sorts in ascending order.
func (sl Slice[T]) SortBy(f func(i, j int) bool) Slice[T] {
	sort.Slice(sl, f)
	return sl
}

// FilterZeroValues returns a new slice with all zero values removed.
//
// The function iterates over the elements in the slice and checks if they are
// zero values using the reflect.DeepEqual function. If an element is not a zero value,
// it is added to the resulting slice. The new slice, containing only non-zero values,
// is returned.
//
// Returns:
//
// - Slice[T]: A new slice containing only non-zero values from the original slice.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 0, 4, 0}
//	nonZeroSlice := slice.FilterZeroValues()
//	fmt.Println(nonZeroSlice)
//
// Output: [1 2 4].
func (sl Slice[T]) FilterZeroValues() Slice[T] {
	return sl.Filter(func(v T) bool { return !reflect.DeepEqual(v, *new(T)) })
}

// FilterZeroValuesInPlace removes all zero values from the current slice.
//
// The function iterates over the elements in the slice and checks if they are
// zero values using the reflect.DeepEqual function. If an element is a zero value,
// it is removed from the slice.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 0, 4, 0}
//	slice.FilterZeroValuesInPlace()
//	fmt.Println(slice)
//
// Output: [1 2 4].
func (sl *Slice[T]) FilterZeroValuesInPlace() Slice[T] {
	sl.FilterInPlace(func(v T) bool { return !reflect.DeepEqual(v, *new(T)) })
	return *sl
}

// ToStringSlice converts the slice into a slice of strings.
func (sl Slice[T]) ToStringSlice() []string {
	result := NewSlice[string](0, sl.Len())
	sl.ForEach(func(t T) { result = result.Append(fmt.Sprint(t)) })

	return result
}

// Join joins the elements in the slice into a single String, separated by the provided separator
// (if any).
func (sl Slice[T]) Join(sep ...T) String {
	var separator string
	if len(sep) != 0 {
		separator = fmt.Sprint(sep[0])
	}

	return String(strings.Join(sl.ToStringSlice(), separator))
}

// SubSlice returns a new slice containing elements from the current slice between the specified start
// and end indices, with an optional step parameter to define the increment between elements.
// The function checks if the start and end indices are within the bounds of the original slice.
// If the end index is negative, it represents the position from the end of the slice.
// If the start index is negative, it represents the position from the end of the slice counted
// from the start index.
//
// Parameters:
//
// - start (int): The start index of the range.
//
// - end (int): The end index of the range.
//
// - step (int, optional): The increment between elements. Defaults to 1 if not provided.
// If negative, the slice is traversed in reverse order.
//
// Returns:
//
// - Slice[T]: A new slice containing elements from the current slice between the start and end
// indices, with the specified step.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9}
//	subSlice := slice.SubSlice(1, 7, 2) // Extracts elements 2, 4, 6
//	fmt.Println(subSlice)
//
// Output: [2 4 6].
func (sl Slice[T]) SubSlice(start, end int, step ...int) Slice[T] {
	_step := 1

	if len(step) != 0 {
		_step = step[0]
	}

	start = sl.normalizeIndex(start, struct{}{})
	end = sl.normalizeIndex(end, struct{}{})

	if (start >= end && _step > 0) || (start <= end && _step < 0) || _step == 0 {
		return NewSlice[T]()
	}

	var loopCondition func(int) bool
	if _step > 0 {
		loopCondition = func(i int) bool { return i < end }
	} else {
		loopCondition = func(i int) bool { return i > end }
	}

	var slice Slice[T]

	for i := start; loopCondition(i); i += _step {
		slice = slice.Append(sl[i])
	}

	return slice
}

// Cut removes a range of elements from the Slice and returns a new Slice.
// It creates two slices: one from the beginning of the original slice up to
// the specified start index (exclusive), and another from the specified end
// index (inclusive) to the end of the original slice. These two slices are
// then concatenated to form the resulting Slice.
//
// Parameters:
//
// - start (int): The start index of the range to be removed.
//
// - end (int): The end index of the range to be removed.
//
// Note:
//
//	The function also supports negative indices. Negative indices are counted
//	from the end of the slice. For example, -1 means the last element, -2
//	means the second-to-last element, and so on.
//
// Returns:
//
//	Slice[T]: A new slice containing elements from the current slice with
//	the specified range removed.
//
// Example:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5}
//	newSlice := slice.Cut(1, 3)
//	// newSlice is [1 4 5]
func (sl Slice[T]) Cut(start, end int) Slice[T] { return sl.Replace(start, end) }

// CutInPlace removes a range of elements from the Slice in-place.
// It modifies the original slice by creating two slices: one from the
// beginning of the original slice up to the specified start index
// (exclusive), and another from the specified end index (inclusive)
// to the end of the original slice. These two slices are then
// concatenated to form the modified original Slice.
//
// Parameters:
//
// - start (int): The start index of the range to be removed.
//
// - end (int): The end index of the range to be removed.
//
// Note:
//
// The function also supports negative indices. Negative indices are counted
// from the end of the slice. For example, -1 means the last element, -2
// means the second-to-last element, and so on.
func (sl *Slice[T]) CutInPlace(start, end int) Slice[T] { return sl.ReplaceInPlace(start, end) }

// Random returns a random element from the slice.
//
// The function uses the crypto/rand package to generate a random index within the bounds of the
// slice. If the slice is empty, the zero value of type T is returned.
//
// Returns:
//
// - T: A random element from the slice.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5}
//	randomElement := slice.Random()
//	fmt.Println(randomElement)
//
// Output: <any random element from the slice>.
func (sl Slice[T]) Random() T {
	if sl.Empty() {
		return *new(T)
	}

	return sl.Get(rand.N(sl.Len()))
}

// Clone returns a copy of the slice.
func (sl Slice[T]) Clone() Slice[T] { return append(sl[:0:0], sl...) }

// LastIndex returns the last index of the slice.
func (sl Slice[T]) LastIndex() int {
	if !sl.Empty() {
		return sl.Len() - 1
	}

	return 0
}

// factorial a utility function that calculates the factorial of a given number.
func factorial(n int) int {
	if n <= 1 {
		return 1
	}

	return n * factorial(n-1)
}

// Permutations returns all possible permutations of the elements in the slice.
//
// The function uses a recursive approach to generate all the permutations of the elements.
// If the slice has a length of 0 or 1, it returns the slice itself wrapped in a single-element
// slice.
//
// Returns:
//
// - []Slice[T]: A slice of Slice[T] containing all possible permutations of the elements in the
// slice.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3}
//	perms := slice.Permutations()
//	for _, perm := range perms {
//	    fmt.Println(perm)
//	}
//	// Output:
//	// [1 2 3]
//	// [1 3 2]
//	// [2 1 3]
//	// [2 3 1]
//	// [3 1 2]
//	// [3 2 1]
func (sl Slice[T]) Permutations() []Slice[T] {
	if sl.Len() <= 1 {
		return []Slice[T]{sl}
	}

	perms := make([]Slice[T], 0, factorial(sl.Len()))

	for i, elem := range sl {
		rest := NewSlice[T](sl.Len() - 1)

		copy(rest[:i], sl[:i])
		copy(rest[i:], sl[i+1:])

		subPerms := rest.Permutations()

		for j := range subPerms {
			subPerms[j] = append(Slice[T]{elem}, subPerms[j]...)
		}

		perms = append(perms, subPerms...)
	}

	return perms
}

// Zip zips the elements of the given slices with the current slice into a new slice of Slice[T]
// elements.
//
// The function combines the elements of the current slice with the elements of the given slices by
// index. The length of the resulting slice of Slice[T] elements is determined by the shortest
// input slice.
//
// Params:
//
// - slices: The slices to be zipped with the current slice.
//
// Returns:
//
// - []Slice[T]: A new slice of Slice[T] elements containing the zipped elements of the input
// slices.
//
// Example usage:
//
//	slice1 := g.Slice[int]{1, 2, 3}
//	slice2 := g.Slice[int]{4, 5, 6}
//	slice3 := g.Slice[int]{7, 8, 9}
//	zipped := slice1.Zip(slice2, slice3)
//	for _, group := range zipped {
//	    fmt.Println(group)
//	}
//	// Output:
//	// [1 4 7]
//	// [2 5 8]
//	// [3 6 9]
func (sl Slice[T]) Zip(ss ...Slice[T]) []Slice[T] {
	minLen := sl.Len()

	for _, slice := range ss {
		if slice.Len() < minLen {
			minLen = slice.Len()
		}
	}

	result := make([]Slice[T], 0, minLen)

	for i := range iter.N(minLen) {
		values := NewSlice[T](0, len(ss)+1).Append(sl.Get(i))
		for _, j := range ss {
			values = values.Append(j.Get(i))
		}

		result = append(result, values)
	}

	return result
}

// Flatten flattens the nested slice structure into a single-level Slice[any].
//
// It recursively traverses the nested slice structure and appends all non-slice elements to a new
// Slice[any].
//
// Returns:
//
// - Slice[any]: A new Slice[any] containing the flattened elements.
//
// Example usage:
//
//	nested := g.Slice[any]{1, 2, g.Slice[int]{3, 4, 5}, []any{6, 7, []int{8, 9}}}
//	flattened := nested.Flatten()
//	fmt.Println(flattened)
//
// Output: [1 2 3 4 5 6 7 8 9].
func (sl Slice[T]) Flatten() Slice[any] {
	flattened := NewSlice[any]()
	flattenRecursive(reflect.ValueOf(sl), &flattened)

	return flattened
}

// flattenRecursive a helper function for recursively flattening nested slices.
func flattenRecursive(val reflect.Value, flattened *Slice[any]) {
	for i := range iter.N(val.Len()) {
		elem := val.Index(i)
		if elem.Kind() == reflect.Interface {
			elem = elem.Elem()
		}

		if elem.Kind() == reflect.Slice {
			flattenRecursive(elem, flattened)
		} else {
			*flattened = append(*flattened, elem.Interface())
		}
	}
}

// Eq returns true if the slice is equal to the provided other slice.
func (sl Slice[T]) Eq(other Slice[T]) bool {
	if sl.Len() != other.Len() {
		return false
	}

	for index, val := range sl {
		if !reflect.DeepEqual(val, other.Get(index)) {
			return false
		}
	}

	return true
}

// String returns a string representation of the slice.
func (sl Slice[T]) String() string {
	var builder strings.Builder

	sl.ForEach(func(v T) { builder.WriteString(fmt.Sprintf("%v, ", v)) })

	return String(builder.String()).TrimRight(", ").Format("Slice[%s]").Std()
}

// Append appends the provided elements to the slice and returns the modified slice.
func (sl Slice[T]) Append(elems ...T) Slice[T] { return append(sl, elems...) }

// AppendInPlace appends the provided elements to the slice and modifies the original slice.
func (sl *Slice[T]) AppendInPlace(elems ...T) Slice[T] {
	*sl = sl.Append(elems...)
	return *sl
}

// Cap returns the capacity of the Slice.
func (sl Slice[T]) Cap() int { return cap(sl) }

// Contains returns true if the slice contains the provided value.
func (sl Slice[T]) Contains(val T) bool { return sl.Index(val) >= 0 }

// ContainsAny checks if the Slice contains any element from another Slice.
func (sl Slice[T]) ContainsAny(values ...T) bool {
	if sl.Empty() || len(values) == 0 {
		return false
	}

	seen := NewMap[any, struct{}](sl.Len())
	sl.ForEach(func(t T) { seen.Set(t, struct{}{}) })

	for _, v := range values {
		if seen.Contains(v) {
			return true
		}
	}

	return false
}

// ContainsAll checks if the Slice contains all elements from another Slice.
func (sl Slice[T]) ContainsAll(values ...T) bool {
	if sl.Empty() || len(values) == 0 {
		return false
	}

	seen := NewMap[any, struct{}](sl.Len())
	sl.ForEach(func(t T) { seen.Set(t, struct{}{}) })

	for _, v := range values {
		if !seen.Contains(v) {
			return false
		}
	}

	return true
}

// Delete removes the element at the specified index from the slice and returns the modified slice.
func (sl Slice[T]) Delete(i int) Slice[T] {
	nsl := sl.Clone()
	nsl.DeleteInPlace(i)

	return nsl.Clip()
}

// DeleteInPlace removes the element at the specified index from the slice and modifies the
// original slice.
func (sl *Slice[T]) DeleteInPlace(i int) Slice[T] {
	i = sl.normalizeIndex(i)
	copy((*sl)[i:], (*sl)[i+1:])
	*sl = (*sl)[:sl.Len()-1]

	return *sl
}

// Empty returns true if the slice is empty.
func (sl Slice[T]) Empty() bool { return sl.Len() == 0 }

// Last returns the last element of the slice.
func (sl Slice[T]) Last() T { return sl.Get(-1) }

// Ne returns true if the slice is not equal to the provided other slice.
func (sl Slice[T]) Ne(other Slice[T]) bool { return !sl.Eq(other) }

// NotEmpty returns true if the slice is not empty.
func (sl Slice[T]) NotEmpty() bool { return sl.Len() != 0 }

// Pop returns the last element of the slice and a new slice without the last element.
func (sl Slice[T]) Pop() (T, Slice[T]) { return sl.Last(), sl.SubSlice(0, -1) }

// Set sets the value at the specified index in the slice and returns the modified slice.
// This method can be used in place, as it modifies the original slice.
//
// Parameters:
//
// - i (int): The index at which to set the new value.
//
// - val (T): The new value to be set at the specified index.
//
// Returns:
//
// - Slice[T]: The modified slice with the new value set at the specified index.
//
// Example usage:
//
// slice := g.Slice[int]{1, 2, 3, 4, 5}
// slice.Set(2, 99)
// fmt.Println(slice)
//
// Output: [1 2 99 4 5].
func (sl Slice[T]) Set(index int, val T) Slice[T] {
	index = sl.normalizeIndex(index)
	sl[index] = val

	return sl
}

// Len returns the length of the slice.
func (sl Slice[T]) Len() int { return len(sl) }

// Less defines the comparison logic between two elements at indices i and j within the slice.
// It utilizes type-based comparisons for elements of various types like Int, int, String, string,
// Float, and float64. The comparison is performed according to the types and their respective
// comparison methods. If the types are not directly comparable, it returns false.
func (sl Slice[T]) Less(i, j int) bool {
	elemI := any(sl.Get(i))
	elemJ := any(sl.Get(j))

	switch elemI := elemI.(type) {
	case Int:
		if elemJ, ok := elemJ.(Int); ok {
			return elemI.Lt(elemJ)
		}
	case String:
		if elemJ, ok := elemJ.(String); ok {
			return elemI.Lt(elemJ)
		}
	case Float:
		if elemJ, ok := elemJ.(Float); ok {
			return elemI.Lt(elemJ)
		}
	case int:
		if elemJ, ok := elemJ.(int); ok {
			return elemI < elemJ
		}
	case string:
		if elemJ, ok := elemJ.(string); ok {
			return elemI < elemJ
		}
	case float64:
		if elemJ, ok := elemJ.(float64); ok {
			return Float(elemI).Lt(Float(elemJ))
		}
	}

	return false
}

// Swap swaps the elements at the specified indices in the slice.
// This method used in place, as it modifies the original slice.
//
// Parameters:
//
// - i (int): The index of the first element to be swapped.
//
// - j (int): The index of the second element to be swapped.
//
// Returns:
//
// - Slice[T]: The modified slice with the elements at the specified indices swapped.
//
// Example usage:
//
// slice := g.Slice[int]{1, 2, 3, 4, 5}
// slice.Swap(1, 3)
// fmt.Println(slice)
//
// Output: [1 4 3 2 5].
func (sl Slice[T]) Swap(i, j int) {
	i = sl.normalizeIndex(i)
	j = sl.normalizeIndex(j)

	sl[i], sl[j] = sl[j], sl[i]
}

// Grow increases the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. If n is negative or too large to
// allocate the memory, Grow panics.
func (sl Slice[T]) Grow(n int) Slice[T] {
	if n < 0 {
		panic("cannot be negative")
	}

	if n -= sl.Cap() - sl.Len(); n > 0 {
		sl = append(sl[:sl.Cap()], make(Slice[T], n)...)[:sl.Len()]
	}

	return sl
}

// Clip removes unused capacity from the slice.
func (sl Slice[T]) Clip() Slice[T] { return sl[:sl.Len():sl.Len()] }

// Std returns a new slice with the same elements as the Slice[T].
func (sl Slice[T]) Std() []T { return sl }

// Clear removes all elements from the Slice and sets its length to 0.
func (sl Slice[T]) Clear() Slice[T] { clear(sl); return sl }

// Print prints the elements of the Slice to the standard output (console)
// and returns the Slice unchanged.
func (sl Slice[T]) Print() Slice[T] { fmt.Println(sl); return sl }

func (sl Slice[T]) normalizeIndex(i int, subslice ...struct{}) int {
	ii := i
	if ii < 0 {
		ii += sl.Len()
	}

	negative := 0
	if len(subslice) != 0 {
		negative = -1
	}

	if ii > sl.Len() || ii < negative {
		panic(fmt.Sprintf("runtime error: slice bounds out of range [%d] with length %d", i, sl.Len()))
	}

	return ii
}

func (sl Slice[T]) extract(start, end int) Slice[T] {
	slice := NewSlice[T](end - start)
	copy(slice, sl[start:end])

	return slice
}
