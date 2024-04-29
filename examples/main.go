package main

import (
	"fmt"

	"github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	// strings
	str := g.NewString("") // declaration and assignation

	str.Random(9).Print()
	str.Hash().MD5().Print()

	str = "test"
	str.Comp().Flate().Decomp().Flate().Unwrap().Print()

	g.NewString("12").ToInt().Ok().Print()

	var jsonSet g.Set[int]

	str.Enc().JSON(g.SetOf(1, 2, 3, 4)).Unwrap().Dec().JSON(&jsonSet).Unwrap()

	fmt.Println(str.Decomp().Flate().Err())
	fmt.Println(str.Decomp().Flate().UnwrapOr("some value"))
	// fmt.Println(str.Dec().Flate().Expect("some custom message on error"))
	// fmt.Println(str.Dec().Flate().Unwrap())

	str = "*(&()&)(*&(*))"
	fmt.Println(str.Dec().Base64().Err())
	fmt.Println(str.Dec().Base64().UnwrapOr("some value"))
	// fmt.Println(str.Dec().Base64().Expect("some custom message on error"))
	// fmt.Println(str.Dec().Base64().Unwrap())

	var str2 g.String = "rest" // declaration and assignation

	fmt.Println(str2)

	a := g.NewString("abc")
	b := g.NewString("bbb")
	c := g.NewString("ccc")
	d := g.NewString("ddd")
	e := g.NewString("eee")

	str3 := a.ReplaceAll("a", "zzz").Upper().Fields().Collect().Random().Split("").Collect()[0].Lower().Std()

	fmt.Println(str3)

	// ints
	n := g.NewInt(52452356235) // declaration and assignation

	fmt.Printf("%v\n", n.Bytes())

	n.Hash().MD5().Print()
	n.Hash().SHA1().Print()
	n.Hash().SHA256().Print()

	n.ToBinary().Print()
	n.ToString().Print()

	rn := g.NewInt(10).Random()
	fmt.Println("random number: ", rn)

	rrn := g.NewInt(10).RandomRange(100)
	fmt.Println("random range number: ", rrn)

	var n2 g.Int = 321

	fmt.Println(n2) // declaration and assignation

	n.Add(n2).Mul(3).Print()

	// floats

	fl := g.NewFloat(12.5456)
	fl.Round().Print() // 13

	// slices
	sl := g.NewSlice[g.String]().Append(a, b, c, d, e) // declaration and assignation

	sl.Shuffle()

	fmt.Println(sl.Get(-1) == "eee")
	fmt.Println(sl.Get(1) == "bbb")
	fmt.Println(sl.Get(-2) == "ddd")

	sl.Iter().Map(g.String.Upper).Collect().Print()

	counter := sl.Append(sl...).Append("ddd").Iter().Counter().Collect()

	counter.SortBy(func(a, b g.Pair[g.String, g.Int]) cmp.Ordering {
		return b.Value.Cmp(a.Value).Then(a.Key.Cmp(b.Key))
	})

	counter.Print() // Output: MapOrd{ddd:3, abc:2, bbb:2, ccc:2, eee:2}

	counter.Iter().ForEach(func(k g.String, v g.Int) { fmt.Println(k.Title(), ":", v) })

	sl.Iter().ForEach(func(v g.String) { v.Print() })

	sl = sl.Iter().Unique().Collect()
	sl.Reverse()

	sl = sl.Iter().
		Filter(
			func(s g.String) bool {
				return s != "bbb"
			}).
		Collect()

	sl.Print()

	fmt.Println(sl.Random())

	sl1 := g.SliceOf(1, 2, 3, 4, 5) // declaration and assignation

	fmt.Println(sl1.Iter().Fold(0, func(index, value int) int { return index + value })) // 15

	sl3 := g.Slice[g.String]{} // declaration and assignation
	sl3 = sl3.Append("aaaaa", "bbbbb")

	fmt.Println(sl3.Last().Count("b")) // 5

	sl4 := g.SliceOf([]string{"root", "toor"}...).Random()
	g.NewString(sl4).Upper().Print()

	sl3.Iter().Map(func(s g.String) g.String { return s + "MAPMAPMAP" }).Collect().Print()

	empsl := g.NewSlice[g.String]()
	fmt.Println(empsl.Empty())

	// maps
	m1 := g.MapFromStd(map[int]string{1: "root", 22: "toor"}) // declaration and assignation
	m1.Iter().Values().Collect().Print()
	m1.Iter().Keys().Collect().Print()

	m2 := g.NewMap[int, string]() // declaration and assignation

	m2[99] = "AAA"
	m2[88] = "BBB"
	m2.Set(77, "CCC")

	m2.Delete(99).Print()
	m2.Iter().Keys().Collect().Print()

	m2.Print()
	fmt.Println(m2.Std())

	fmt.Println(m2.Invert().Iter().Values().Collect().Get(0))        // return int type
	fmt.Println(m2.Invert().Iter().Keys().Collect().Get(0).(string)) // return any type, need assert to type

	m3 := g.Map[string, string]{"test": "rest"} // declaration and assignation
	fmt.Println(m3.Contains("test"))

	ub := g.NewBytes("abcdef\u0301\u031dg")
	ub.NormalizeNFC().Reverse().Print()

	g.NewString("abcde丂g").Reverse().Print()

	l := g.String("hello")
	l.Similarity("world").Print()

	hbs := g.Bytes("Hello, 世界!")
	hbs.Reverse().ToString().Print() // "!界世 ,olleH"

	hbs = g.Bytes("hello, world!")

	hbs.Replace([]byte("l"), []byte("L"), 2).ToString().Print() // "heLLo, world!"

	hs1 := g.String("kitten")
	hs2 := g.String("sitting")
	similarity := hs1.Similarity(hs2) // g.Float(57.14285714285714)

	similarity.Print()

	g.NewString("&aacute;").Dec().HTML().Print()

	to := g.String("Hello, 世界!")

	to.Enc().Hex().Print()
	to.Enc().Hex().Dec().Hex().Unwrap().Print()

	to.Enc().Octal().Print()
	to.Enc().Octal().Dec().Octal().Unwrap().Print()

	to.Enc().Binary().Chunks(8).Join(" ").Print()
	to.Enc().Binary().Dec().Binary().Unwrap().Print()

	toi := g.Int(1234567890)

	toi.ToBinary().Print()
	toi.ToOctal().Print()
	toi.ToHex().Print()

	ascii := g.String("💛💚💙💜")
	fmt.Println(ascii.IsASCII())

	reg := g.NewString("some text")
	fmt.Println(reg.ContainsRegexp(`\w+`).Unwrap())

	fmt.Println(g.String("example.com").EndsWith(".com", ".net"))

	g.NewString("Hello").Format("%s world").Print()
}
