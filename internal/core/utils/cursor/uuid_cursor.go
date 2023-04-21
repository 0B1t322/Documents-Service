package cursor

import (
	"encoding/base64"
	"github.com/google/uuid"
)

func UUIDToCursor(id uuid.UUID) string {
	// Skip error because it's always nil
	bytes, _ := id.MarshalBinary()
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func CursorToUUID(cursor string) (uuid.UUID, error) {
	bytes, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return uuid.Nil, err
	}

	cur, err := uuid.FromBytes(bytes)
	if err != nil {
		return uuid.Nil, &CursorDecodeError{
			ShortReason: "cursor is not uuid",
		}
	}

	return cur, nil
}
