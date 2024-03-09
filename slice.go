package g

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

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

// MapSlice applies the given function to each element of a Slice and returns a new Slice
// containing the transformed values.
//
// Parameters:
//
// - sl: The input Slice.
//
// - fn: The function to apply to each element of the input Slice.
//
// Returns:
//
// A new Slice containing the results of applying the function to each element of the input Slice.
func MapSlice[T, U any](sl Slice[T], fn func(T) U) Slice[U] { return mapSlice(sl.Iter(), fn).Collect() }

// Iter returns an iterator (*liftIter) for the Slice, allowing for sequential iteration
// over its elements. It is commonly used in combination with higher-order functions,
// such as 'ForEach', to perform operations on each element of the Slice.
//
// Returns:
//
// A pointer to a liftIter, which can be used for sequential iteration over the elements of the Slice.
//
// Example usage:
//
//	slice := g.Slice[int]{1, 2, 3, 4, 5}
//	iterator := slice.Iter()
//	iterator.ForEach(func(element int) {
//		// Perform some operation on each element
//		fmt.Println(element)
//	})
//
// The 'Iter' method provides a convenient way to traverse the elements of a Slice
// in a functional style, enabling operations like mapping or filtering.
func (sl Slice[T]) Iter() seqSlice[T] { return liftSlice(sl) }

// func (sl Slice[T]) Iter() *liftIter[T] { return lift(sl) }

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
func (sl Slice[T]) Counter() MapOrd[any, uint] {
	result := NewMapOrd[any, uint](sl.Len())
	sl.Iter().ForEach(func(t T) {
		i := result.GetOrDefault(t, 0)
		i++
		result.Set(t, i)
	})

	return result
}

// Fill fills the slice with the specified value.
// This function is useful when you want to create an Slice with all elements having the same
// value.
// This method modifies the original slice in place.
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
func (sl Slice[T]) Fill(val T) {
	for i := range sl.Len() {
		sl.Set(i, val)
	}
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
// Output: Slice[1, 2, 3, 4, 5, 6, 7, 8, 9].
func (sl Slice[T]) Flatten() Slice[any] {
	flattened := NewSlice[any]()
	flattenRecursive(reflect.ValueOf(sl), &flattened)

	return flattened
}

// flattenRecursive a helper function for recursively flattening nested slices.
func flattenRecursive(val reflect.Value, flattened *Slice[any]) {
	for i := range val.Len() {
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

// ToMapHashed returns a map with the hashed version of each element as the key.
func (sl Slice[T]) ToMapHashed() Map[String, T] {
	result := NewMap[String, T](sl.Len())

	sl.Iter().ForEach(func(t T) {
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

	clone := sl.Clone()
	clone.Shuffle()

	return clone[0:sequence]
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
func (sl *Slice[T]) InsertInPlace(i int, values ...T) { sl.ReplaceInPlace(i, i, values...) }

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
func (sl *Slice[T]) ReplaceInPlace(i, j int, values ...T) {
	i = sl.normalizeIndex(i)
	j = sl.normalizeIndex(j)

	if i > j {
		*sl = (*sl)[:0]
		return
	}

	if i == j {
		if len(values) > 0 {
			*sl = (*sl)[:i].Append(append(values, (*sl)[i:]...)...)
		}

		return
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
func (sl *Slice[T]) AddUniqueInPlace(elems ...T) {
	for _, elem := range elems {
		if !sl.Contains(elem) {
			*sl = sl.Append(elem)
		}
	}
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

	sl.Iter().ForEach(func(t T) {
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

	maxi := sl.Get(0)

	var greater func(a, b any) bool

	switch any(maxi).(type) {
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

	sl.Iter().ForEach(func(t T) {
		if greater(t, maxi) {
			maxi = t
		}
	})

	return maxi
}

// Min returns the minimum element in the slice, assuming elements are comparable.
func (sl Slice[T]) Min() T {
	if sl.Empty() {
		return *new(T)
	}

	mini := sl.Get(0)

	var less func(a, b any) bool

	switch any(mini).(type) {
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

	sl.Iter().ForEach(func(t T) {
		if less(t, mini) {
			mini = t
		}
	})

	return mini
}

// Shuffle shuffles the elements in the slice randomly.
// This method modifies the original slice in place.
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
func (sl Slice[T]) Shuffle() {
	n := sl.Len()

	for i := n - 1; i > 0; i-- {
		j := rand.N(i + 1)
		sl.Swap(i, j)
	}
}

// Reverse reverses the order of the elements in the slice.
// This method modifies the original slice in place.
//
// Returns:
//
// - Slice[T]: The modified slice with the elements reversed.
//
// Example usage:
//
// slice := g.Slice[int]{1, 2, 3, 4, 5}
// slice.Reverse()
// fmt.Println(slice)
//
// Output: [5 4 3 2 1].
func (sl Slice[T]) Reverse() {
	for i, j := 0, sl.Len()-1; i < j; i, j = i+1, j-1 {
		sl.Swap(i, j)
	}
}

// Sort sorts the elements in the slice in increasing order. It modifies the original
// slice in place. For proper functionality, the type T used in the slice must support
// comparison via the Less method.
func (sl Slice[T]) Sort() { sort.Sort(sl) }

// SortBy sorts the elements in the slice using the provided comparison function.
// It modifies the original slice in place. It requires the elements to be of a type
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
// Example usage:
//
// sl := NewSlice[int](1, 5, 3, 2, 4)
// sl.SortBy(func(a, j int) bool { return sl[i] < sl[j] }) // sorts in ascending order.
func (sl Slice[T]) SortBy(fn func(a, b T) bool) Slice[T] {
	sort.Slice(sl, func(i, j int) bool {
		return fn(sl[i], sl[j])
	})

	return sl
}

// ToStringSlice converts the slice into a slice of strings.
func (sl Slice[T]) ToStringSlice() []string {
	result := NewSlice[string](0, sl.Len())
	sl.Iter().ForEach(func(t T) { result = result.Append(fmt.Sprint(t)) })

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
func (sl *Slice[T]) CutInPlace(start, end int) { sl.ReplaceInPlace(start, end) }

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

	sl.Iter().ForEach(func(v T) { builder.WriteString(fmt.Sprintf("%v, ", v)) })

	return String(builder.String()).TrimRight(", ").Format("Slice[%s]").Std()
}

// Append appends the provided elements to the slice and returns the modified slice.
func (sl Slice[T]) Append(elems ...T) Slice[T] { return append(sl, elems...) }

// AppendInPlace appends the provided elements to the slice and modifies the original slice.
func (sl *Slice[T]) AppendInPlace(elems ...T) { *sl = sl.Append(elems...) }

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
	sl.Iter().ForEach(func(t T) { seen.Set(t, struct{}{}) })

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
	sl.Iter().ForEach(func(t T) { seen.Set(t, struct{}{}) })

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
func (sl *Slice[T]) DeleteInPlace(i int) {
	i = sl.normalizeIndex(i)
	copy((*sl)[i:], (*sl)[i+1:])
	*sl = (*sl)[:sl.Len()-1]
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
// This method modifies the original slice in place.
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
func (sl Slice[T]) Set(index int, val T) {
	index = sl.normalizeIndex(index)
	sl[index] = val
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
// This method modifies the original slice in place.
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
