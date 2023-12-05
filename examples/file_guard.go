package main

import (
	"time"

	"gitlab.com/x0xO/g"
)

func main() {
	exit := make(chan struct{})

	// Specify the file name
	fname := g.String("test_file.txt")

	// Create and guard the file
	f := g.NewFile(fname).Guard()

	// Goroutine to release the guard after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		f.Close()
	}()

	// Goroutine to read and print from the file
	go func() {
		// Create a new file and guard it for reading
		g.NewFile(fname).Guard().Read().Unwrap().Print()
		exit <- struct{}{}
	}()

	// Append data to the original file
	f.Append("test string")

	// Wait for the exit signal
	<-exit

	// Delete the file
	f.Remove()
}
