package multicloser

import (
	"errors"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
	"testing"
)

func testEmpty(t *testing.T) {
	var m MultiCloser
	err := m.Close()
	require.NoError(t, err)
}

func TestNil(t *testing.T) {
	var m MultiCloser
	require.Panics(t, func() { m.Defer(nil) })
}

func TestClose(t *testing.T) {
	var m MultiCloser

	var ii []int
	m.Defer(func() error {
		ii = append(ii, 1)
		return nil
	})
	m.Defer(func() error {
		ii = append(ii, 2)
		return nil
	})
	m.Defer(func() error {
		ii = append(ii, 3)
		return nil
	})

	err := m.Close()
	require.NoError(t, err)
	require.Equal(t, []int{3, 2, 1}, ii)

	// second time should not run the functions
	ii = nil
	err = m.Close()
	require.NoError(t, err)
	require.Nil(t, ii)
}

func TestCloseErrors(t *testing.T) {
	var m MultiCloser

	err1 := errors.New("1")
	err2 := errors.New("2")
	merr := multierror.Append(err2, err1)

	m.Defer(func() error {
		return err1
	})
	m.Defer(func() error {
		return nil
	})
	m.Defer(func() error {
		return err2
	})

	err := m.Close()
	require.Equal(t, merr, err)
}

func TestClosePanic(t *testing.T) {
	var m MultiCloser

	var ii []int
	m.Defer(func() error {
		ii = append(ii, 1)
		return nil
	})
	m.Defer(func() error {
		ii = append(ii, 2)
		panic("")
	})
	m.Defer(func() error {
		ii = append(ii, 3)
		return nil
	})

	require.Panics(t, func() { m.Close() })
	require.Equal(t, []int{3, 2, 1}, ii)
}

func TestWrap(t *testing.T) {
	var i int
	err := Wrap(func() { i = 1 })()
	require.NoError(t, err)
	require.Equal(t, 1, i)
}

func TestWrapf(t *testing.T) {
	err := Wrapf(func() error { return errors.New("err") }, "wrapped : %w")()
	require.Error(t, err)
	require.Equal(t, "wrapped : err", err.Error())
}
