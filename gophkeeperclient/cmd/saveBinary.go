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

// saveBinaryCmd represents the saveBinary command
var saveBinaryCmd = &cobra.Command{
	Use:   "saveBinary",
	Short: "Save a new binary data",
	Long: `
This command allows to the authenticated user to save new binary data.
Usage: gophkeeperclient saveBinary --title=<title_for_saved_data> --body=<binary_data> --comment=<comment_for_saved_data>.`,
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
		bin, ok := vault.Bin[saveBin.Title]
		// local version exists - return it.
		if ok {
			// we save new version - so we take current version + 1
			saveBin.Version = bin.Version + 1
		} else {
			// not found - version = 1 - first new
			saveBin.Version = 1
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
		response, err := c.PostBin(ctxWTKN, &pb.PostBinRequest{BinData: &saveBin})
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

		// save data to local
		vault.Bin[saveBin.Title] = models.ProtoToModelsBin(&saveBin)
		// successful response
		fmt.Println(response.GetStatus())
	},
}

var (
	saveBin pb.Bin
)

func init() {
	rootCmd.AddCommand(saveBinaryCmd)
	saveBinaryCmd.Flags().StringVarP(&saveBin.Title, "title", "t", "", "Binary data title to save.")
	saveBinaryCmd.Flags().BytesBase64VarP(&saveBin.Body, "body", "b", nil, "Binary data to save.")
	saveBinaryCmd.Flags().StringVarP(&saveBin.Comment, "comment", "c", "", "Comment for the saved binary data (optional).")
	saveBinaryCmd.MarkFlagRequired("title")
	saveBinaryCmd.MarkFlagRequired("body")
}
