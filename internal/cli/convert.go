package cli

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rpearce/ature/internal/temp"
)

// conversion describes one subcommand that converts between two scales.
type conversion struct {
	name     string
	from, to temp.Unit
}

// conversions is the complete list of subcommands.
var conversions = []conversion{
	{"ctof", temp.Celsius, temp.Fahrenheit},
	{"ctok", temp.Celsius, temp.Kelvin},
	{"ftoc", temp.Fahrenheit, temp.Celsius},
	{"ftok", temp.Fahrenheit, temp.Kelvin},
	{"ktoc", temp.Kelvin, temp.Celsius},
	{"ktof", temp.Kelvin, temp.Fahrenheit},
}

// errNotFinite is returned for inputs such as "inf" or "nan" that parse as
// floats but are not temperatures.
var errNotFinite = errors.New("temperature must be a finite number")

// negativeNumberFlag matches pflag's "unknown shorthand flag: '1' in -10"
// error when the rejected cluster starts with a digit or a dot, which is
// how a negative number given without "--" surfaces.
var negativeNumberFlag = regexp.MustCompile(`in -[0-9.]`)

// newConvertCmd builds the subcommand for one conversion.
func newConvertCmd(c conversion) *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.name + " <value>",
		Short: "Convert " + c.from.String() + " to " + c.to.String(),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			value, err := parseTemp(args[0])
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(cmd.OutOrStdout(), formatTemp(temp.Convert(value, c.from, c.to)))

			return err
		},
	}

	cmd.SetFlagErrorFunc(negativeNumberHint("ature " + c.name))

	return cmd
}

// negativeNumberHint returns a Cobra flag-error function. "ature ctof -10"
// fails flag parsing before the value is ever seen, so the function points at
// the POSIX "--" convention, showing example as the command to type. Other
// flag errors, such as a typo like --foo, pass through untouched.
func negativeNumberHint(example string) func(*cobra.Command, error) error {
	return func(_ *cobra.Command, err error) error {
		if !negativeNumberFlag.MatchString(err.Error()) {
			return err
		}

		return fmt.Errorf("%w (negative values must follow \"--\", for example: %s -- -10)", err, example)
	}
}

// parseTemp parses s as a finite temperature value. Surrounding whitespace
// is ignored so values piped in from other tools still work.
func parseTemp(s string) (float64, error) {
	s = strings.TrimSpace(s)

	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// NumError's message repeats the function name and input; keep only
		// the reason ("invalid syntax" or "value out of range").
		if numErr, ok := errors.AsType[*strconv.NumError](err); ok {
			err = numErr.Err
		}

		return 0, fmt.Errorf("parse temperature %q: %w", s, err)
	}

	if math.IsInf(v, 0) || math.IsNaN(v) {
		return 0, errNotFinite
	}

	return v, nil
}
