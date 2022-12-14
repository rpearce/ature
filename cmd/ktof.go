package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var ktofCmd = &cobra.Command{
	Use:   "ktof",
	Short: "Convert Kelvin to Fahrenheit",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("no temperature was provided")
		}

		if value, err := strconv.ParseFloat(strings.TrimSpace(args[0]), 64); err == nil {
			fmt.Println((value-273.15)*(9.0/5.0) + 32)
		} else {
			return fmt.Errorf("temperature argument could not be parsed")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(ktofCmd)
}
