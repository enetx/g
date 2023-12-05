package main

import (
	"fmt"

	"gitlab.com/x0xO/g"
)

func main() {
	var (
		file       = g.NewFile("somebigfile.txt")
		position   int64
		content    g.String
		lineToRead = 10
	)

	position, content = file.SeekToLine(position, lineToRead)

	fmt.Println(position)
	fmt.Println(content)
}
