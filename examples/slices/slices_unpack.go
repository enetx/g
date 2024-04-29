package main

import (
	"fmt"

	"github.com/enetx/g"
)

func main() {
	// Example 1: Using the Split method
	account := g.String("e@mail.com:password")

	// Define variables to store the result of splitting
	var mail, password g.String

	// Split the account string by ":" and unpack the result into mail and password variables
	account.Split(":").Collect().Unpack(&mail, &password)

	// Print the result
	fmt.Println(mail, password)

	// Example 2: Using the Unpack method with a slice
	numbers := g.Slice[int]{1, 2, 3, 4, 5}

	// Define variables to store the unpacked values
	var a, b, c int

	// Unpack the first three elements of the numbers slice into variables a, b, and c
	numbers.Unpack(&a, &b, &c)

	// Print the result
	fmt.Println(a, b, c)

	// Example 3: Using the Unpack method with a slice of custom struct
	// Define a custom struct
	type character struct {
		Name string
		Age  int
	}

	// Create a slice of custom struct
	characters := g.Slice[character]{
		{Name: "Tom", Age: 6},
		{Name: "Jerry", Age: 2},
	}

	// Define variables to store the unpacked values
	var tom, jerry character

	// Unpack the first two elements of the characters slice into variables tom and jerry
	characters.Unpack(&tom, &jerry)

	// Print the names of tom and jerry
	fmt.Println(tom.Name, jerry.Name)
}
