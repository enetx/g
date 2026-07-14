package main

import (
	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
	"github.com/enetx/g/f"
)

func main() {
	// TryMap applies a fallible transform to each element and enters the Result
	// pipeline (SeqResult[U]). The result type U is chosen by the transform, so a
	// single primitive can be parsed into many different target types. The first
	// Err short-circuits the terminal (TryCollect / SumBy / ...).

	// String -> Int
	i := SliceOf[String]("1", "2", "3").
		Iter().
		TryMap(String.TryInt). // SeqResult[Int]
		TryCollect()           // Result[Slice[Int]]
	Println("Int:     {}", i) // Ok(Slice[1, 2, 3])

	// String -> Float, then reduce with SumBy
	fl := SliceOf[String]("1.5", "2.25", "0.25").
		Iter().
		TryMap(String.TryFloat). // SeqResult[Float]
		SumBy(f.Id)              // Result[Float]
	Println("Float:   {}", fl) // Ok(4)

	// String -> bool
	b := SliceOf[String]("true", "0", "T", "false").
		Iter().
		TryMap(String.TryBool). // SeqResult[bool]
		TryCollect()            // Result[Slice[bool]]
	Println("Bool:    {}", b) // Ok(Slice[true, false, true, false])

	// String -> uint
	u := SliceOf[String]("10", "0xff", "0b1010").
		Iter().
		TryMap(String.TryUint). // SeqResult[uint]
		TryCollect()            // Result[Slice[uint]]
	Println("Uint:    {}", u) // Ok(Slice[10, 255, 10])

	// Error path: the first bad element short-circuits the whole batch.
	bad := SliceOf[String]("1", "x", "3").
		Iter().
		TryMap(String.TryInt).
		TryCollect()
	Println("Err:     {}", bad) // Err(invalid integer: "x")

	// From a MapOrd: fn receives (K, V) and collapses them into a single Result[U].
	// Here U is String — the value is parsed and recombined with its key.
	mo := NewMapOrd[String, String]()
	mo.Insert("a", "10")
	mo.Insert("b", "20")
	kv := mo.Iter().
		TryMap(func(k, v String) Result[String] {
			return v.TryInt().Map(func(n Int) String { return k + "=" + (n * 2).String() })
		}).
		TryCollect() // Result[Slice[String]]
	Println("MapOrd:  {}", kv) // Ok(Slice[a=20, b=40])

	// From a Set, summing parsed values.
	sset := SetOf[String]("3", "4", "5").
		Iter().
		TryMap(String.TryInt).
		SumBy(f.Id) // Result[Int]
	Println("Set sum: {}", sset) // Ok(12)

	// From a Heap.
	h := HeapOf(cmp.Cmp[String], "2", "1", "3").
		Iter().
		TryMap(String.TryInt).
		TryCollect()
	Println("Heap:    {}", h) // Ok(Slice[1, 2, 3])

	// Map -> Map: keep the key by choosing U = Pair[K, U2], then rebuild the map.
	// TryMap yields Result[Pair], TryCollect gathers Result[Slice[Pair]], and the
	// *FromPairs constructors go straight into Result.Map as first-class functions,
	// so the whole "fallible map -> map" reads as one short chain.
	src := NewMapOrd[String, String]()
	src.Insert("a", "10")
	src.Insert("b", "20")

	toPairs := func(k, v String) Result[Pair[String, Int]] {
		return v.TryInt().Map(func(n Int) Pair[String, Int] { return PairOf(k, n) })
	}

	// into an (unordered) Map
	resMap := src.Iter().TryMap(toPairs).TryCollect().Map(MapFromPairs)
	Println("-> Map:     {}", resMap) // Ok(Map{a:10, b:20})

	// into a MapOrd — insertion order is preserved through the whole chain
	resOrd := src.Iter().TryMap(toPairs).TryCollect().Map(MapOrdFromPairs)
	Println("-> MapOrd:  {}", resOrd) // Ok(MapOrd{a:10, b:20})

	// into a concurrent MapSafe
	resSafe := src.Iter().TryMap(toPairs).TryCollect().Map(MapSafeFromPairs)
	Println("-> MapSafe: {}", resSafe) // Ok(MapSafe{a:10, b:20})

	// The first Err short-circuits the whole rebuild, whichever target is used.
	src.Insert("c", "x")
	resErr := src.Iter().TryMap(toPairs).TryCollect().Map(MapFromPairs)
	Println("-> Err:     {}", resErr) // Err(invalid integer: "x")
}
