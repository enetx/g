package main

import (
	"time"

	"gitlab.com/x0xO/g/pkg/rand"
)

func main() {
	rand.N(100 * time.Millisecond)
}
