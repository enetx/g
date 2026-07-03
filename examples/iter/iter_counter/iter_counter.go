package main

import (
	"fmt"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
	"github.com/enetx/g/f"
)

func main() {
	// ── Identity counting ─────────────────────────────────────────────────────
	// CounterBy with an identity fn counts the elements themselves. Keys keep
	// their real type, order is first-seen.
	slice := SliceOf(1, 2, 3, 1, 2, 1)

	slice.Iter().
		CounterBy(f.Id).
		Collect().
		Println() // MapOrd{1:3, 2:2, 3:1}

	// Sort a histogram by key descending:
	slice.Iter().
		CounterBy(f.Id).
		SortBy(func(a, b Pair[int, Int]) cmp.Ordering { return cmp.Cmp(b.Key, a.Key) }).
		Collect().
		Println() // MapOrd{3:1, 2:2, 1:3}

	// ── CounterBy: counting by a category key ────────────────────────────────
	// fn is a classifier, not a value transform: it answers "which bucket does
	// this element belong to". Elements whose keys collide are merged with
	// their counts summed. An injective fn (i*2, i+1) only relabels keys —
	// use the identity fn above for plain occurrence counting.

	// 1. HTTP statuses by class: how many 2xx / 4xx / 5xx
	statuses := SliceOf(200, 201, 200, 404, 500, 503, 200)

	statuses.Iter().
		CounterBy(func(s int) int { return s / 100 }).
		Collect().
		Println() // MapOrd{2:4, 4:1, 5:2}

	// 2. Words by length
	words := SliceOf[String]("go", "rust", "zig", "gleam", "c")

	words.Iter().
		CounterBy(String.Len).
		Collect().
		Println() // MapOrd{2:1, 4:1, 3:1, 5:1, 1:1}

	// 3. Emails by domain
	emails := SliceOf[String]("a@gmail.com", "b@mail.ru", "c@gmail.com", "d@gmail.com")

	emails.Iter().
		CounterBy(func(e String) String {
			return e.Split("@").Collect().Last().UnwrapOrDefault()
		}).
		Collect().
		Println() // MapOrd{gmail.com:3, mail.ru:1}

	// 4. Numbers by parity — the key can be any comparable type, even bool
	SliceOf(1, 2, 3, 4, 5, 6, 7).
		Iter().
		CounterBy(func(i int) bool { return i%2 == 0 }).
		Collect().
		Println() // MapOrd{false:4, true:3}

	// 5. Structs by field — count log entries per level, most frequent first
	type entry struct {
		level String
		msg   String
	}

	logs := SliceOf(
		entry{level: "INFO", msg: "started"},
		entry{level: "ERROR", msg: "boom"},
		entry{level: "INFO", msg: "listening"},
		entry{level: "WARN", msg: "slow query"},
		entry{level: "INFO", msg: "done"},
	)

	byLevel := logs.Iter().
		CounterBy(func(e entry) String { return e.level }).
		SortBy(func(a, b Pair[String, Int]) cmp.Ordering { return b.Value.Cmp(a.Value) }).
		Collect()

	byLevel.Println() // MapOrd{INFO:3, ERROR:1, WARN:1}

	// Top-1 category, fully typed end to end:
	if top := byLevel.Iter().First(); top.IsSome() {
		fmt.Println(top.Some().Key, "->", top.Some().Value) // INFO -> 3
	}
}
