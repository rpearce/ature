package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rpearce/ature/internal/cli"
)

// run executes ature with args on a fresh command tree and returns what was
// written to stdout and stderr along with the error from Execute.
func run(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()

	if args == nil {
		// A nil slice makes Cobra fall back to os.Args, which inside go test
		// are the test binary's own flags.
		args = []string{}
	}

	var out, errOut bytes.Buffer
	root := cli.NewRootCmd("1.2.3")
	root.SetOut(&out)
	root.SetErr(&errOut)
	root.SetArgs(args)

	err = root.Execute()

	return out.String(), errOut.String(), err
}

func TestConvertCommands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want string
	}{
		{"ctof", []string{"ctof", "16"}, "60.8\n"},
		{"ctok", []string{"ctok", "0"}, "273.15\n"},
		{"ftoc", []string{"ftoc", "212"}, "100\n"},
		{"ftok", []string{"ftok", "32"}, "273.15\n"},
		{"ktoc", []string{"ktoc", "0"}, "-273.15\n"},
		{"ktof", []string{"ktof", "273.15"}, "32\n"},
		{"float noise is rounded", []string{"ctof", "37"}, "98.6\n"},
		{"negative after double dash", []string{"ctof", "--", "-40"}, "-40\n"},
		{"repeating decimal", []string{"ftoc", "--", "-10"}, "-23.333333\n"},
		{"large value has no exponent", []string{"ctof", "1000000"}, "1800032\n"},
		{"result near zero prints as zero", []string{"ftoc", "31.9999999"}, "0\n"},
		{"surrounding whitespace is ignored", []string{"ctof", " 16 "}, "60.8\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stdout, stderr, err := run(t, tt.args...)
			if err != nil {
				t.Fatalf("ature %v: unexpected error: %v", tt.args, err)
			}
			if stdout != tt.want {
				t.Errorf("ature %v: stdout = %q, want %q", tt.args, stdout, tt.want)
			}
			if stderr != "" {
				t.Errorf("ature %v: stderr = %q, want nothing", tt.args, stderr)
			}
		})
	}
}

func TestConvertCommandErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		wantErr string // substring of the returned error
	}{
		{"no value", []string{"ctof"}, "accepts 1 arg(s), received 0"},
		{"two values", []string{"ctof", "16", "99"}, "accepts 1 arg(s), received 2"},
		{"not a number", []string{"ctof", "abc"}, `parse temperature "abc": invalid syntax`},
		{"out of range", []string{"ctof", "1e400"}, `parse temperature "1e400": value out of range`},
		{"nan", []string{"ctof", "nan"}, "finite"},
		{"inf", []string{"ctof", "inf"}, "finite"},
		{"negative inf", []string{"ktoc", "--", "-inf"}, "finite"},
		{"negative without double dash", []string{"ctof", "-10"}, `negative values must follow "--", for example: ature ctof -- -10`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stdout, stderr, err := run(t, tt.args...)
			if err == nil {
				t.Fatalf("ature %v: expected an error, got none", tt.args)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ature %v: error %q does not contain %q", tt.args, err, tt.wantErr)
			}
			if stdout != "" {
				t.Errorf("ature %v: stdout = %q, want nothing", tt.args, stdout)
			}
			if stderr != "" {
				t.Errorf("ature %v: stderr = %q, want nothing (main reports errors)", tt.args, stderr)
			}
		})
	}
}

func TestRootErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		wantErr string // substring of the returned error
	}{
		{"negative without a command", []string{"-10"}, `negative values must follow "--", for example: ature ctof -- -10`},
		{"value after double dash without a command", []string{"--", "-10"}, `unknown command "-10" for "ature"`},
		// Regression guard: Cobra's own unknown-command check, with its
		// suggestions, must keep working for ordinary typos.
		{"unknown command keeps suggestions", []string{"ctoff", "16"}, "Did you mean this?"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stdout, stderr, err := run(t, tt.args...)
			if err == nil {
				t.Fatalf("ature %v: expected an error, got none (stdout %q)", tt.args, stdout)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ature %v: error %q does not contain %q", tt.args, err, tt.wantErr)
			}
			if stdout != "" {
				t.Errorf("ature %v: stdout = %q, want nothing", tt.args, stdout)
			}
			if stderr != "" {
				t.Errorf("ature %v: stderr = %q, want nothing (main reports errors)", tt.args, stderr)
			}
		})
	}
}

