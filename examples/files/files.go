package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
)

func main() {
	// f := g.NewFile("some/dir/that/dont/exist/file.txt")

	// // words := f.Iter().Unwrap().ByWords()
	// // runes := f.Iter().Unwrap().ByRunes()
	// // bytes := f.Iter().Unwrap().ByBytes()
	// lines := f.Iter().Unwrap().ByLines()

	// for lines.Next() {
	// 	fmt.Println(lines.ToString())
	// }

	// if lines.Err() != nil {
	// 	fmt.Println(lines.Err())
	// }

	// or

	// f := g.NewFile("somebigfile.txt")

	// dbg.Dbg(f.Name())

	// for line := f.Iter().Unwrap().ByLines(); line.Next(); {
	// 	// for line := f.Iterr().Expect().ByLines(); line.Next(); {
	// 	fmt.Println(line.ToString())
	// }

	//////////////////////////////////////////////////

	f := g.NewFile("some/dir/that/dont/exist/file.txt")

	f.Append("one").Unwrap().Append("\n")
	f.Append("two").Ok().Append("\n")

	f.Close()

	f.Read().Unwrap().Print()

	stat := f.Stat().Unwrap()
	fmt.Printf("Name(): %v\n", f.Name())
	fmt.Printf("IsDir(): %v\n", stat.IsDir())
	fmt.Printf("Size(): %v\n", stat.Size())
	fmt.Printf("Mode(): %v\n", stat.Mode())
	fmt.Printf("ModeTime(): %v\n", stat.ModTime())

	fmt.Println(f.Exist())
	f.Dir().Unwrap().Path().Unwrap().Print()
	f.Path().Unwrap().Print()

	f = f.Rename("aaa/aaa/aaa/fff.txt").Ok().Copy(f.Dir().Ok().Join("copy_of_aaa.txt").Ok()).Ok()
	f.Name().Print()

	f.Ext().Print()
	f.MimeType().Unwrap().Print()

	fmt.Println("--------------------------------------------------------------")

	f.Path().Unwrap().Print()

	dir, file := f.Split()
	fmt.Println(dir.Path().Unwrap(), file.Name())

	newFile := g.NewFile(dir.Join(file.Name()).Ok())
	newFile.Path().Ok().Print()
}
