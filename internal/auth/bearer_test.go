package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBearerToken(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("Authorization", "Bearer abc.def.ghi")

		token, err := GetBearerToken(headers)

		assert.NoError(t, err)
		assert.Equal(t, "abc.def.ghi", token)
	})

	t.Run("missing header", func(t *testing.T) {
		headers := http.Header{}

		token, err := GetBearerToken(headers)

		assert.Error(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("wrong scheme", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("Authorization", "Basic abc123")

		token, err := GetBearerToken(headers)

		assert.Error(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("empty token", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("Authorization", "Bearer ")

		token, err := GetBearerToken(headers)

		assert.Error(t, err)
		assert.Equal(t, "", token)
	})
}
