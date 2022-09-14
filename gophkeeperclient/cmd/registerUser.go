package cmd

import (
	"context"
	"fmt"
	"github.com/EestiChameleon/gophkeeper/gophkeeperclient/grpcclient"
	clstor "github.com/EestiChameleon/gophkeeper/gophkeeperclient/storage"
	pb "github.com/EestiChameleon/gophkeeper/proto"
	"github.com/spf13/cobra"
	"log"
	"os/user"
	"time"
)

// registerUserCmd represents the registerUser command
var registerUserCmd = &cobra.Command{
	Use:   "registerUser",
	Short: "Register new user in the service.",
	Long: `
This command register a new user.
Usage: gophkeeperclient registerUser --login=<login> --password=<password>.`,
	Run: func(cmd *cobra.Command, args []string) {
		// get current user from os/user. Like this we can locally identify if the user changed.
		u, err := user.Current()
		if err != nil {
			log.Fatalln(err)
			return
		}

		// request with 3s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		c, err := grpcclient.DialUp()
		if err != nil {
			log.Fatalln(err)
			return
		}

		// send request to server
		response, err := c.RegisterUser(ctx, &registerUser)
		if err != nil {
			log.Println(`[ERROR]:`, err)
			fmt.Println("request failed. please try again.")
			return
		}

		// save local pair - localUserName -> JWT
		clstor.Users[u.Username] = response.GetJwt()
		// init for the new user local storage
		clstor.Local[u.Username] = clstor.MakeVault()
		// return response
		fmt.Println(response.GetStatus())
	},
}

var (
	registerUser pb.RegisterUserRequest
)

func init() {
	rootCmd.AddCommand(registerUserCmd)
	registerUserCmd.Flags().StringVarP(&registerUser.ServiceLogin, "login", "l", "", "New user login value.")
	registerUserCmd.Flags().StringVarP(&registerUser.ServicePass, "password", "p", "", "New user password value.")
	registerUserCmd.MarkFlagRequired("login")
	registerUserCmd.MarkFlagRequired("password")
}
