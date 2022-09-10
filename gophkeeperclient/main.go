/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"fmt"
	"github.com/EestiChameleon/gophkeeper/gophkeeperclient/cmd"
	"github.com/EestiChameleon/gophkeeper/gophkeeperclient/grpcclient"
	"github.com/EestiChameleon/gophkeeper/gophkeeperclient/storage"
	"log"
)

func main() {
	fmt.Println("Start service")
	if err := storage.InitStorage(); err != nil {
		log.Fatalln(err)
	}

	fmt.Print("Command status: ")
	cmd.Execute()

	fmt.Println("Update service data")
	if err := storage.UpdateFiles(); err != nil {
		log.Fatalln(err)
	}
	log.Println("close server connection")
	if grpcclient.ActiveConnection() {
		if err := grpcclient.ConnDown(); err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println("Stop service: successful")
}
