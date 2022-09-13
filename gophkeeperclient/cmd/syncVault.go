/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"github.com/EestiChameleon/gophkeeper/gophkeeperclient/grpcclient"
	clserv "github.com/EestiChameleon/gophkeeper/gophkeeperclient/service"
	clstor "github.com/EestiChameleon/gophkeeper/gophkeeperclient/storage"
	pb "github.com/EestiChameleon/gophkeeper/proto"
	"google.golang.org/grpc/metadata"
	"log"
	"os/user"
	"time"

	"github.com/spf13/cobra"
)

// syncVaultCmd represents the syncVault command
var syncVaultCmd = &cobra.Command{
	Use:   "syncVault",
	Short: "Synchronize your local vault with the server database.",
	Long: `
This command provides latest data from the server database.
Then the database data, with version higher, that the local version, is saved to local storage.
During the saving of local data to the database, in case of version conflict(database version is higher/newer), you will be alerted by a warning.
Usage: gophkeeperclient syncVault`,
	Run: func(cmd *cobra.Command, args []string) {
		// get current user from os/user. Like this we can locally identify if the user changed.
		u, err := user.Current()
		if err != nil {
			log.Fatalln(err)
			return
		}
		jwt, ok := clstor.Users[u.Username]
		if !ok {
			fmt.Println("User not authenticated.")
			return
		}

		// request with 3s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		// Add token to gRPC Request. ctx WithToKeN
		ctxWTKN := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+jwt)

		c, err := grpcclient.DialUp()
		if err != nil {
			log.Fatalln(err)
			return
		}

		// send data to server and receive all users data.
		response, err := c.SyncVault(ctxWTKN, &syncData)
		if err != nil {
			log.Println(`[ERROR]:`, err)
			fmt.Println("request failed. please try again.")
			return
		}

		// check server response
		if response.GetStatus() != "success" {
			log.Println(response.GetStatus())
			fmt.Println("request failed. please try again.")
			return
		}

		//check for latest version data
		syncVault := clserv.CombineVault(clstor.Local[u.Username], clserv.VaultSyncConvert(response))

		// save actual data
		clstor.Local[u.Username] = syncVault
		fmt.Println(response.GetStatus())
	},
}

var (
	syncData pb.SyncVaultRequest
)

func init() {
	rootCmd.AddCommand(syncVaultCmd)
}
