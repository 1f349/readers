package readers

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestNewCommentReader(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		r := NewCommentReader(strings.NewReader("Hello world! # this is a comment"), []string{"#"})
		all, err := io.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, "Hello world! ", string(all))
		all, err = io.ReadAll(r.Comment())
		assert.NoError(t, err)
		assert.Equal(t, "# this is a comment", string(all))
	})
	t.Run("overflow start", func(t *testing.T) {
		r := NewCommentReader(strings.NewReader(strings.Repeat("Hello world!", 2048)+" # this is a comment"), []string{"#"})
		all, err := io.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, strings.Repeat("Hello world!", 2048)+" ", string(all))
		all, err = io.ReadAll(r.Comment())
		assert.NoError(t, err)
		assert.Equal(t, "# this is a comment", string(all))
	})
	t.Run("overflow comment", func(t *testing.T) {
		r := NewCommentReader(strings.NewReader(strings.Repeat("Hello world!", 2048)+" # "+strings.Repeat("this is a comment", 2048)), []string{"#"})
		all, err := io.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, strings.Repeat("Hello world!", 2048)+" ", string(all))
		all, err = io.ReadAll(r.Comment())
		assert.NoError(t, err)
		assert.Equal(t, "# "+strings.Repeat("this is a comment", 2048), string(all))
	})
}

func FuzzCommentReader(f *testing.F) {
	f.Fuzz(func(t *testing.T, a string) {
		r := NewCommentReader(strings.NewReader(a), []string{"#"})
		b1, err := io.ReadAll(r)
		assert.NoError(t, err)
		b2, err := io.ReadAll(r.Comment())
		assert.NoError(t, err)
		assert.Equal(t, string(b1)+string(b2), a)
	})
}
