/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gophkeeperclient",
	Short: "GophKeeper is a service to store and protect your important data",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

GophKeeper is a service, that gives you the possibilities to save you data and retrieve it from different devices. 
Service is synchronized between all you devices, where you are authenticated.
This application is a CLI tool to interact with the service.
Type -help to see more.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ready")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
