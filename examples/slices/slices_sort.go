package main

import (
	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	type Order struct {
		Product  String
		Customer String
		Price    Float
	}

	orders := Slice[Order]{
		{"foo", "alice", 1.458},
		{"bar", "bob", 3.256},
		{"baz", "carol", 4.391},
		{"foo", "alice", 2.681},
		{"bar", "carol", 1.866},
		{"foo", "bob", 4.825},
	}

	// Sort by customer first, product second, and last by higher price
	orders.SortBy(func(a, b Order) cmp.Ordering {
		return a.Customer.Cmp(b.Customer).
			Then(a.Product.Cmp(b.Product)).
			Then(b.Price.Cmp(a.Price))
	})

	orders.Iter().ForEach(func(v Order) {
		Print("{} {} {.RoundDecimal(2)}\n", v.Product, v.Customer, v.Price)
	})

	// Output:

	// foo alice 2.68
	// foo alice 1.46
	// bar bob 3.26
	// foo bob 4.83
	// bar carol 1.87
	// baz carol 4.39
}
