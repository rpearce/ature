package cli

import (
	"math"
	"testing"
)

func TestFormatTemp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   float64
		want string
	}{
		{"integer", 32, "32"},
		{"one decimal", 60.8, "60.8"},
		{"float noise is rounded away", 98.60000000000001, "98.6"},
		{"repeating decimal is cut at six places", -23.333333333333336, "-23.333333"},
		{"large value uses no exponent", 1800032, "1800032"},
		{"very large value stays fixed-point", 1e15, "1000000000000000"},
		{"negative", -40, "-40"},
		{"zero", 0, "0"},
		{"negative zero", math.Copysign(0, -1), "0"},
		{"tiny negative rounds to zero, not minus zero", -0.0000001, "0"},
		{"seventh decimal rounds up", 0.0000006, "0.000001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := formatTemp(tt.in); got != tt.want {
				t.Errorf("formatTemp(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
