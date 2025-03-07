package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	type response struct {
		Page   Int           `json:"page"`
		Fruits Slice[String] `json:"fruits"`
	}

	res := response{Page: 1, Fruits: Slice[String]{"apple", "peach", "pear"}}

	s := String("").Encode().JSON(res).Unwrap().Println()

	var res2 response

	s.Decode().JSON(&res2)
	fmt.Println(res2.Page, res2.Fruits.Get(-2))
}
