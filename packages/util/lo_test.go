package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTake(t *testing.T) {
	type testStruct struct {
		i int
		s string
	}
	e1 := testStruct{i: 1, s: "one"}
	e2 := testStruct{i: 2, s: "two"}
	e3 := testStruct{i: 3, s: "three"}
	e4 := testStruct{i: 4, s: "four"}
	eo := testStruct{i: 1, s: "vienas"}
	all := []testStruct{e1, e2, e3, eo, e4}

	elem, list, result := Take(all, func(ts testStruct) bool { return ts.i == 1 })
	require.True(t, result)
	require.Equal(t, e1, elem)
	require.Equal(t, []testStruct{e2, e3, eo, e4}, list)

	elem, list, result = Take(all, func(ts testStruct) bool { return ts.i == 2 })
	require.True(t, result)
	require.Equal(t, e2, elem)
	require.Equal(t, []testStruct{e1, e3, eo, e4}, list)

	elem, list, result = Take(all, func(ts testStruct) bool { return ts.i == 4 })
	require.True(t, result)
	require.Equal(t, e4, elem)
	require.Equal(t, []testStruct{e1, e2, e3, eo}, list)

	_, list, result = Take(all, func(ts testStruct) bool { return ts.i == 10 })
	require.False(t, result)
	require.Equal(t, all, list)

	_, list, result = Take([]testStruct{}, func(ts testStruct) bool { return ts.i == 1 })
	require.False(t, result)
	require.Equal(t, []testStruct{}, list)

	elem, list, result = Take([]testStruct{e1}, func(ts testStruct) bool { return ts.i == 1 })
	require.True(t, result)
	require.Equal(t, e1, elem)
	require.Equal(t, []testStruct{}, list)

	_, list, result = Take([]testStruct{e1}, func(ts testStruct) bool { return ts.i == 2 })
	require.False(t, result)
	require.Equal(t, []testStruct{e1}, list)
}
