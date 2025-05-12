package main

import (
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	// Create a new ordered map and get an entry for "alpha"
	mo := NewMapOrd[string, int]()

	fmt.Println("=== MapOrdEntry Examples ===")

	// 1) OrSet: insert 10 if "alpha" is not present
	mo.Entry("alpha").OrSet(10)
	fmt.Println(`After OrSet("alpha", 10):`, mo)

	// 2) Transform: multiply the existing "alpha" value by 2
	mo.Entry("alpha").Transform(func(v *int) { *v *= 2 })
	fmt.Println(`After Transform(*2) on "alpha":`, mo)

	// 3) OrSetBy: lazy insertion won't run since "alpha" already exists
	mo.Entry("alpha").OrSetBy(func() int {
		fmt.Println("This won't run because the key exists")
		return 99
	})
	fmt.Println(`After OrSetBy on existing "alpha":`, mo)

	// 4) OrDefault: insert the zero value (0) for new key "beta"
	mo.Entry("beta").OrDefault()
	fmt.Println(`After OrDefault on "beta":`, mo)

	// 5) Set: unconditionally set "beta" to 42
	mo.Entry("beta").Set(42)
	fmt.Println(`After Set("beta", 42):`, mo)

	// 6) Get: retrieve the value as an Option and print if present
	if opt := mo.Entry("alpha").Get(); opt.IsSome() {
		fmt.Println(`Get("alpha"):`, opt.Unwrap())
	}

	// 7) Delete: remove "alpha" from the map
	mo.Entry("alpha").Delete()
	fmt.Println(`After Delete("alpha"):`, mo)
}
