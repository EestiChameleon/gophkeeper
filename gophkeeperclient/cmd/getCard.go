/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getCardCmd represents the getCard command
var getCardCmd = &cobra.Command{
	Use:   "getCard",
	Short: "Get a card data by title",
	Long: `
This command returns to the authenticated user the card data requested by title.
Usage: gophkeeperclient getCard --title=<title>.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getCard called")
	},
}

func init() {
	rootCmd.AddCommand(getCardCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCardCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCardCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
