package request_test

import (
	"testing"

	"github.com/joaopandolfi/blackwhale/remotes/request"
	"github.com/stretchr/testify/assert"
)

func TestInjectHeaderValues(t *testing.T) {

	header := map[string]string{}

	request.InjectAuthorization("ABC", header)
	request.InjectAcceptGzip(header)

	expected := map[string]string{
		"Authorization":   "Bearer ABC",
		"Accept-Encoding": "gzip",
	}

	assert.Equal(t, header, expected)

	request.InjectSendGzip(header)
	expected["Content-Encoding"] = "gzip"

	assert.Equal(t, header, expected)
}
