// Package temp converts temperatures between Celsius, Fahrenheit, and Kelvin.
package temp

import "strconv"

// Unit is a temperature scale.
type Unit uint8

// The supported temperature scales.
const (
	Celsius Unit = iota
	Fahrenheit
	Kelvin
)

// String returns the scale's name, for example "Celsius".
func (u Unit) String() string {
	switch u {
	case Celsius:
		return "Celsius"
	case Fahrenheit:
		return "Fahrenheit"
	case Kelvin:
		return "Kelvin"
	default:
		panic(unknownUnit(u))
	}
}

// Convert converts value from one scale to another.
//
// Every conversion routes through Celsius, so each scale has exactly one
// formula in each direction and round trips stay consistent. Convert panics
// if either unit is not one of the declared constants.
func Convert(value float64, from, to Unit) float64 {
	return fromCelsius(toCelsius(value, from), to)
}

func toCelsius(v float64, u Unit) float64 {
	switch u {
	case Celsius:
		return v
	case Fahrenheit:
		return (v - 32) * 5 / 9
	case Kelvin:
		return v - 273.15
	default:
		panic(unknownUnit(u))
	}
}

func fromCelsius(c float64, u Unit) float64 {
	switch u {
	case Celsius:
		return c
	case Fahrenheit:
		return c*9/5 + 32
	case Kelvin:
		return c + 273.15
	default:
		panic(unknownUnit(u))
	}
}

// unknownUnit builds the panic message for a Unit outside the declared
// constants. It names the value so the caller can see what was passed.
func unknownUnit(u Unit) string {
	return "temp: unknown unit " + strconv.Itoa(int(u))
}
