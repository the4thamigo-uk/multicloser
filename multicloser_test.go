package multicloser_test

import (
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/the4thamigo-uk/multicloser"
	"testing"
)

func TestEmpty(t *testing.T) {
	m := multicloser.New()
	err := m.Close()
	require.NoError(t, err)
}

func TestNil(t *testing.T) {
	m := multicloser.New()
	require.Panics(t, func() { m.Defer(nil) })
}

func TestClose(t *testing.T) {
	m := multicloser.New()

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
	m := multicloser.New()

	err1 := errors.New("1")
	err2 := errors.New("2")

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
	require.Error(t, err)
	require.ErrorIs(t, err, err1)
	require.ErrorIs(t, err, err2)
}

func TestClosePanic(t *testing.T) {
	m := multicloser.New()

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
	err := multicloser.Wrap(func() { i = 1 })()
	require.NoError(t, err)
	require.Equal(t, 1, i)
}

func TestWrapf(t *testing.T) {
	err := multicloser.Wrapf(func() error { return errors.New("err") }, "wrapped : %w")()
	require.Error(t, err)
	require.Equal(t, "wrapped : err", err.Error())
}
