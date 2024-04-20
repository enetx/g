package main

import (
	"fmt"

	"github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	type Order struct {
		Product  g.String
		Customer g.String
		Price    g.Float
	}

	orders := g.Slice[Order]{
		{"foo", "alice", 1.00},
		{"bar", "bob", 3.00},
		{"baz", "carol", 4.00},
		{"foo", "alice", 2.00},
		{"bar", "carol", 1.00},
		{"foo", "bob", 4.00},
	}

	// Sort by customer first, product second, and last by higher price
	orders.SortBy(func(a, b Order) cmp.Ordering {
		return a.Customer.Cmp(b.Customer).
			Then(a.Product.Cmp(b.Product)).
			Then(b.Price.Cmp(a.Price))
	})

	orders.Iter().ForEach(func(v Order) {
		fmt.Printf("%s %s %.2f\n", v.Product, v.Customer, v.Price)
	})

	// Output:

	// foo alice 2.00
	// foo alice 1.00
	// bar bob 3.00
	// foo bob 4.00
	// bar carol 1.00
	// baz carol 4.00
}
