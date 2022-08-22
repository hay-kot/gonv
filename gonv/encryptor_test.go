package gonv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EncryptFile(t *testing.T) {
	path := envFactory(t, map[string]string{
		"FOO": "bar",
		"BAR": "baz",
	})

	err := EncryptFile(path, "pw")
	assert.NoError(t, err)

	contents, err := os.ReadFile(path)
	assert.NoError(t, err)

	assert.NotEqual(t, string(contents), "FOO=bar\nBAR=baz\n")
}

func Test_DecryptFile(t *testing.T) {
	path := envFactory(t, map[string]string{
		"FOO": "bar",
		"BAR": "baz",
	})

	err := EncryptFile(path, "pw")
	assert.NoError(t, err)

	err = DecryptFile(path, "pw")
	assert.NoError(t, err)

	contents, err := os.ReadFile(path)
	assert.NoError(t, err)

	assert.Equal(t, string(contents), "FOO=bar\nBAR=baz\n")
}
