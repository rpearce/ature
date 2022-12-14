/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// ftocCmd represents the ftoc command
var ftocCmd = &cobra.Command{
	Use:   "ftoc",
	Short: "Convert Fahrenheit to Celsius",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("no temperature was provided")
		}

		if value, err := strconv.ParseFloat(strings.TrimSpace(args[0]), 64); err == nil {
			fmt.Println((value - 32) * (5.0 / 9.0))
		} else {
			return fmt.Errorf("temperature argument could not be parsed")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(ftocCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ftocCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ftocCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
