package main

import (
	"fmt"

	"github.com/enetx/g"
)

func main() {
	type response struct {
		Page   g.Int           `json:"page"`
		Fruits g.Slice[string] `json:"fruits"`
	}

	res := response{Page: 1, Fruits: g.SliceOf("apple", "peach", "pear")}

	s := g.NewString("").Encode().JSON(res).Unwrap().Append("\n").Print()

	var res2 response

	s.Decode().JSON(&res2)
	fmt.Println(res.Page, res.Fruits.Get(-2))
}
