package main

import (
	"encoding/json"
	"fmt"

	. "github.com/enetx/g"
)

func main() {
	type User struct {
		Name   String                `json:"name"`
		Age    Option[Int]           `json:"age,omitzero"`
		Email  Option[String]        `json:"email,omitempty"`
		Tags   Option[Slice[String]] `json:"tags,omitzero"`
		Scores Option[Slice[Int]]    `json:"scores,omitempty"`
	}

	// All fields present
	alice := User{
		Name:   "Alice",
		Age:    Some[Int](30),
		Email:  Some[String]("alice@example.com"),
		Tags:   Some(SliceOf[String]("admin", "active")),
		Scores: Some(SliceOf[Int](95, 87, 100)),
	}

	data, _ := json.Marshal(alice)
	fmt.Println(string(data))
	// {"name":"Alice","age":30,"email":"alice@example.com","tags":["admin","active"],"scores":[95,87,100]}

	// None fields: omitzero omits the field, omitempty outputs null
	bob := User{
		Name:   "Bob",
		Age:    None[Int](),           // omitzero → field omitted
		Email:  None[String](),        // omitempty → "email":null
		Tags:   None[Slice[String]](), // omitzero → field omitted
		Scores: None[Slice[Int]](),    // omitempty → "scores":null
	}

	data, _ = json.Marshal(bob)
	fmt.Println(string(data))
	// {"name":"Bob","email":null,"scores":null}

	// Unmarshal with slices
	input := []byte(`{"name":"Carol","age":28,"tags":["user","new"],"scores":[80,90]}`)

	var carol User
	json.Unmarshal(input, &carol)

	fmt.Println(carol.Tags)            // Some(Slice[user, new])
	fmt.Println(carol.Scores)          // Some(Slice[80, 90])
	fmt.Println(carol.Email)           // None
	fmt.Println(carol.Tags.Unwrap())   // Slice[user, new]
	fmt.Println(carol.Scores.Unwrap()) // Slice[80, 90]
	fmt.Println(carol.Email.IsNone())  // true

	// Unmarshal null and missing fields
	input2 := []byte(`{"name":"Dan","email":null,"scores":null}`)

	var dan User
	json.Unmarshal(input2, &dan)

	fmt.Println(dan.Age)    // None
	fmt.Println(dan.Email)  // None
	fmt.Println(dan.Tags)   // None
	fmt.Println(dan.Scores) // None
}
