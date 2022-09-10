/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getTextCmd represents the getText command
var getTextCmd = &cobra.Command{
	Use:   "getText",
	Short: "Get a text data by title",
	Long: `
This command returns to the authenticated user the text data requested by title.
Usage: gophkeeperclient getText --title=<title>.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getText called")
	},
}

func init() {
	rootCmd.AddCommand(getTextCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getTextCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getTextCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
