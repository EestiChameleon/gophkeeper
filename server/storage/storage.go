package storage

import (
	"github.com/EestiChameleon/gophkeeper/models"
	"github.com/EestiChameleon/gophkeeper/server/storage/postgre"
	"github.com/EestiChameleon/gophkeeper/server/storage/testdb"
)

var (
	Vault Vaulter
)

type Vaulter interface {
	UserAdd(login, pass string) (int, error)
	UserLogin(log string) (*models.User, error)
	PairInt
	TextInt
	BinInt
	CardInt
}

type PairInt interface {
	PairByTitle(title string, usrID int) (*models.Pair, error)
	PairAdd(uID int, title, login, pass, comment string, v uint32) error
	PairDelete(title string, uID int) error
}

type TextInt interface {
	TextByTitle(title string, usrID int) (*models.Text, error)
	TextAdd(uID int, title, body, comment string, v uint32) error
	TextDelete(title string, uID int) error
}

type BinInt interface {
	BinByTitle(title string, usrID int) (*models.Bin, error)
	BinAdd(uID int, title string, body []byte, comment string, v uint32) error
	BinDelete(title string, uID int) error
}

type CardInt interface {
	CardByTitle(title string, usrID int) (*models.Card, error)
	CardAdd(uID int, title, number, expdate, comment string, v uint32) error
	CardDelete(title string, uID int) error
}

// Init initializes the DB connection.
func Init() (err error) {
	Vault, err = postgre.Run()
	if err != nil {
		return err
	}

	return nil
}

func Close() error {
	return postgre.ShutDown()
}

// InitTest initializes the test DB for tests.
func InitTest() (err error) {
	Vault, err = testdb.Run()
	return err
}
