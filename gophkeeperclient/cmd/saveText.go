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

// saveTextCmd represents the saveText command
var saveTextCmd = &cobra.Command{
	Use:   "saveText",
	Short: "Save a new text data",
	Long: `
This command allows to the authenticated user to save new text data.
Usage: gophkeeperclient saveText --title=<title_for_saved_text> --body=<text_content_to_save> --comment=<comment_for_saved_text>.`,
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
		text, ok := vault.Text[saveText.Title]
		// local version exists - return it.
		if ok {
			// we save new version - so we take current version + 1
			saveText.Version = text.Version + 1
		} else {
			// not found - version = 1 - first new
			saveText.Version = 1
		}

		log.Println(saveText.Version)

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
		response, err := c.PostText(ctxWTKN, &pb.PostTextRequest{Text: &saveText})
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
		vault.Text[saveText.Title] = &saveText
		// return pair data
		fmt.Println(response.GetStatus())

	},
}

var (
	saveText pb.Text
)

func init() {
	rootCmd.AddCommand(saveTextCmd)
	saveTextCmd.Flags().StringVarP(&saveText.Title, "title", "t", "", "Text title to save.")
	saveTextCmd.Flags().StringVarP(&saveText.Body, "body", "b", "", "Text to save.")
	saveTextCmd.Flags().StringVarP(&saveText.Comment, "comment", "c", "", "Comment for the saved text (optional).")
	saveTextCmd.MarkFlagRequired("title")
	saveTextCmd.MarkFlagRequired("body")
}
