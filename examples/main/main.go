package main

import (
	"fmt"
	"github.com/enetx/g/rand"

	"github.com/enetx/g/rx"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	// strings
	str := NewString("") // declaration and assignation

	rand.String(9).Println()
	str.Hash().MD5().Println()

	str.Compress().Flate().Decompress().Flate().Unwrap().Println()

	NewString("12").TryInt().Ok().Println()

	str.Encode().JSON().Unwrap().Decode().JSON().Unwrap()

	fmt.Println(str.Decompress().Flate().Err())
	fmt.Println(str.Decompress().Flate().UnwrapOr("some value"))
	// fmt.Println(str.Decompress().Flate().Expect("some custom message on error"))
	// fmt.Println(str.Decompress().Flate().Unwrap())

	str = "*(&()&)(*&(*))"
	fmt.Println(str.Decode().Base64().Err())
	fmt.Println(str.Decode().Base64().UnwrapOr("some value"))
	// fmt.Println(str.Decode().Base64().Expect("some custom message on error"))
	// fmt.Println(str.Decode().Base64().Unwrap())

	var str2 String = "rest" // declaration and assignation

	fmt.Println(str2)

	a := NewString("abc")
	b := NewString("bbb")
	c := NewString("ccc")
	d := NewString("ddd")
	e := NewString("eee")

	str3 := SliceOf(rand.Choice(SliceOf(a.ReplaceAll("a", "zzz").Upper().Fields()...).Iter().Collect().Slice()).UnwrapOrDefault().Split("")...).Iter().Collect().Slice()[0].Lower().
		Std()

	fmt.Println(str3)

	// ints
	n := NewInt(52452356235) // declaration and assignation

	fmt.Printf("%v\n", n.BytesBE())

	n.Binary().Println()
	n.String().Println()

	rn := rand.N(NewInt(10))
	fmt.Println("random number: ", rn)

	rrn := rand.RangeInclusive(NewInt(10), 100)
	fmt.Println("random range number: ", rrn)

	var n2 Int = 321

	fmt.Println(n2) // declaration and assignation

	n.Add(n2).Mul(3).Println()

	// floats

	fl := NewFloat(12.5456)
	fl.Round().Println() // 13

	// slices
	sl := Slice[String]{a, b, c, d, e} // declaration and assignation

	rand.Shuffle(sl)

	fmt.Println(sl.Get(-1).Some() == "eee")
	fmt.Println(sl.Get(1).Some() == "bbb")
	fmt.Println(sl.Get(-2).Some() == "ddd")

	sl.Iter().Map(String.Upper).Collect().Slice().Println()

	slc := sl.Clone()
	slc.Push(sl...)
	slc.Push("ddd")

	// Typed counter: keys stay String, first-seen order, no any casts.
	counter := slc.Iter().CounterBy(func(s String) String { return s }).Collect().MapOrd[String, Int]()

	counter.SortBy(func(a, b Pair[String, Int]) cmp.Ordering {
		return b.Value.Cmp(a.Value).Then(a.Key.Cmp(b.Key))
	})

	counter.Println() // Output: MapOrd{ddd:3, abc:2, bbb:2, ccc:2, eee:2}

	counter.Iter().ForEach(func(k String, v Int) { fmt.Println(k.Title(), ":", v) })

	sl.Iter().ForEach(func(v String) { v.Println() })

	sl = sl.Iter().Unique().Collect().Slice()
	sl.Reverse()

	sl = sl.Iter().
		Filter(
			func(s String) bool {
				return s != "bbb"
			},
		).
		Collect().
		Slice()

	sl.Println()

	fmt.Println(rand.Choice(sl).UnwrapOrDefault())

	sl1 := SliceOf(1, 2, 3, 4, 5) // declaration and assignation

	fmt.Println(sl1.Iter().Fold(0, func(index, value int) int { return index + value })) // 15

	sl3 := Slice[String]{} // declaration and assignation
	sl3.Push("aaaaa", "bbbbb")

	fmt.Println(sl3.Last().Some().Count("b")) // 5

	sl4 := rand.Choice(SliceOf([]string{"root", "toor"}...)).UnwrapOrDefault()
	NewString(sl4).Upper().Println()

	sl3.Iter().Map(func(s String) String { return s + "MAPMAPMAP" }).Collect().Slice().Println()

	empsl := NewSlice[String]()
	fmt.Println(empsl.IsEmpty())

	// maps
	m1 := Map[int, string](map[int]string{1: "root", 22: "toor"}) // declaration and assignation
	m1.Iter().Values().Collect().Slice().Println()
	m1.Iter().Keys().Collect().Slice().Println()

	m2 := NewMap[int, string]() // declaration and assignation

	m2[99] = "AAA"
	m2[88] = "BBB"
	m2.Insert(77, "CCC")

	m2.Remove(99)
	m2.Println()
	m2.Iter().Keys().Collect().Slice().Println()

	m2.Println()
	fmt.Println(m2.Std())

	m3 := Map[string, string]{"test": "rest"} // declaration and assignation
	fmt.Println(m3.Contains("test"))

	ub := Bytes("abcdef\u0301\u031dg")
	ub.NormalizeNFC().Reverse().Println()

	NewString("abcde丂g").Reverse().Println()

	l := String("hello")
	l.Similarity("world").Println()

	hbs := Bytes("Hello, 世界!")
	hbs.Reverse().String().Println() // "!界世 ,olleH"

	hbs = Bytes("hello, world!")

	hbs.Replace([]byte("l"), []byte("L"), 2).String().Println() // "heLLo, world!"

	hs1 := String("kitten")
	hs2 := String("sitting")
	similarity := hs1.Similarity(hs2) // g.Float(57.14285714285714)

	similarity.Println()

	NewString("&aacute;").Decode().HTML().Println()

	to := String("Hello, 世界!")

	to.Encode().Hex().Println()
	to.Encode().Hex().Decode().Hex().Unwrap().Println()

	to.Encode().Octal().Println()
	to.Encode().Octal().Decode().Octal().Unwrap().Println()

	SliceOf(to.Encode().Binary().Chunks(8)...).Iter().Collect().Slice().Join(" ").Println()
	to.Encode().Binary().Decode().Binary().Unwrap().Println()

	toi := Int(1234567890)

	toi.Binary().Println()
	toi.Octal().Println()
	toi.Hex().Println()

	ascii := String("💛💚💙💜")
	fmt.Println(ascii.IsASCII())

	match := rx.New(`\w+`).Unwrap().IsMatch(NewString("some text"))
	fmt.Println(match)

	fmt.Println(String("example.com").EndsWithAny(".com", ".net"))

	NewString("Hello").Format("{} world").Println()
}
