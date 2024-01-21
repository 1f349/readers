package readers

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go/scanner"
	"io"
	"strings"
	"testing"
)

func TestNewLineReader(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		r := NewLineReader(strings.NewReader("Hello world!\nHello new line!\n"))
		for r.Next() {
			a, err := io.ReadAll(r)
			fmt.Println(len(a), string(a))
			assert.NoError(t, err)
		}
		assert.NoError(t, r.Err())
	})
	t.Run("really long", func(t *testing.T) {
		r := NewLineReader(strings.NewReader(strings.Repeat("Hello world! ", 1024) + "\na\n"))
		for r.Next() {
			a, err := io.ReadAll(r)
			fmt.Println(len(a), string(a))
			assert.NoError(t, err)
		}
		assert.NoError(t, r.Err())
	})
	t.Run("against scanner", func(t *testing.T) {
		r := NewLineReader(strings.NewReader(strings.Repeat("Hello world! ", 409600) + "\na\n"))
		for r.Next() {
			a, err := io.ReadAll(r)
			fmt.Println(len(a), string(a))
			assert.NoError(t, err)
		}
		assert.NoError(t, r.Err())

		sc := bufio.NewScanner(strings.NewReader(strings.Repeat("Hello world! ", 409600) + "\na\n"))
		assert.False(t, sc.Scan())
		assert.Error(t, sc.Err(), scanner.Error{})
	})
}
