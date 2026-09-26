package cli

import (
	"strconv"
	"strings"
)

// formatTemp renders a converted temperature for display.
//
// The value is printed in fixed-point notation rounded to six decimal
// places, trailing zeros are removed, and negative zero is normalized to
// zero, so 98.60000000000001 prints as 98.6 and 1.800032e+06 prints as
// 1800032.
func formatTemp(v float64) string {
	s := strconv.FormatFloat(v, 'f', 6, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimSuffix(s, ".")
	if s == "-0" {
		return "0"
	}
	return s
}
