package testdb

import (
	"database/sql"
	"github.com/EestiChameleon/gophkeeper/models"
	pb "github.com/EestiChameleon/gophkeeper/proto"
	"github.com/EestiChameleon/gophkeeper/server/storage/postgre"
	"log"
)

var (
	TestUser = &models.User{
		ID:       7,
		Login:    "user7",
		Password: "8c96c3884a827355aed2c0f744594a52", //service.EncryptPass("pass7")
	}

	TestPair = &models.Pair{
		ID:        1,
		UserID:    7,
		Title:     "testPair",
		Login:     "testLogin",
		Pass:      "testPass",
		Comment:   "testComment",
		Version:   7,
		DeletedAt: sql.NullTime{},
	}

	TestText = &models.Text{
		ID:        2,
		UserID:    7,
		Title:     "testText",
		Body:      "testBody",
		Comment:   "testComment",
		Version:   7,
		DeletedAt: sql.NullTime{},
	}

	TestBin = &models.Bin{
		ID:        3,
		UserID:    7,
		Title:     "testBin",
		Body:      []byte("testBody"),
		Comment:   "testComment",
		Version:   7,
		DeletedAt: sql.NullTime{},
	}

	TestCard = &models.Card{
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
	log.Printf("Test UserAdd: login %s, encrypted password %s", login, pass)
	return 7, nil
}

// UserAdd imitates user creation method. Returns id = 7.
func (t *TestVault) UserLogin(login string) (*models.User, error) {
	log.Printf("Test UserLogin: login %s", login)
	if login != TestUser.Login {
		return nil, postgre.ErrNotFound
	}

	return TestUser, nil
}

// PairByTitle provides test pair data.
// All int values = 7. All string values = "test" + fieldName. Like Title = "testTitle".
func (t *TestVault) PairByTitle(title string, usrID int) (*models.Pair, error) {
	log.Printf("Test PairByTitle: title %s, user %d", title, usrID)
	if title != TestPair.Title || TestPair.DeletedAt.Valid {
		return nil, postgre.ErrNotFound
	}
	return TestPair, nil
}

func (t *TestVault) PairAdd(uID int, title, login, pass, comment string, v uint32) error {
	log.Printf("Test PairAdd: %v, %v, %v, %v, %v, %v", uID, title, login, pass, comment, v)
	return nil
}

func (t *TestVault) PairDelete(title string, uID int) error {
	log.Printf("Test PairDelete: %v, %v", title, uID)
	if title == TestPair.Title {
		TestPair.DeletedAt.Valid = true
	}
	return nil
}

func (t *TestVault) TextByTitle(title string, usrID int) (*models.Text, error) {
	if title != TestText.Title || TestText.DeletedAt.Valid {
		return nil, postgre.ErrNotFound
	}
	return TestText, nil
}

func (t *TestVault) TextAdd(uID int, title, body, comment string, v uint32) error {
	log.Printf("Test TextAdd: %v, %v, %v, %v", uID, title, body, comment)
	return nil
}

func (t *TestVault) TextDelete(title string, uID int) error {
	log.Printf("Test TextDelete: %v, %v", title, uID)
	if title == TestText.Title {
		TestText.DeletedAt.Valid = true
	}
	return nil
}

func (t *TestVault) BinByTitle(title string, usrID int) (*models.Bin, error) {
	if title != TestBin.Title || TestBin.DeletedAt.Valid {
		return nil, postgre.ErrNotFound
	}
	return TestBin, nil
}

func (t *TestVault) BinAdd(uID int, title string, body []byte, comment string, v uint32) error {
	log.Printf("Test BinAdd: %v, %v, %v, %v", uID, title, body, comment)
	return nil
}

func (t *TestVault) BinDelete(title string, uID int) error {
	log.Printf("Test BinDelete: %v, %v", title, uID)
	if title == TestBin.Title {
		TestBin.DeletedAt.Valid = true
	}
	return nil
}

func (t *TestVault) CardByTitle(title string, usrID int) (*models.Card, error) {
	if title != TestCard.Title || TestCard.DeletedAt.Valid {
		return nil, postgre.ErrNotFound
	}
	return TestCard, nil
}

func (t *TestVault) CardAdd(uID int, title, number, expdate, comment string, v uint32) error {
	log.Printf("Test CardAdd: %v, %v, %v, %v, %v, %v", uID, title, number, expdate, comment, v)
	return nil
}

func (t *TestVault) CardDelete(title string, uID int) error {
	log.Printf("Test CardDelete: %v, %v", title, uID)
	if title == TestCard.Title {
		TestCard.DeletedAt.Valid = true
	}
	return nil
}

func (t *TestVault) AllUserLatestData(usrID int) (*models.ActualProtoData, error) {
	if usrID == 7 {
		return &models.ActualProtoData{
			Pairs: []*pb.Pair{models.ModelsToProtoPair(TestPair)},
			Texts: []*pb.Text{models.ModelsToProtoText(TestText)},
			Bins:  []*pb.Bin{models.ModelsToProtoBin(TestBin)},
			Cards: []*pb.Card{models.ModelsToProtoCard(TestCard)},
		}, nil
	}

	return new(models.ActualProtoData), nil
}
