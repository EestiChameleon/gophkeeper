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
	"google.golang.org/grpc/status"
	"log"
	"os/user"
	"time"

	"github.com/spf13/cobra"
)

// getBinaryCmd represents the getBinary command
var getBinaryCmd = &cobra.Command{
	Use:   "getBinary",
	Short: "Get a binary data by title",
	Long: `
This command returns to the authenticated user the binary data requested by title.
Usage: gophkeeperclient getBinary --title=<title>.`,
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
		binData, ok := vault.Bin[getBin.Title]
		// local version exists - return it.
		if ok {
			msg := fmt.Sprintf("Title: %s\nBody: %s\nComment: %s\nMake sure you have the latest version by synchronizing your vault.",
				binData.Title, binData.Body, binData.Comment)
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
		response, err := c.GetBin(ctxWTKN, &getBin)
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
		vault.Bin[response.BinData.Title] = clserv.ProtoToModelsBin(response.BinData)
		// return pair data
		msg := fmt.Sprintf("Title: %s\nBody: %s\nComment: %s\nMake sure you have the latest version by synchronizing your vault.",
			response.BinData.Title, response.BinData.Body, response.BinData.Comment)
		fmt.Println(response.GetStatus())
		fmt.Println(msg)
	},
}

var (
	getBin pb.GetBinRequest
)

func init() {
	rootCmd.AddCommand(getBinaryCmd)
	getBinaryCmd.Flags().StringVarP(&getBin.Title, "title", "t", "", "Text title to search for.")
	getBinaryCmd.MarkFlagRequired("title")
}
