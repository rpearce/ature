// Package cli builds the ature command tree.
package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewRootCmd returns a fresh root command with every conversion subcommand
// attached. version is what the --version flag reports.
//
// Each call builds an independent tree, so tests can run in parallel with
// their own output writers.
func NewRootCmd(version string) *cobra.Command {
	root := &cobra.Command{
		Use:     "ature",
		Short:   "Convert temperature values",
		Version: version,

		// main reports errors exactly once; Cobra must not also print the
		// error or the usage block.
		SilenceUsage:  true,
		SilenceErrors: true,

		// Keep shell completion working but out of the help listing.
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},

		// Plain "ature" prints help. Arguments reach this function only when
		// they follow "--", because Cobra stops looking for a command there
		// and rejects any other stray argument itself; without this check
		// "ature -- -10" would print help and exit 0.
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())
			}

			return cmd.Help()
		},
	}

	// Subcommands without their own handler inherit this one, so it covers
	// "ature -10" and the help and completion commands.
	root.SetFlagErrorFunc(negativeNumberHint("ature " + conversions[0].name))

	for _, c := range conversions {
		root.AddCommand(newConvertCmd(c))
	}

	return root
}
