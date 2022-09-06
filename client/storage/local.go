package storage

import (
	"encoding/json"
	"github.com/EestiChameleon/gophkeeper/client/cfg"
	"github.com/EestiChameleon/gophkeeper/models"
	"log"
	"os"
)

var (
	Local         map[string]*models.Vault
	CurrentUserID string
)

// InitLocalStorage function initializes the local storage and file where we save the user's data.
func InitLocalStorage() (map[string]*models.Vault, error) {
	localStorage := make(map[string]*models.Vault)

	// create/open file
	f, err := os.OpenFile(cfg.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer f.Close()

	// read file
	bytes, err := os.ReadFile(cfg.FileStoragePath)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// parse file data
	if len(bytes) != 0 {
		return nil, json.Unmarshal(bytes, &localStorage)
	}

	return localStorage, nil
}

// InitVault initializes a new instance of VaultData.
func InitVault() *models.Vault {
	return &models.Vault{
		Pair: make(map[string]*models.Pair),
		Text: make(map[string]*models.Text),
		Bin:  make(map[string]*models.Bin),
		Card: make(map[string]*models.Card),
	}
}

// InitNewUserLocalStorage initializes a new VaultData for the provided user id. First it checks for user id existence.
// Like this we have a map[userID]userVaultData.
func InitNewUserLocalStorage(usrID string) {
	_, ok := Local[usrID]
	if !ok {
		Local[usrID] = InitVault()
	}
}

func GetUserVaultLocalStorage(usrID string) (*models.Vault, bool) {
	data, ok := Local[usrID]
	return data, ok
}

// Shutdown method closes the storage file with saving the latest data.
func ShutdownMemory() error {
	return UpdateFile()
}

// UpdateFile method rewrite the storage file with the latest data.
func UpdateFile() error {
	// open & rewrite file
	f, err := os.OpenFile(cfg.FileStoragePath, os.O_WRONLY, 0777)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()

	// prepare data
	jsonByte, err := json.Marshal(Local)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = f.Write(jsonByte)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