func TestHugeValueStaysFixedPoint(t *testing.T) {
	t.Parallel()

	stdout, _, err := run(t, "ctok", "1e300")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.ContainsAny(stdout, "eE+") || strings.Contains(stdout, "Inf") {
		t.Errorf("ature ctok 1e300 printed %q, want a plain fixed-point number", stdout)
	}
	if !strings.HasSuffix(stdout, "\n") {
		t.Errorf("ature ctok 1e300 printed %q, want a trailing newline", stdout)
	}
}

func TestVersionFlag(t *testing.T) {
	t.Parallel()

	for _, flag := range []string{"--version", "-v"} {
		t.Run(flag, func(t *testing.T) {
			t.Parallel()

			stdout, stderr, err := run(t, flag)
			if err != nil {
				t.Fatalf("ature %s: unexpected error: %v", flag, err)
			}
			if stdout != "ature version 1.2.3\n" {
				t.Errorf("ature %s: stdout = %q, want %q", flag, stdout, "ature version 1.2.3\n")
			}
			if stderr != "" {
				t.Errorf("ature %s: stderr = %q, want nothing", flag, stderr)
			}
		})
	}
}

func TestRootWithoutArgsPrintsHelp(t *testing.T) {
	t.Parallel()

	stdout, _, err := run(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stdout, "Usage:") {
		t.Errorf("ature with no args printed %q, want the help text", stdout)
	}
}

func TestHelpListsCommandsAndHidesCompletion(t *testing.T) {
	t.Parallel()

	stdout, _, err := run(t, "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, want := range []string{
		"ctof        Convert Celsius to Fahrenheit",
		"ctok        Convert Celsius to Kelvin",
		"ftoc        Convert Fahrenheit to Celsius",
		"ftok        Convert Fahrenheit to Kelvin",
		"ktoc        Convert Kelvin to Celsius",
		"ktof        Convert Kelvin to Fahrenheit",
	} {
		if !strings.Contains(stdout, want) {
			t.Errorf("help output is missing %q:\n%s", want, stdout)
		}
	}

	if strings.Contains(stdout, "completion") {
		t.Errorf("help output should hide the completion command:\n%s", stdout)
	}
}

func TestSubcommandHelpShowsValuePlaceholder(t *testing.T) {
	t.Parallel()

	stdout, _, err := run(t, "ctof", "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stdout, "ature ctof <value>") {
		t.Errorf("ctof help printed %q, want it to contain %q", stdout, "ature ctof <value>")
	}
}

func TestCompletionStillWorksWhenHidden(t *testing.T) {
	t.Parallel()

	stdout, _, err := run(t, "completion", "zsh")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stdout, "compdef") {
		t.Errorf("ature completion zsh printed %q, want a zsh completion script", stdout)
	}
}

func TestFlagErrorHint(t *testing.T) {
	t.Parallel()

	const hint = `negative values must follow "--"`

	tests := []struct {
		name     string
		args     []string
		wantHint bool
	}{
		{"negative integer", []string{"ctof", "-10"}, true},
		{"negative decimal", []string{"ctof", "-.5"}, true},
		{"unknown long flag", []string{"ctof", "--foo"}, false},
		{"unknown short flag", []string{"ctof", "-x"}, false},
		{"negative without a command", []string{"-10"}, true},
		{"unknown short flag without a command", []string{"-x"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, _, err := run(t, tt.args...)
			if err == nil {
				t.Fatalf("ature %v: expected an error, got none", tt.args)
			}
			if got := strings.Contains(err.Error(), hint); got != tt.wantHint {
				t.Errorf("ature %v: error %q; hint present = %v, want %v", tt.args, err, got, tt.wantHint)
			}
		})
	}
}
