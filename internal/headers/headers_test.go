package headers

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestHeaders(t *testing.T) {

    // Test: Valid single header
    headers := NewHeaders()
    data := []byte("Host: localhost:42069\r\n\r\n")
    n, done, err := headers.Parse(data)
    require.NoError(t, err)
    require.NotNil(t, headers)
    assert.Equal(t, "localhost:42069", headers["host"])
    assert.Equal(t, 23, n)
    assert.False(t, done)

    // Test: Valid single header with extra spacing
    headers = NewHeaders()
    data = []byte("    Host:  localhost:42069  \r\n\r\n")
    n, done, err = headers.Parse(data)
    require.NoError(t, err)
    require.NotNil(t, headers)
    assert.Equal(t, "localhost:42069", headers["host"])
    assert.Equal(t, 30, n)
    assert.False(t, done)

    // Test: Single header with invalid spacing
    headers = NewHeaders()
    data = []byte("Host : localhost:42069\r\n\r\n")
    n, done, err = headers.Parse(data)
    require.Error(t, err)
    assert.Equal(t, 0, n)
    assert.False(t, done)

    // Test: Valid two headers with existing headers
    headers = map[string]string {"Host": "localhost:42069"}
    data = []byte("Content-Type: text/json\r\nAccept: */*\r\n\r\n")
    n, done, err = headers.Parse(data)
    require.NoError(t, err)
    require.NotNil(t, headers)
    assert.Equal(t, "localhost:42069", headers["Host"])
    assert.Equal(t, "text/json", headers["content-type"])
    assert.Equal(t, 25, n)
    assert.False(t, done)

    // Test: Single header with invalid character
    headers = NewHeaders()
    data = []byte("Host@: localhost:42069\r\n\r\n")
    n, done, err = headers.Parse(data)
    require.Error(t, err)
    assert.Equal(t, 0, n)
    assert.False(t, done)

    // Test: Valid done
    headers = NewHeaders()
    data = []byte("\r\n other data")
    n, done, err = headers.Parse(data)
    require.NoError(t, err)
    require.NotNil(t, headers)
    assert.Empty(t, headers)
    assert.Equal(t, 2, n)
    assert.True(t, done)

    // Test: Valid two headers with existing headers
    headers = map[string]string {"multiples": "itemOne"}
    data = []byte("multiples: itemTwo\r\n")
    n, done, err = headers.Parse(data)
    require.NoError(t, err)
    require.NotNil(t, headers)
    assert.Equal(t, "itemOne, itemTwo", headers["multiples"])
    assert.Equal(t, 20, n)
    assert.False(t, done)

}
