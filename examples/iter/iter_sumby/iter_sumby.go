package main

import (
	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
	"github.com/enetx/g/f"
)

func main() {
	// SumBy projects every element to a numeric value and adds them up. The
	// result type is chosen by the projection, independent of the element type.

	// Slice: sum the lengths of the words (String -> Int)
	SliceOf[String]("a", "bb", "ccc", "dddd").
		Iter().
		SumBy(func(s String) Int { return s.Len() }).
		Println() // 10

	// Slice with a type-changing projection: Int elements -> Float sum
	SliceOf[Int](1, 2, 3, 4).
		Iter().
		SumBy(func(v Int) Float { return Float(v) / 2 }).
		Println() // 5

	// Set: order is undefined, but summation is commutative
	SetOf[Int](10, 20, 30).
		Iter().
		SumBy(f.Id).
		Println() // 60

	// Deque
	DequeOf[Int](1, 2, 3).
		Iter().
		SumBy(func(v Int) Int { return v * v }).
		Println() // 14

	// Heap
	HeapOf(cmp.Cmp[Int], 5, 1, 3).
		Iter().
		SumBy(f.Id).
		Println() // 9

	// MapOrd: the projection receives both key and value
	mo := NewMapOrd[String, Int]()
	mo.Insert("a", 1)
	mo.Insert("b", 2)
	mo.Insert("c", 3)
	mo.Iter().
		SumBy(func(_ String, v Int) Int { return v }).
		Println() // 6

	// Zipped pairs (SeqPairs): dot-product of two slices
	xs := SliceOf[Int](1, 2, 3)
	ys := SliceOf[Int](4, 5, 6)
	xs.Iter().
		Zip(ys.Iter()).
		SumBy(func(x, y Int) Int { return x * y }).
		Println() // 32

	// SeqResult via TryMap: parse fallibly, then SumBy short-circuits on Err.
	// TryMap enters the Result pipeline.
	ok := SliceOf[String]("1", "2", "3").
		Iter().
		TryMap(String.TryInt). // SeqResult[Int]
		SumBy(f.Id)
	Println("{}", ok) // Ok(6)

	bad := SliceOf[String]("1", "x", "3").
		Iter().
		TryMap(String.TryInt).
		SumBy(f.Id)

	Println("{}", bad) // Err(invalid integer: "x")
}
