package main

import (
	"fmt"
	"strings"

	"gitlab.com/x0xO/g"
	"gitlab.com/x0xO/g/filters"
)

func main() {
	slice := g.Slice[int]{1, 2, 3, 4}

	fmt.Println(slice.Iter().All(func(i int) bool { return i != 5 }))
	fmt.Println(slice.Iter().Any(func(i int) bool { return i == 5 }))

	slice.Iter().Range(func(val int) bool {
		fmt.Println(val)
		return val != 3
	})

	slice = g.Slice[int]{1, 2, 3, 1, 2, 1}
	slice.Counter().Print()

	result := g.Slice[int]{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result.Delete(1).Print()  // Slice[1, 3, 4, 5, 6, 7, 8, 9, 10]
	result.Delete(-9).Print() // Slice[1, 3, 4, 5, 6, 7, 8, 9, 10]

	result.SubSlice(1, -3).Print()            // Slice[2, 3, 4, 5, 6, 7]
	result.SubSlice(-3, result.Len()).Print() // Slice[8, 9, 10]
	result.SubSlice(-3, -1).Print()           // Slice[8, 9]
	result.SubSlice(-1, result.Len()).Print() // Slice[10]

	result = result.RandomSample(5)

	result.Clone().Append(999).Print()
	result.Print()
	fmt.Printf("%#v\n", result.Std())

	filled := g.NewSlice[int](10).Fill(88)
	filled.Print()

	slice = g.Slice[int]{1, 2, 3, 4, 5}.Print()

	slice.Cut(1, 3).Print()  // Slice[1, 4, 5]
	slice.Cut(-4, 3).Print() // Slice[1, 4, 5]
	slice.Cut(-3, 4).Print() // Slice[1, 2, 5]

	// InPlace Methods
	sipl := g.NewSlice[int]()

	sipl.AppendInPlace(1)
	sipl.AppendInPlace(2)
	sipl.AppendInPlace(3)

	sipl.DeleteInPlace(1)
	sipl.Fill(999999)

	sipl.InsertInPlace(0, 22, 33, 44)
	sipl.AddUniqueInPlace(22, 22, 22, 33, 44, 55)

	sipl.Print()

	slicea := g.Slice[string]{"a", "b", "c", "d"}
	slicea.InsertInPlace(2, "e", "f")
	slicea.Print()

	slice = g.Slice[int]{1, 2, 0, 4, 0, 3, 0, 0, 0, 0}
	slice = slice.Iter().Exclude(filters.IsZero).Collect()

	slice.DeleteInPlace(0)
	slice.Print()

	sll := g.NewSlice[int](0, 100000)
	sll = sll.Append(1).Clip()

	fmt.Println(sll.Cap())

	g.TransformSlice([]string{"AAA", "BBB"}, g.NewString).Iter().Map(g.String.Lower).Collect().Print()
	g.SliceOf([]string{"AAA", "BBB"}...).Iter().Map(strings.ToLower).Collect().Print()

	g.SliceOf(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11).Iter().
		Filter(isPrime).
		ForEach(func(n int) { fmt.Printf("%d is a prime number\n", n) })
}

func isPrime(n int) bool {
	for i := 2; i < n/2; i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}
