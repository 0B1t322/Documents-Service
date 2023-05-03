package transform

import (
	. "github.com/0B1t322/Documents-Service/sessions/internal/core/operations/text"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnit_Transform(t *testing.T) {
	t.Run(
		"Insert Delete Text", func(t *testing.T) {
			start := "abcdef"

			t.Run(
				"Insert before delete", func(t *testing.T) {
					op1 := NewInsertText(1, "12")
					op2 := NewDeleteText(2, "cd")

					expected := "a12bef"
					op2Apply := op2.Apply(start)

					require.Equal(t, expected, InsertDelete(op1, op2).(InsertText).Apply(op2Apply))

				},
			)

			t.Run(
				"Insert after delete", func(t *testing.T) {
					op1 := NewInsertText(1, "12")
					op2 := NewDeleteText(0, "ab")

					expected := "12cdef"
					applyed := op2.Apply(start)
					require.Equal(t, expected, InsertDelete(op1, op2).(InsertText).Apply(applyed))
				},
			)
		},
	)

	t.Run(
		"Insert Insert Text", func(t *testing.T) {
			start := "abcdef"

			t.Run(
				"Insert before insert", func(t *testing.T) {
					op1 := NewInsertText(0, "1")
					op2 := NewInsertText(1, "2")

					expected := "1a2bcdef"

					applyed := op2.Apply(start)

					require.Equal(t, expected, InsertInsert(op1, op2).(InsertText).Apply(applyed))
				},
			)

			t.Run(
				"Insert before insert", func(t *testing.T) {
					op1 := NewInsertText(1, "1")
					op2 := NewInsertText(0, "2")

					expected := "2a1bcdef"

					applyed := op2.Apply(start)

					require.Equal(t, expected, InsertInsert(op1, op2).(InsertText).Apply(applyed))
				},
			)
		},
	)

	t.Run(
		"Delete Insert", func(t *testing.T) {
			start := "abcdef"

			t.Run(
				"Delete before insert", func(t *testing.T) {
					op1 := NewDeleteText(0, "a")
					op2 := NewInsertText(1, "1")

					require.Equal(t, "1bcdef", DeleteInsert(op1, op2).(DeleteText).Apply(op2.Apply(start)))
				},
			)

			t.Run(
				"Delete after insert", func(t *testing.T) {
					op1 := NewDeleteText(2, "c")
					op2 := NewInsertText(1, "1")

					require.Equal(t, "a1bdef", DeleteInsert(op1, op2).(DeleteText).Apply(op2.Apply(start)))
				},
			)
		},
	)

	t.Run(
		"Delete Delete Text", func(t *testing.T) {
			start := "abcdef"

			t.Run(
				"Delete before delete", func(t *testing.T) {
					op1 := NewDeleteText(0, "ab")
					op2 := NewDeleteText(3, "de")

					expected := "cf"

					require.Equal(t, expected, DeleteDelete(op1, op2).(DeleteText).Apply(op2.Apply(start)))
				},
			)

			t.Run(
				"Delete after delete", func(t *testing.T) {
					op1 := NewDeleteText(3, "de")
					op2 := NewDeleteText(0, "ab")

					expected := "cf"

					require.Equal(t, expected, DeleteDelete(op1, op1).(DeleteText).Apply(op2.Apply(start)))
				},
			)
		},
	)
}
