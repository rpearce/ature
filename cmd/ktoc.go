package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var ktocCmd = &cobra.Command{
	Use:   "ktoc",
	Short: "Convert Kelvin to Celsius",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("no temperature was provided")
		}

		if value, err := strconv.ParseFloat(strings.TrimSpace(args[0]), 64); err == nil {
			fmt.Println(value - 273.15)
		} else {
			return fmt.Errorf("temperature argument could not be parsed")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(ktocCmd)
}
