package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
	
	// Test: Valid single header with extra whitespace
	headers = NewHeaders()
	data = []byte("       Host:     localhost:42069                  \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 52, n)
	assert.False(t, done)

	// Test: Valid 2 headers with existing headers
	headers = NewHeaders()
	data = []byte("Host:   localhost:42069\r\n    Something: hey \r\n  Content-Type:   application/json\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 25, n)
	assert.False(t, done)

	num := n
	n, done, err = headers.Parse(data[n:])
	n += num
	require.NoError(t, err)
	assert.Equal(t, "hey", headers.Get("Something"))
	assert.Equal(t, 46, n)
	assert.False(t, done)
	
	num = n
	n, done, err = headers.Parse(data[n:])
	n += num
	require.NoError(t, err)
	assert.Equal(t, "application/json", headers.Get("Content-Type"))
	assert.Equal(t, 82, n)
	assert.False(t, done)

	num = n
	n, done, err = headers.Parse(data[n:])
	n += num
	require.NoError(t, err)
	assert.Equal(t, 82, n)
	assert.True(t, done)

	// Test: Valid done
	headers = NewHeaders()
	data = []byte("\r\n    ")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.True(t, done)

	// Test: Invalid characters
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, n, 0)
	assert.False(t, done)

	// Test: Multiple of the same header
	headers = NewHeaders()
	data = []byte("  Set-Person: bob  \r\n  Set-Person: fred \r\n Set-Person: joe\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.False(t, done)

	num = n
	n, done, err = headers.Parse(data[n:])
	n += num
	require.NoError(t, err)
	assert.False(t, done)

	num = n
	n, done, err = headers.Parse(data[n:])
	n += num
	require.NoError(t, err)
	assert.False(t, done)
	assert.Equal(t, "bob, fred, joe", headers.Get("Set-Person"))
}
