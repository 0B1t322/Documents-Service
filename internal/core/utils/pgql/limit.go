package pgql

// Limit return limit arg
// if passed limit > 0, return this limit
// else return nil
func Limit(limit uint) any {
	if limit > 0 {
		return limit
	}

	return nil
}
