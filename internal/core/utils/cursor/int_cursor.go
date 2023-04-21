package cursor

import (
	"fmt"
	"strconv"
)

func IntToCursor(cursor int) string {
	return fmt.Sprint(cursor)
}

func CursorToInt(cursor string) (int, error) {
	c, err := strconv.Atoi(cursor)
	if err != nil {
		return 0, &CursorDecodeError{
			ShortReason: "cursor is not a int",
		}
	}
	return c, nil
}
