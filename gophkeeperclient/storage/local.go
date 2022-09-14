package storage

import (
	"encoding/json"
	"github.com/EestiChameleon/gophkeeper/gophkeeperclient/cfg"
	"github.com/EestiChameleon/gophkeeper/models"
	pb "github.com/EestiChameleon/gophkeeper/proto"
	"io/ioutil"
	"log"
	"os"
)

var (
	Users = make(map[string]string)     //UserLocalName: JWT from server. UserLocalName is obtained via os/user -> user.Current()
	Local map[string]*models.VaultProto // user id :vault
)

/*
// to think about?
type LocalStorer interface {
	Save(string, []byte)        // storageName string, dataJSON []byte => Save("pair", [01010101])
	Get(string, string) ClientData // Get("pair", "title") => 1-where(pair, text, bin, card), 2 - what
	Delete(string, string)      // Delete("pair", "title") => 1-where(pair, text, bin, card), 2 - what
}

// ClientData is used to simplify the data operations.
type ClientData struct {
	Title   string //data title
	Data    []byte //login&pass / text / binary / cardNumber&expDate after marshal
	Comment string //data comment
}
*/

// InitStorage function initializes the storage data (check files & parse to local memory).
func InitStorage() (err error) {
	if err = initUsers(); err != nil {
		return err
	}

	if err = initLocal(); err != nil {
		return err
	}

	return nil
}

// initUsers reads or creates the local user auth info file. Then parse the content to local memory.
func initUsers() error {
	// create/open file
	fu, err := os.OpenFile(cfg.UsersFileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Println(err)
		return err
	}
	defer fu.Close()

	// read file
	ubytes, err := os.ReadFile(cfg.UsersFileStoragePath)
	if err != nil {
		log.Println(err)
		return err
	}

	// parse file data
	if len(ubytes) != 0 {
		return json.Unmarshal(ubytes, &Users)
	}

	return nil
}

// initLocal reads or creates the local users data storage file. Then parse the content to local memory.
func initLocal() error {
	// second - open vault file
	Local = make(map[string]*models.VaultProto)

	// create/open file
	fv, err := os.OpenFile(cfg.VaultFileStoragePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Println(err)
		return err
	}
	defer fv.Close()

	// read the whole file at once
	vbytes, err := ioutil.ReadFile(cfg.VaultFileStoragePath)
	if err != nil {
		panic(err)
	}

	// parse file data
	if len(vbytes) != 0 {
		return json.Unmarshal(vbytes, &Local)
	}

	return nil
}

// MakeVaultProto initializes a new instance of VaultProto.
func MakeVaultProto() *models.VaultProto {
	return &models.VaultProto{
		Pair: make(map[string]*pb.Pair),
		Text: make(map[string]*pb.Text),
		Bin:  make(map[string]*pb.Bin),
		Card: make(map[string]*pb.Card),
	}
}

// UpdateFiles rewrites local files with actual data.
func UpdateFiles() error {
	// prepare users data
	usersJSONByte, err := json.Marshal(Users)
	if err != nil {
		log.Println(err)
		return err
	}
	if err = UpdateFile(cfg.UsersFileStoragePath, usersJSONByte); err != nil {
		return err
	}

	// prepare vault data
	vaultJSONByte, err := json.Marshal(Local)
	if err != nil {
		log.Println(err)
		return err
	}

	if err = UpdateFile(cfg.VaultFileStoragePath, vaultJSONByte); err != nil {
		return err
	}

	return nil
}

// UpdateFile method rewrite the file with the passed data.
func UpdateFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0644)
}
