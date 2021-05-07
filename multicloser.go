// multicloser provides a way to defer a number of functions (func() error), to be
// executed when the Close() function of the multiCloser is called. Errors returned from
// any of the deferred function invocations, are merged into a single error returned by the
// Close() method.
//
// The intent of this library is to provide a capability that is similar to the testing.Cleanup()
// mechanism, so you can defer functions to a scope, other than at the end of the current function.
// This allows you to move the creation and cleanup logic for resources, into separate functions
// and this helps to keep your higher level functions (e.g. main()) simpler. Additionally it avoids
// having to deal with errors in multiple defer commands.
package multicloser

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"sync"
)

type (
	Closer interface {
		Close() error
		Defer(func() error)
		Deferf(func() error, string)
	}

	multiCloser struct {
		ff  []func() error
		mtx sync.Mutex
	}
)

// New creates a new MultiCloser
func New() Closer {
	return new(multiCloser)
}

// Close executes all the deferred functions in reverse order.
// If a panic occurs when calling any of the deferred functions,
// the other functions will execute also.
func (m *multiCloser) Close() (err error) {
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

// Defer queues a function to be called in Close().
// Passing nil to this function will cause a panic.
func (m *multiCloser) Defer(f func() error) {

	if f == nil {
		panic("nil function indicates a programming error")
	}

	m.mtx.Lock()
	m.ff = append(m.ff, f)
	m.mtx.Unlock()
}

// Deferf queues a function to be called in Close(),
// but wraps any resulting error with the provided format string.
func (m *multiCloser) Deferf(f func() error, format string) {
	m.Defer(Wrapf(f, format))
}

// Wrapf decorates the error returned from the function with the specified
// format string.
func Wrapf(f func() error, format string) func() error {
	return func() error {
		err := f()
		if err != nil {
			return fmt.Errorf(format, err)
		}
		return nil
	}
}

// Wrap lifts a function that does not return an error into one that returns nil
func Wrap(f func()) func() error {
	return func() error {
		f()
		return nil
	}
}
