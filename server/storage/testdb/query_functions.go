package testdb

import (
	"database/sql"
	"github.com/EestiChameleon/gophkeeper/models"
	"github.com/EestiChameleon/gophkeeper/server/storage/postgre"
	"log"
)

var (
	testUser = &models.User{
		ID:       7,
		Login:    "user7",
		Password: "8c96c3884a827355aed2c0f744594a52", //service.EncryptPass("pass7")
	}

	testPair = &models.Pair{
		ID:        1,
		UserID:    7,
		Title:     "testPair",
		Login:     "testLogin",
		Pass:      "testPass",
		Comment:   "testComment",
		Version:   7,
		DeletedAt: sql.NullTime{},
	}

	testText = &models.Text{
		ID:        2,
		UserID:    7,
		Title:     "testText",
		Body:      "testBody",
		Comment:   "testComment",
		Version:   7,
		DeletedAt: sql.NullTime{},
	}

	testBin = &models.Bin{
		ID:        3,
		UserID:    7,
		Title:     "testBin",
		Body:      []byte("testBody"),
		Comment:   "testComment",
		Version:   7,
		DeletedAt: sql.NullTime{},
	}

	testCard = &models.Card{
		ID:             4,
		UserID:         7,
		Title:          "testCard",
		Number:         "testNumber",
		ExpirationDate: "testExpDate",
		Comment:        "testComment",
		Version:        7,
		DeletedAt:      sql.NullTime{},
	}
)

type TestVault struct{}

// UserAdd imitates user creation method. Returns id = 7.
func (t *TestVault) UserAdd(login, pass string) (int, error) {
	log.Printf("Test UserAdd: login %s, password %s", login, pass)
	return 7, nil
}

// UserAdd imitates user creation method. Returns id = 7.
func (t *TestVault) UserLogin(login string) (*models.User, error) {
	log.Printf("Test UserLogin: login %s", login)
	if login != testUser.Login {
		return nil, postgre.ErrNotFound
	}

	return testUser, nil
}

// PairByTitle provides test pair data.
// All int values = 7. All string values = "test" + fieldName. Like Title = "testTitle".
func (t *TestVault) PairByTitle(title string, usrID int) (*models.Pair, error) {
	log.Printf("Test PairByTitle: title %s, user %d", title, usrID)
	if title != testPair.Title {
		return nil, postgre.ErrNotFound
	}
	return testPair, nil
}

func (t *TestVault) PairAdd(uID int, title, login, pass, comment string, v uint32) error {
	log.Printf("Test PairAdd: %v, %v, %v, %v, %v, %v", uID, title, login, pass, comment, v)
	return nil
}

func (t *TestVault) PairDelete(title string, uID int) error {
	log.Printf("Test PairDelete: %v, %v", title, uID)
	if title == testPair.Title {
		testPair.DeletedAt.Valid = true
	}
	return nil
}

func (t *TestVault) TextByTitle(title string, usrID int) (*models.Text, error) {
	if title != testText.Title {
		return nil, postgre.ErrNotFound
	}
	return testText, nil
}

func (t *TestVault) TextAdd(uID int, title, body, comment string, v uint32) error {
	log.Printf("Test TextAdd: %v, %v, %v, %v", uID, title, body, comment)
	return nil
}

func (t *TestVault) TextDelete(title string, uID int) error {
	log.Printf("Test TextDelete: %v, %v", title, uID)
	if title == testText.Title {
		testText.DeletedAt.Valid = true
	}
	return nil
}

func (t *TestVault) BinByTitle(title string, usrID int) (*models.Bin, error) {
	if title != testBin.Title {
		return nil, postgre.ErrNotFound
	}
	return testBin, nil
}

func (t *TestVault) BinAdd(uID int, title string, body []byte, comment string, v uint32) error {
	log.Printf("Test BinAdd: %v, %v, %v, %v", uID, title, body, comment)
	return nil
}

func (t *TestVault) BinDelete(title string, uID int) error {
	log.Printf("Test BinDelete: %v, %v", title, uID)
	if title == testBin.Title {
		testBin.DeletedAt.Valid = true
	}
	return nil
}

func (t *TestVault) CardByTitle(title string, usrID int) (*models.Card, error) {
	if title != testCard.Title {
		return nil, postgre.ErrNotFound
	}
	return testCard, nil
}

func (t *TestVault) CardAdd(uID int, title, number, expdate, comment string, v uint32) error {
	log.Printf("Test CardAdd: %v, %v, %v, %v, %v, %v", uID, title, number, expdate, comment, v)
	return nil
}

func (t *TestVault) CardDelete(title string, uID int) error {
	log.Printf("Test CardDelete: %v, %v", title, uID)
	if title == testCard.Title {
		testCard.DeletedAt.Valid = true
	}
	return nil
}
