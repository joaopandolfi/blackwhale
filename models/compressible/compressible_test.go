package compressible

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompression(t *testing.T) {

	rawData := "veeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeery long striiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiing"

	data := []byte(rawData)

	result, err := ZipBytes(data)
	assert.Nil(t, err, "zipping bytes")
	assert.True(t, len(data) > len(result), "compressed result must be smaller than original")
}

type _test_compress struct {
	payload string
	Compressible
}

func (c *_test_compress) Reffs() []*string {
	return []*string{&c.payload}
}

func (c *_test_compress) Compress() error {
	return c.Compressible.Compress(c.Reffs())
}

func (c *_test_compress) Decompress() error {
	return c.Compressible.Decompress(c.Reffs())
}

func TestCompressible(t *testing.T) {
	var c _test_compress

	originalValue := "veeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeery long striiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiing"

	c.payload = originalValue

	assert.Nil(t, c.Compress(), c.payload)
	assert.NotEqual(t, c.payload, originalValue)
	assert.True(t, len(originalValue) > len(c.payload), "compressed result must be smaller than original")
	assert.Nil(t, c.Decompress(), c.payload)
	assert.Equal(t, c.payload, originalValue)
}
