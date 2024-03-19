package main

import (
	"io"
	"log"
	"os"

	"github.com/enetx/g"
)

func main() {
	fn := RedirectLogOutput()
	defer fn()

	g.NewFile("file_not_exist.txt").Stat().Unwrap()
}

func RedirectLogOutput() func() {
	// Define the name of the log file
	logfile := "logfile.txt"

	// Open or create the log file in append mode with read and write permissions
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		panic(err)
	}

	// Save the existing stdout for later use
	out := os.Stdout

	// Create a MultiWriter that writes to both the original stdout and the log file
	mw := io.MultiWriter(out, f)

	// Create a pipe to capture stdout and stderr
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	// Replace stdout and stderr with the write end of the pipe
	// os.Stdout = w
	os.Stderr = w

	// Set the log package's output to the MultiWriter, so log.Print writes to both stdout and the log file
	log.SetOutput(mw)

	// Create a channel to control program exit
	exit := make(chan bool)

	go func() {
		// Copy all reads from the pipe to the MultiWriter, which writes to stdout and the log file
		io.Copy(mw, r)
		// Signal that copying is finished by sending true to the channel
		exit <- true
	}()

	// Return a function that can be deferred to clean up and ensure writes finish before program exit
	return func() {
		// Close the write end of the pipe and wait for copying to finish
		w.Close()
		<-exit

		// Close the log file after all writes have finished
		f.Close()
	}
}
