package main

import (
	"fmt"
	"math"
	"strings"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
	"github.com/enetx/g/f"
)

func main() {
	ws := SliceOf[String]("d", "b", "c", "a")

	// ws.SortBy(String.Cmp)
	// or
	ws.SortBy(cmp.Cmp)

	ws.Println()

	wsl := SliceOf[String](
		"aa a",
		"b b",
		"o o oo oooo",
		"aa a a aa a",
		"aaa aaaa",
		"a a a aaaa",
		"aaa aaa",
		"three three",
		"one",
		"four",
	)

	// wsl.SortBy(func(a, b String) cmp.Ordering { return a.Cmp(b) }) // or
	// wsl.SortBy(String.Cmp) // or
	wsl.SortBy(cmp.Cmp)
	wsl.Println() // Slice[a a a aaaa, aa a, aa a a aa a, aaa aaa, aaa aaaa, b b, four, o o oo oooo, one, three three]

	wsl.SortBy(func(a, b String) cmp.Ordering { return a.Cmp(b).Reverse() })
	wsl.Println() // Slice[three three, one, o o oo oooo, four, b b, aaa aaaa, aaa aaa, aa a a aa a, aa a, a a a aaaa]

	wsl.SortBy(func(a, b String) cmp.Ordering {
		return a.Fields().Collect().Len().Cmp(b.Fields().Collect().Len()).
			Then(a.Len().Cmp(b.Len()))
	})

	wsl.Println() // Slice[one, four, b b, aa a, aaa aaa, aaa aaaa, three three, a a a aaaa, o o oo oooo, aa a a aa a]

	slice := Slice[int]{1, 2, 3, 4}

	slice.Iter().Filter(f.Gt(2)).Collect().Println() // Slice[3, 4]
	slice.Iter().Filter(f.Eq(2)).Collect().Println() // Slice[2]
	slice.Iter().Filter(f.Ne(2)).Collect().Println() // Slice[1, 3, 4]

	fmt.Println(slice.Iter().All(func(i int) bool { return i != 5 }))
	fmt.Println(slice.Iter().Any(func(i int) bool { return i == 5 }))

	slice.Iter().Range(func(val int) bool {
		fmt.Println(val)
		return val != 3
	})

	result := Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result.Delete(1)
	result.Println() // Slice[1, 3, 4, 5, 6, 7, 8, 9, 10]
	result.Delete(-9)
	result.Println() // Slice[1, 3, 4, 5, 6, 7, 8, 9, 10]

	result.SubSlice(1, -3).Println()            // Slice[2, 3, 4, 5, 6, 7]
	result.SubSlice(-3, result.Len()).Println() // Slice[8, 9, 10]
	result.SubSlice(-3, -1).Println()           // Slice[8, 9]
	result.SubSlice(-1, result.Len()).Println() // Slice[10]

	result = result.RandomSample(5)

	result = result.Clone()
	result.Push(999)
	result.Println()
	fmt.Printf("%#v\n", result.Std())

	filled := NewSlice[int](10)
	filled.Fill(88)

	filled.Println()

	slice = Slice[int]{1, 2, 3, 4, 5}.Println()

	sc := slice.Clone()
	sc.Delete(1, 3)
	sc.Println() // Slice[1, 4, 5]

	sc = slice.Clone()
	sc.Delete(-4, 3)
	sc.Println() // Slice[1, 4, 5]

	sc = slice.Clone()
	sc.Delete(-3, 4)
	sc.Println() // Slice[1, 2, 5]

	// InPlace Methods
	sipl := NewSlice[int]()

	sipl.Push(1)
	sipl.Push(2)
	sipl.Push(3)

	sipl.Delete(1)
	sipl.Fill(999999)

	sipl.Insert(0, 22, 33, 44)
	sipl.PushUnique(22, 22, 22, 33, 44, 55)

	sipl.Println()

	slicea := Slice[string]{"a", "b", "c", "d"}
	slicea.Insert(2, "e", "f")
	slicea.Println()

	slice = Slice[int]{1, 2, 0, 4, 0, 3, 0, 0, 0, 0}
	slice = slice.Iter().Exclude(f.IsZero).Collect()

	slice.Delete(0)
	slice.Println()

	sll := NewSlice[int](0, 100000)
	sll.Push(1)
	sll.Clip()

	fmt.Println(sll.Cap())

	TransformSlice([]string{"AAA", "BBB"}, NewString).Iter().Map(String.Lower).Collect().Println()
	SliceOf([]string{"AAA", "BBB"}...).Iter().Map(strings.ToLower).Collect().Println()

	SliceOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11).Iter().
		Filter(isPrime).
		ForEach(func(n int) { fmt.Printf("%d is a prime number\n", n) })
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}

	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}
