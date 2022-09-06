package gophkeeper

import (
	"fmt"

	"github.com/EestiChameleon/gophkeeper/client/pkg/gophkeeper"
	"github.com/spf13/cobra"
)

var reverseCmd = &cobra.Command{
	Use:     "reverse",
	Aliases: []string{"rev"},
	Short:   "Reverses a string",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res := gophkeeper.Reverse(args[0])
		fmt.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(reverseCmd)
}
