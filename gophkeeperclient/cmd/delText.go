package cmd

import (
	"context"
	"fmt"
	"github.com/EestiChameleon/gophkeeper/gophkeeperclient/grpcclient"
	clstor "github.com/EestiChameleon/gophkeeper/gophkeeperclient/storage"
	pb "github.com/EestiChameleon/gophkeeper/proto"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"os/user"
	"time"

	"github.com/spf13/cobra"
)

// delTextCmd represents the delText command
var delTextCmd = &cobra.Command{
	Use:   "delText",
	Short: "Delete the text data by title",
	Long: `
This command allows to the authenticated user to delete the text data.
Usage: gophkeeperclient delText --title=<title>.`,
	Run: func(cmd *cobra.Command, args []string) {
		user, err := user.Current()
		if err != nil {
			log.Fatalln(err)
		}
		jwt, ok := clstor.Users[user.Username]
		if !ok {
			fmt.Println("User not authenticated.")
			return
		}
		// search local version
		vault, ok := clstor.Local[user.Username]
		if !ok {
			fmt.Println("User not found. Please register.")
			return
		}
		_, ok = vault.Text[delText.Title]
		// local version doesn't exist: nothing to delete.
		if ok {
			msg := fmt.Sprintf("Nothing found for title: %s\nMake sure you have the latest version by synchronizing your vault.",
				delText.Title)
			fmt.Println(msg)
			return
		}
		// local version found - delete on server and then delete local version

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
		response, err := c.DelText(ctxWTKN, &delText)
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

		// check server response
		if response.GetStatus() != "success" {
			if response.GetStatus() == "not found" {
				fmt.Println("Nothing found for the requested title:", delText.Title)
				return
			}
			log.Println(response.GetStatus())
			fmt.Println("request failed. please try again.")
			return
		}

		// successful response
		// delete local version
		delete(vault.Text, delText.Title)

		fmt.Println(response.GetStatus())
	},
}

var (
	delText pb.DelTextRequest
)

func init() {
	rootCmd.AddCommand(delTextCmd)
	delTextCmd.Flags().StringVarP(&delPair.Title, "title", "t", "", "Text title to delete.")
	delTextCmd.MarkFlagRequired("title")
}
