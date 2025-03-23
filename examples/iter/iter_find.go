package main

import . "github.com/enetx/g"

func main() {
	SliceOf[Int](1, 1, 1, 3, 4, 4, 8, 8, 9, 9).
		Iter().
		Find(func(v Int) bool { return v%2 == 0 }).
		Some().
		Println() // 4

	m := Map[Int, Int]{1: 11, 2: 22, 3: 33}
	m.
		Iter().
		Find(func(_, v Int) bool { return v == 22 }).
		Some().
		Key.
		Println() // 2

	m.ToMapOrd().
		Iter().
		Find(func(_, v Int) bool { return v == 33 }).
		Some().
		Key.
		Println() // 3
}
