package cmd

import (
	"context"
	"fmt"
	"github.com/EestiChameleon/gophkeeper/gophkeeperclient/grpcclient"
	clserv "github.com/EestiChameleon/gophkeeper/gophkeeperclient/service"
	clstor "github.com/EestiChameleon/gophkeeper/gophkeeperclient/storage"
	pb "github.com/EestiChameleon/gophkeeper/proto"
	"github.com/spf13/cobra"
	"log"
	"os/user"
	"time"
)

// loginUserCmd represents the loginUser command
var loginUserCmd = &cobra.Command{
	Use:   "loginUser",
	Short: "Login user to the service",
	Long: `
This command login user.
Usage: gophkeeperclient loginUser --login=<login> --password=<password>.`,
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

		// send data to server and receive JWT in case of success. then save it in Users
		response, err := c.LoginUser(ctx, &loginUser)
		if err != nil {
			log.Println(`[ERROR]:`, err)
			fmt.Println("request failed. please try again.")
			return
		}

		// save local username = JWT for server
		clstor.Users[u.Username] = response.GetJwt()
		// save local user data
		clstor.Local[u.Username] = clserv.VaultSyncConvert(response.AllData)

		fmt.Println(response.GetStatus())
	},
}

var (
	loginUser pb.LoginUserRequest
)

func init() {
	rootCmd.AddCommand(loginUserCmd)
	loginUserCmd.Flags().StringVarP(&loginUser.ServiceLogin, "login", "l", "", "New user login value.")
	loginUserCmd.Flags().StringVarP(&loginUser.ServicePass, "password", "p", "", "New user password value.")
	loginUserCmd.MarkFlagRequired("login")
	loginUserCmd.MarkFlagRequired("password")
}
