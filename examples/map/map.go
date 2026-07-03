package main

import (
	. "github.com/enetx/g"
	"github.com/enetx/g/f"
)

func main() {
	gos2 := NewMap[int, Slice[int]]()

	for i := range 5 {
		gos2.Insert(i, gos2.Get(i).UnwrapOrDefault().Append(i))
	}

	for i := range 10 {
		gos2.Insert(i, gos2.Get(i).UnwrapOrDefault().Append(i))
	}

	gos2.Println()

	//////////////////////////////////////////////////////////////////////////

	god := NewMap[int, Slice[int]]()

	for i := range 10 {
		god[i] = god.Get(i).UnwrapOrDefault().Append(i)
	}

	for i := range 10 {
		god[i] = god.Get(i).UnwrapOrDefault().Append(i)
	}

	god.Println()

	// Map iterators can change BOTH key and value types in one pass (Go 1.27)
	users := Map[int, String]{1: "alice", 2: "bob"}
	users.Iter().
		Map(func(id int, name String) (String, int) { return name, id }). // swap key/value types
		Collect().
		Println() // Map{alice:1, bob:2} (order may vary)

	// FilterByValue with a one-place f predicate — no lambda wrapper needed
	scores := Map[String, Int]{"alice": 90, "bob": 65, "carol": 78}
	scores.Iter().
		FilterByValue(f.Gt(Int(70))).
		Collect().
		Println() // Map{alice:90, carol:78} (order may vary)

}
