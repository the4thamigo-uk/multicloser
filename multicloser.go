package multicloser

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"sync"
)

type (
	MultiCloser struct {
		ff  []func() error
		mtx sync.Mutex
	}
)

var cls MultiCloser

// Close executes `Close() on the singleton `MultiCloser`
func Close() (err error) {
	return cls.Close()
}

// Defer executes `Defer() on the singleton `MultiCloser`
func Defer(f func() error) {
	cls.Defer(f)
}

// Close executes all the deferred functions in reverse order.
// If a panic occurs when calling any of the deferred functions,
// the other functions will execute also.
func (m *MultiCloser) Close() (err error) {
	m.mtx.Lock()
	ff := m.ff
	m.ff = nil
	m.mtx.Unlock()

	for _, f := range ff {
		defer func(f func() error) {
			if e := f(); e != nil {
				err = multierror.Append(err, e)
			}
		}(f)
	}
	return
}

// Defer queues a function to be called in `Close()`.
// Passing `nil` to this function will cause a panic.
func (m *MultiCloser) Defer(f func() error) {

	if f == nil {
		panic("nil function indicates a programming error")
	}

	m.mtx.Lock()
	m.ff = append(m.ff, f)
	m.mtx.Unlock()
}

// Wrapf decorates the error returned from the function with the specified
// format string.
func Wrapf(f func() error, format string) func() error {
	return func() error {
		return fmt.Errorf(format, f())
	}
}

func Wrap(f func()) func() error {
	return func() error {
		f()
		return nil
	}
}
