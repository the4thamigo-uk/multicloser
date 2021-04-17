package multicloser

import (
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

// WrapDefer executes `WrapDefer() on the singleton `MultiCloser`
func WrapDefer(f func()) {
	cls.WrapDefer(f)
}

// Close executes all the deferred functions in reverse order
func (m *MultiCloser) Close() (err error) {
	m.mtx.Lock()
	ff := m.ff
	m.ff = nil
	m.mtx.Unlock()

	for i := len(ff) - 1; i >= 0; i-- {
		f := ff[i]
		if e := f(); e != nil {
			err = multierror.Append(err, e)
		}
	}
	return err
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

// WrapDefer queues a function to be called in `Close()`
func (m *MultiCloser) WrapDefer(f func()) {
	m.Defer(wrap(f))
}

func wrap(f func()) func() error {
	return func() error {
		f()
		return nil
	}
}