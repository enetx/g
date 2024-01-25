package main

import (
	"fmt"
	"time"

	"gitlab.com/x0xO/g"
)

func main() {
	// Create a channel to signal the main goroutine when the other goroutine completes
	exit := make(chan struct{})

	// Specify the file name
	fname := g.String("test_file.txt")

	// Create and guard the file
	f := g.NewFile(fname).Guard()

	// Append data to the original file
	f.Append("test string")

	// Goroutine to read and print from the file
	go func() {
		fmt.Println("Waiting for guard release")
		// Create a new file and guard it for reading
		g.NewFile(fname).Guard().Read().Unwrap().Print()
		// Signal the main goroutine that the reading is complete
		exit <- struct{}{}
	}()

	// Goroutine to release the guard after 2 seconds
	go func() {
		// Sleep for 2 seconds
		time.Sleep(2 * time.Second)
		// Close the file, releasing the guard
		f.Close()

		fmt.Println("Guard released")
	}()

	// Wait for the exit signal
	<-exit

	// Delete the file
	f.Remove()
}
