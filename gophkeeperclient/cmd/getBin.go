/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getBinCmd represents the getBin command
var getBinCmd = &cobra.Command{
	Use:   "getBin",
	Short: "Get a bin data by title",
	Long: `
This command returns to the authenticated user the binary data requested by title.
Usage: gophkeeperclient getBin --title=<title>.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getBin called")
	},
}

func init() {
	rootCmd.AddCommand(getBinCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getBinCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getBinCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
