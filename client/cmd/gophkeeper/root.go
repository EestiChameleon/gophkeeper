package gophkeeper

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "gophkeeper",
	Version: version,
	Short:   "gophkeeper - a simple CLI to store user data (login&pass pair, card data, text, etc)",
	Long: `gophkeeper is a super fancy CLI (kidding)
   
One can use gophkeeper to store or retrieve important data straight from the terminal`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
