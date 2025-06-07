package main

import (
	"sync"

	. "github.com/enetx/g"
	"github.com/enetx/g/box"
)

type Data struct {
	Counter int
}

func main() {
	Println("=== Using Box (thread-safe, copy-on-write) ===")

	b := box.New(&Data{Counter: 0})
	var wg sync.WaitGroup

	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Update(func(d *Data) *Data {
				cp := *d
				cp.Counter++
				return &cp
			})
		}()
	}

	wg.Wait()
	Println("Final counter (Box): {}", b.Load().Counter)

	// ------------------------------------------------

	Println("\n=== Without Box (unsafe, data race) ===")

	data := &Data{Counter: 0}
	wg = sync.WaitGroup{}

	for range 1000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data.Counter++ //  This is not thread-safe
		}()
	}

	wg.Wait()
	Println("Final counter (without Box): {}", data.Counter)

	// === Using Box (thread-safe, copy-on-write) ===
	// Final counter (Box): 1000
	//
	// === Without Box (unsafe, data race) ===
	// Final counter (without Box): 943
}
