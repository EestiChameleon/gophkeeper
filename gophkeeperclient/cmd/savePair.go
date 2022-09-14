/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"github.com/EestiChameleon/gophkeeper/gophkeeperclient/grpcclient"
	clstor "github.com/EestiChameleon/gophkeeper/gophkeeperclient/storage"
	"github.com/EestiChameleon/gophkeeper/models"
	pb "github.com/EestiChameleon/gophkeeper/proto"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"os/user"
	"time"

	"github.com/spf13/cobra"
)

// savePairCmd represents the savePair command
var savePairCmd = &cobra.Command{
	Use:   "savePair",
	Short: "Save a new pair of login&password",
	Long: `
This command allows to the authenticated user to save new pair data.
Usage: gophkeeperclient savePair --title=<title_for_saved_login&password> --login=<login_to_save> --password=<password_to_save> --comment=<comment_for_saved_login&password>.`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for user auth
		user, err := user.Current()
		if err != nil {
			log.Fatalln(err)
		}
		jwt, ok := clstor.Users[user.Username]
		if !ok {
			fmt.Println("User not authenticated.")
			return
		}
		// search user local vault
		vault, ok := clstor.Local[user.Username]
		if !ok {
			fmt.Println("User not found. Please register.")
			return
		}
		// search for local version
		pair, ok := vault.Pair[savePair.Title]
		// local version exists - return it.
		if ok {
			// we save new version - so we take current version + 1
			savePair.Version = pair.Version + 1
		} else {
			// not found - version = 1 - first new
			savePair.Version = 1
		}

		// request with 3s timeout. ctx WithTimeOut
		ctxWTO, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		c, err := grpcclient.DialUp()
		if err != nil {
			log.Fatalln(err)
			return
		}

		// Add token to gRPC Request. ctx WithToKeN
		ctxWTKN := metadata.AppendToOutgoingContext(ctxWTO, "authorization", "Bearer "+jwt)

		// send data to server and receive JWT in case of success. then save it in Users
		response, err := c.PostPair(ctxWTKN, &pb.PostPairRequest{Pair: &savePair})
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				// Error was not a status error
				fmt.Println("request failed. please try again.")
			}
			msg := fmt.Sprintf("Request failed.\nStatusCode: %v\nMessage: %s", st.Code(), st.Message())
			fmt.Println(msg)
			return
		}

		// successful response
		// save pair to local
		vault.Pair[savePair.Title] = models.ProtoToModelsPair(&savePair)
		// return pair data
		fmt.Println(response.GetStatus())

	},
}

var (
	savePair pb.Pair
)

func init() {
	rootCmd.AddCommand(savePairCmd)
	savePairCmd.Flags().StringVarP(&savePair.Title, "title", "t", "", "Pair title to save.")
	savePairCmd.Flags().StringVarP(&savePair.Login, "login", "l", "", "Login to save.")
	savePairCmd.Flags().StringVarP(&savePair.Pass, "password", "p", "", "Password to save.")
	savePairCmd.Flags().StringVarP(&savePair.Comment, "comment", "c", "", "Comment for the saved pair. Optional.")
	savePairCmd.MarkFlagRequired("title")
	savePairCmd.MarkFlagRequired("login")
	savePairCmd.MarkFlagRequired("password")
}
