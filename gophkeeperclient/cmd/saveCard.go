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

// saveCardCmd represents the saveCard command
var saveCardCmd = &cobra.Command{
	Use:   "saveCard",
	Short: "Save a new card data",
	Long: `
This command allows to the authenticated user to save new card data.
Usage: gophkeeperclient saveCard --title=<title_for_saved_card> --number=<card_number_to_save> --expdate=<card_expiration_date> --comment=<comment_for_saved_card>.`,
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
		card, ok := vault.Bin[saveCard.Title]
		// local version exists - return it.
		if ok {
			// we save new version - so we take current version + 1
			saveCard.Version = card.Version + 1
		}
		// not found - version = 1 - first new
		saveCard.Version = 1

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
		response, err := c.PostCard(ctxWTKN, &pb.PostCardRequest{Card: &saveCard})
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
			log.Println(response.GetStatus())
			msg := fmt.Sprintf("Request failed.\nStatus: %v\nPlease try again", response.GetStatus())
			fmt.Println(msg)
			return
		}
		// save data to local
		vault.Card[saveCard.Title] = &saveCard
		// successful response
		fmt.Println("New data successfully saved")
	},
}

var (
	saveCard pb.Card
)

func init() {
	rootCmd.AddCommand(saveCardCmd)
	saveCardCmd.Flags().StringVarP(&saveCard.Title, "title", "t", "", "Card title to save.")
	saveCardCmd.Flags().StringVarP(&saveCard.Number, "number", "n", "", "Card number to save.")
	saveCardCmd.Flags().StringVarP(&saveCard.Expdate, "expdate", "xd", "", "Card expiration date to save.")
	saveCardCmd.Flags().StringVarP(&saveCard.Comment, "comment", "c", "", "Comment for the saved card data (optional).")
	saveCardCmd.MarkFlagRequired("title")
	saveCardCmd.MarkFlagRequired("number")
	saveCardCmd.MarkFlagRequired("expdate")
}
