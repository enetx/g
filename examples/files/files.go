package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	f := NewFile("some/dir/that/dont/exist/file.txt")

	f.Append("one").Unwrap().Append("\n")
	f.Append("two").Ok().Append("\n")

	f.Close()

	f.Read().Unwrap().Println()

	stat := f.Stat().Unwrap()
	fmt.Printf("Name(): %v\n", f.Name())
	fmt.Printf("IsDir(): %v\n", f.IsDir())
	fmt.Printf("f.IsLink(): %v\n", f.IsLink())
	fmt.Printf("Size(): %v\n", stat.Size())
	fmt.Printf("Mode(): %v\n", stat.Mode())
	fmt.Printf("ModeTime(): %v\n", stat.ModTime())

	fmt.Println(f.Exist())
	f.Dir().Unwrap().Path().Unwrap().Println()
	f.Path().Unwrap().Println()

	f = f.Rename("aaa/aaa/aaa/fff.txt").Ok().Copy(f.Dir().Ok().Join("copy_of_aaa.txt").Ok()).Ok()
	f.Name().Println()

	f.Ext().Println()
	f.MimeType().Unwrap().Println()

	fmt.Println("--------------------------------------------------------------")

	f.Path().Unwrap().Println()

	dir, file := f.Split()
	fmt.Println(dir.Path().Unwrap(), file.Name())

	newFile := NewFile(dir.Join(file.Name()).Ok())
	newFile.Path().Ok().Println()
}
