package temp_test

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/rpearce/ature/internal/temp"
)

// tolerance absorbs floating-point noise; 37 C is 98.60000000000001 F in
// float64, which is correct but not equal to 98.6.
const tolerance = 1e-9

func TestConvert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    float64
		from, to temp.Unit
		want     float64
	}{
		{"freezing C to F", 0, temp.Celsius, temp.Fahrenheit, 32},
		{"freezing C to K", 0, temp.Celsius, temp.Kelvin, 273.15},
		{"freezing F to C", 32, temp.Fahrenheit, temp.Celsius, 0},
		{"freezing F to K", 32, temp.Fahrenheit, temp.Kelvin, 273.15},
		{"freezing K to C", 273.15, temp.Kelvin, temp.Celsius, 0},
		{"freezing K to F", 273.15, temp.Kelvin, temp.Fahrenheit, 32},
		{"boiling C to F", 100, temp.Celsius, temp.Fahrenheit, 212},
		{"boiling C to K", 100, temp.Celsius, temp.Kelvin, 373.15},
		{"boiling F to K", 212, temp.Fahrenheit, temp.Kelvin, 373.15},
		{"minus forty C to F", -40, temp.Celsius, temp.Fahrenheit, -40},
		{"minus forty F to C", -40, temp.Fahrenheit, temp.Celsius, -40},
		{"absolute zero K to C", 0, temp.Kelvin, temp.Celsius, -273.15},
		{"absolute zero K to F", 0, temp.Kelvin, temp.Fahrenheit, -459.67},
		{"body temperature C to F", 37, temp.Celsius, temp.Fahrenheit, 98.6},
		{"identity C", 21.5, temp.Celsius, temp.Celsius, 21.5},
		{"identity F", 70.7, temp.Fahrenheit, temp.Fahrenheit, 70.7},
		{"identity K", 300, temp.Kelvin, temp.Kelvin, 300},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := temp.Convert(tt.value, tt.from, tt.to)
			if math.Abs(got-tt.want) > tolerance {
				t.Errorf("Convert(%v, %v, %v) = %v, want %v", tt.value, tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestConvertRoundTrip(t *testing.T) {
	t.Parallel()

	units := []temp.Unit{temp.Celsius, temp.Fahrenheit, temp.Kelvin}
	values := []float64{-300, -40, 0, 0.1, 36.6, 1000}

	for _, from := range units {
		for _, to := range units {
			for _, v := range values {
				got := temp.Convert(temp.Convert(v, from, to), to, from)
				if math.Abs(got-v) > tolerance {
					t.Errorf("%v -> %v -> %v round trip of %v = %v", from, to, from, v, got)
				}
			}
		}
	}
}

func TestUnknownUnitPanics(t *testing.T) {
	t.Parallel()

	unknown := temp.Unit(99)

	tests := []struct {
		name string
		call func()
	}{
		{"Convert from unknown", func() { temp.Convert(0, unknown, temp.Celsius) }},
		{"Convert to unknown", func() { temp.Convert(0, temp.Celsius, unknown) }},
		{"String of unknown", func() { _ = unknown.String() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			defer func() {
				r := recover()
				if r == nil {
					t.Fatalf("%s did not panic", tt.name)
				}
				if !strings.Contains(fmt.Sprint(r), "99") {
					t.Errorf("%s panicked with %v, want a message naming unit 99", tt.name, r)
				}
			}()

			tt.call()
		})
	}
}

func TestUnitString(t *testing.T) {
	t.Parallel()

	tests := map[temp.Unit]string{
		temp.Celsius:    "Celsius",
		temp.Fahrenheit: "Fahrenheit",
		temp.Kelvin:     "Kelvin",
	}

	for u, want := range tests {
		if got := u.String(); got != want {
			t.Errorf("Unit(%d).String() = %q, want %q", u, got, want)
		}
	}
}
