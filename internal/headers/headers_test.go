package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderParse(t *testing.T) {
	// Test: Valid double header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nFooFoo: barbar\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("HOST"))
	assert.Equal(t, "barbar", headers.Get("FooFoo"))
	assert.Equal(t, 39, n)
	assert.False(t, done)

	// Test: Valid double header
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nFooFoo: barbar\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("HOST"))
	assert.Equal(t, "barbar", headers.Get("FooFoo"))
	assert.Equal(t, 39, n)
	assert.True(t, done)

	// Test: Invalid char in field key
	headers = NewHeaders()
	data = []byte("h©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.ErrorIs(t, err, ErrInvalidFieldKey)
	require.NotNil(t, headers)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
