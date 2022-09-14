package cmd

import (
	"context"
	"fmt"
	"github.com/EestiChameleon/gophkeeper/gophkeeperclient/grpcclient"
	clserv "github.com/EestiChameleon/gophkeeper/gophkeeperclient/service"
	clstor "github.com/EestiChameleon/gophkeeper/gophkeeperclient/storage"
	pb "github.com/EestiChameleon/gophkeeper/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
	"log"
	"os/user"
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
		//ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		//defer cancel()

		c, err := grpcclient.DialUp()
		if err != nil {
			log.Fatalln(err)
			return
		}

		// send data to server and receive JWT in case of success. then save it in Users.
		logResp, err := c.LoginUser(context.Background(), &loginUser)
		if err != nil {
			log.Println(`[ERROR]:`, err)
			fmt.Println("request failed. please try again.")
			return
		}

		fmt.Println("login status: ", logResp.GetStatus())

		// save JWT
		clstor.Users[u.Username] = logResp.Jwt
		// check for nil vault
		// update to latest data
		locV, ok := clstor.Local[u.Username]
		if !ok {
			// local storage not initiated
			locV = clstor.MakeVault()
			clstor.Local[u.Username] = locV
		}

		// after successful login - get JWT and send to server to synchronize data.
		ctxWTKN := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+logResp.Jwt)

		syncResp, err := c.SyncVault(ctxWTKN, &syncData)
		if err != nil {
			log.Println(`[ERROR]:`, err)
			fmt.Println("request failed. please try again.")
			return
		}

		fmt.Println("Get latest data from server: ", syncResp.GetStatus())

		fmt.Print("Synchronizing: ")
		updVault := clserv.CombineVault(clstor.Local[u.Username], clserv.VaultSyncConvert(syncResp))
		//save actual data
		clstor.Local[u.Username] = updVault

		fmt.Println(syncResp.GetStatus())
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
