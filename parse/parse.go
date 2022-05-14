// Package parse implements methods to help with parsing data.

package parse

import "strconv"

// Int parses a string into an integer, with a default fallback for any empty/error cases.
func Int(s string, dfault int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return dfault
	}
	return i
}
