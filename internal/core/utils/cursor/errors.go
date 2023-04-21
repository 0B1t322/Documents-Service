package cursor

import (
	"fmt"
)

type CursorDecodeError struct {
	ShortReason string
}

func (c CursorDecodeError) Error() string {
	return fmt.Sprintf("%s: %s", c.description(), c.ShortReason)
}

func (CursorDecodeError) description() string {
	return "Can't decode cursor"
}
