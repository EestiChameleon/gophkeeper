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

// getCardCmd represents the getCard command
var getCardCmd = &cobra.Command{
	Use:   "getCard",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		card, ok := vault.Card[getCard.Title]
		// local version exists - return it.
		if ok {
			msg := fmt.Sprintf("Title: %s\nNumber: %s\nExpdate: %s\nComment: %s\nMake sure you have the latest version by synchronizing your vault.",
				card.Title, card.Number, card.Expdate, card.Comment)
			fmt.Println(msg)
			return
		}
		// local version not found - search on server

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
		response, err := c.GetCard(ctxWTKN, &getCard)
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				// Error was not a status error
				fmt.Println("request failed. please try again.")
			}
			msg := fmt.Sprintf("success\nStatusCode: %v\nMessage: %s", st.Code(), st.Message())
			fmt.Println(msg)
			return
		}

		// successful response
		// save pair to local
		vault.Card[response.Card.Title] = response.Card
		// return pair data
		msg := fmt.Sprintf("Title: %s\nNumber: %s\nExpiration date: %s\nComment: %s\nMake sure you have the latest version by synchronizing your vault.",
			response.Card.Title, response.Card.Number, response.Card.Expdate, response.Card.Comment)
		fmt.Println(response.GetStatus())
		fmt.Println(msg)
	},
}

var (
	getCard pb.GetCardRequest
)

func init() {
	rootCmd.AddCommand(getCardCmd)
	getCardCmd.Flags().StringVarP(&getCard.Title, "title", "t", "", "Card title to search for.")
	getCardCmd.MarkFlagRequired("title")

}
