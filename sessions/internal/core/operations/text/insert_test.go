package text

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnit_Insert(t *testing.T) {
	t.Run(
		"Apply", func(t *testing.T) {
			to := "Hlo world"
			op := NewInsertText(1, "el")

			require.Equal(t, "Hello world", op.Apply(to))
		},
	)
}
