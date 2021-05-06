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

func run() (err error) {

	// close all resources managed by the multicloser instance at the end of scope
	mc := multicloser.New()
	defer func() { err = mc.Close() }()

	fn, err := createTempFile(mc, "hello world")
	if err != nil {
		return err
	}
	// no defer needed to clean up!

	f, err := openFile(mc, fn)
	if err != nil {
		return err
	}
	// no defer needed to clean up!

	// read from file
	b := make([]byte, 100)
	n, err := f.Read(b)
	if err != nil {
		return err
	}

	fmt.Println(string(b[:n]))
	return nil
}

func createTempFile(mc multicloser.Closer, msg string) (string, error) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Write([]byte(msg))

	// we can use the multicloser in a similar way to testing.Cleanup(),
	// which means we can move logic from run() into separate functions.
	// Here we delete the file when we have finished with it.
	mc.Defer(func() error {
		fmt.Printf("Removing temporary file %s\n", f.Name())
		return os.Remove(f.Name())
	})

	return f.Name(), nil
}

func openFile(mc multicloser.Closer, fn string) (*os.File, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	// we ensure we close the file when the multicloser closes
	mc.Deferf(func() error {
		fmt.Printf("Closing file for reading %s\n", f.Name())
		return f.Close()
	}, "failed to close file : %w")
	return f, nil
}
