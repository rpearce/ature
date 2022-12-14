package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var ctofCmd = &cobra.Command{
	Use:   "ctof",
	Short: "Convert Celsius to Fahrenheit",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("no temperature was provided")
		}

		if value, err := strconv.ParseFloat(strings.TrimSpace(args[0]), 64); err == nil {
			fmt.Println(1.8*value + 32)
		} else {
			return fmt.Errorf("temperature argument could not be parsed")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(ctofCmd)
}
