package main

import (
	"fmt"
	"regexp"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	// strings
	str := NewString("") // declaration and assignation

	str.Random(9).Print()
	str.Hash().MD5().Print()

	str = "test"
	str.Compress().Flate().Decompress().Flate().Unwrap().Print()

	NewString("12").ToInt().Ok().Print()

	var jsonSet Set[int]

	str.Encode().JSON(SetOf(1, 2, 3, 4)).Unwrap().Decode().JSON(&jsonSet).Unwrap()

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

	str3 := a.ReplaceAll("a", "zzz").Upper().Fields().Collect().Random().Split("").Collect()[0].Lower().Std()

	fmt.Println(str3)

	// ints
	n := NewInt(52452356235) // declaration and assignation

	fmt.Printf("%v\n", n.Bytes())

	n.Hash().MD5().Print()
	n.Hash().SHA1().Print()
	n.Hash().SHA256().Print()

	n.Binary().Print()
	n.String().Print()

	rn := NewInt(10).Random()
	fmt.Println("random number: ", rn)

	rrn := NewInt(10).RandomRange(100)
	fmt.Println("random range number: ", rrn)

	var n2 Int = 321

	fmt.Println(n2) // declaration and assignation

	n.Add(n2).Mul(3).Print()

	// floats

	fl := NewFloat(12.5456)
	fl.Round().Print() // 13

	// slices
	sl := NewSlice[String]().Append(a, b, c, d, e) // declaration and assignation

	sl.Shuffle()

	fmt.Println(sl.Get(-1) == "eee")
	fmt.Println(sl.Get(1) == "bbb")
	fmt.Println(sl.Get(-2) == "ddd")

	sl.Iter().Map(String.Upper).Collect().Print()

	counter := sl.Append(sl...).Append("ddd").Iter().Counter().Collect()

	counter.SortBy(func(a, b Pair[String, Int]) cmp.Ordering {
		return b.Value.Cmp(a.Value).Then(a.Key.Cmp(b.Key))
	})

	counter.Print() // Output: MapOrd{ddd:3, abc:2, bbb:2, ccc:2, eee:2}

	counter.Iter().ForEach(func(k String, v Int) { fmt.Println(k.Title(), ":", v) })

	sl.Iter().ForEach(func(v String) { v.Print() })

	sl = sl.Iter().Unique().Collect()
	sl.Reverse()

	sl = sl.Iter().
		Filter(
			func(s String) bool {
				return s != "bbb"
			}).
		Collect()

	sl.Print()

	fmt.Println(sl.Random())

	sl1 := SliceOf(1, 2, 3, 4, 5) // declaration and assignation

	fmt.Println(sl1.Iter().Fold(0, func(index, value int) int { return index + value })) // 15

	sl3 := Slice[String]{} // declaration and assignation
	sl3 = sl3.Append("aaaaa", "bbbbb")

	fmt.Println(sl3.Last().Count("b")) // 5

	sl4 := SliceOf([]string{"root", "toor"}...).Random()
	NewString(sl4).Upper().Print()

	sl3.Iter().Map(func(s String) String { return s + "MAPMAPMAP" }).Collect().Print()

	empsl := NewSlice[String]()
	fmt.Println(empsl.Empty())

	// maps
	m1 := Map[int, string](map[int]string{1: "root", 22: "toor"}) // declaration and assignation
	m1.Iter().Values().Collect().Print()
	m1.Iter().Keys().Collect().Print()

	m2 := NewMap[int, string]() // declaration and assignation

	m2[99] = "AAA"
	m2[88] = "BBB"
	m2.Set(77, "CCC")

	m2.Delete(99).Print()
	m2.Iter().Keys().Collect().Print()

	m2.Print()
	fmt.Println(m2.Std())

	fmt.Println(m2.Invert().Iter().Values().Collect().Get(0))        // return int type
	fmt.Println(m2.Invert().Iter().Keys().Collect().Get(0).(string)) // return any type, need assert to type

	m3 := Map[string, string]{"test": "rest"} // declaration and assignation
	fmt.Println(m3.Contains("test"))

	ub := NewBytes("abcdef\u0301\u031dg")
	ub.NormalizeNFC().Reverse().Print()

	NewString("abcde丂g").Reverse().Print()

	l := String("hello")
	l.Similarity("world").Print()

	hbs := Bytes("Hello, 世界!")
	hbs.Reverse().String().Print() // "!界世 ,olleH"

	hbs = Bytes("hello, world!")

	hbs.Replace([]byte("l"), []byte("L"), 2).String().Print() // "heLLo, world!"

	hs1 := String("kitten")
	hs2 := String("sitting")
	similarity := hs1.Similarity(hs2) // g.Float(57.14285714285714)

	similarity.Print()

	NewString("&aacute;").Decode().HTML().Print()

	to := String("Hello, 世界!")

	to.Encode().Hex().Print()
	to.Encode().Hex().Decode().Hex().Unwrap().Print()

	to.Encode().Octal().Print()
	to.Encode().Octal().Decode().Octal().Unwrap().Print()

	to.Encode().Binary().Chunks(8).Join(" ").Print()
	to.Encode().Binary().Decode().Binary().Unwrap().Print()

	toi := Int(1234567890)

	toi.Binary().Print()
	toi.Octal().Print()
	toi.Hex().Print()

	ascii := String("💛💚💙💜")
	fmt.Println(ascii.IsASCII())

	reg := NewString("some text")
	fmt.Println(reg.RxMatch(regexp.MustCompile(`\w+`)))

	fmt.Println(String("example.com").EndsWithAny(".com", ".net"))

	NewString("Hello").Format("%s world").Print()
}
