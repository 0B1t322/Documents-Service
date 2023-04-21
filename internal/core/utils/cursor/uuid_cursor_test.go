package cursor_test

import (
	"github.com/0B1t322/Documents-Service/internal/core/utils/cursor"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFunc_UUIDCursor(t *testing.T) {
	t.Run(
		"Decode_Encode", func(t *testing.T) {
			id := uuid.New()

			cur := cursor.UUIDToCursor(id)

			decoded, err := cursor.CursorToUUID(cur)
			require.NoError(t, err)

			require.Equal(t, id, decoded)
		},
	)

	t.Run(
		"Decode_Encode_Nil", func(t *testing.T) {
			id := uuid.Nil

			cur := cursor.UUIDToCursor(id)

			decoded, err := cursor.CursorToUUID(cur)
			require.NoError(t, err)

			require.Equal(t, id, decoded)
		},
	)
}
