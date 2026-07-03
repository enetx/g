package main

import (
	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	si := SliceOf(1, 2, 3).Iter()
	Println("si.Next(): {}", si.Next())
	Println("si.Next(): {}", si.Next())
	Println("si.Next(): {}", si.Next())
	Println("si.Next(): {}", si.Next())

	sei := SetOf(1, 2, 3, 3, 3, 1, 2, 2, 2).Iter()
	Println("sei.Next(): {}", sei.Next())
	Println("sei.Next(): {}", sei.Next())
	Println("sei.Next(): {}", sei.Next())
	Println("sei.Next(): {}", sei.Next())

	di := DequeOf(1, 2, 3).Iter()
	Println("di.Next(): {}", di.Next())
	Println("di.Next(): {}", di.Next())
	Println("di.Next(): {}", di.Next())
	Println("di.Next(): {}", di.Next())

	h := NewHeap(cmp.Cmp[int])
	h.Push(1)
	h.Push(2)
	h.Push(3)

	hi := h.Iter()
	Println("hi.Next(): {}", hi.Next())
	Println("hi.Next(): {}", hi.Next())
	Println("hi.Next(): {}", hi.Next())
	Println("hi.Next(): {}", hi.Next())

	m := NewMap[string, string]()
	m.Insert("a", "aa")
	m.Insert("b", "bb")
	m.Insert("c", "cc")

	mi := m.Iter()
	Println("mi.Next(): {}", mi.Next())
	Println("mi.Next(): {}", mi.Next())
	Println("mi.Next(): {}", mi.Next())
	Println("mi.Next(): {}", mi.Next())

	mo := NewMapOrd[string, string]()
	mo.Insert("a", "aa")
	mo.Insert("b", "bb")
	mo.Insert("c", "cc")

	moi := mo.Iter()
	Println("moi.Next(): {}", moi.Next())
	Println("moi.Next(): {}", moi.Next())
	Println("moi.Next(): {}", moi.Next())
	Println("moi.Next(): {}", moi.Next())

	ri := SeqResult[int](func(yield func(Result[int]) bool) {
		yield(Ok(1))
		yield(Ok(2))
		yield(Ok(3))
		yield(Err[int](nil))
	})

	Println("ri.Next(): {}", ri.Next())
	Println("ri.Next(): {}", ri.Next())
	Println("ri.Next(): {}", ri.Next())
	Println("ri.Next(): {}", ri.Next())
	Println("ri.Next(): {}", ri.Next())
}
