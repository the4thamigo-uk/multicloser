package main

import (
	"fmt"
	"github.com/the4thamigo-uk/multicloser"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {

	// close all resources managed by the global instance at the end of scope
	// Note that you can create your own instance if you prefer.
	defer multicloser.Close()

	// create a temporary file that will be deleted by the multicloser
	fn, err := createTempFile("hello world")
	if err != nil {
		return err
	}

	// read the file
	b, err := os.ReadFile(fn)
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}

func createTempFile(msg string) (string, error) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Write([]byte(msg))

	// we can use the multicloser in a similar way to testing.Cleanup(),
	// which means we can move logic from run() into separate functions.
	multicloser.Defer(func() error {
		return os.Remove(f.Name())
	})
	return f.Name(), nil
}
