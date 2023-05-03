package text

import "testing"

func TestUnit_Delete(t *testing.T) {
	t.Run(
		"Apply", func(t *testing.T) {
			to := "Hellelo world"
			op := NewDeleteText(4, "el")

			t.Log(op.Apply(to))
		},
	)
}
