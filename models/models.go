package models

import (
	"database/sql"
	pb "github.com/EestiChameleon/gophkeeper/proto"
)

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Pair struct {
	ID        int          `json:"id"`
	UserID    int          `json:"user_id"`
	Title     string       `json:"pair_title"`
	Login     string       `json:"pair_login"`
	Pass      string       `json:"pair_pass"`
	Comment   string       `json:"pair_comment"`
	Version   uint32       `json:"pair_version"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type Text struct {
	ID        int          `json:"id"`
	UserID    int          `json:"user_id"`
	Title     string       `json:"text_title"`
	Body      string       `json:"text_body"`
	Comment   string       `json:"text_comment"`
	Version   uint32       `json:"text_version"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type Bin struct {
	ID        int          `json:"id"`
	UserID    int          `json:"user_id"`
	Title     string       `json:"bin_title"`
	Body      []byte       `json:"bin_body"`
	Comment   string       `json:"bin_comment"`
	Version   uint32       `json:"bin_version"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type Card struct {
	ID        int          `json:"id"`
	UserID    int          `json:"user_id"`
	Title     string       `json:"card_title"`
	Number    string       `json:"card_number"`
	ExpDate   string       `json:"expdate"`
	Comment   string       `json:"card_comment"`
	Version   uint32       `json:"card_version"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type Vault struct {
	Pair map[string]*Pair `json:"pair"`
	Text map[string]*Text `json:"text"`
	Bin  map[string]*Bin  `json:"bin"`
	Card map[string]*Card `json:"card"`
}

type ActualData struct {
	Pairs []*Pair `json:"pairs"`
	Texts []*Text `json:"texts"`
	Bins  []*Bin  `json:"bins"`
	Cards []*Card `json:"cards"`
}

type VaultProto struct {
	Pair map[string]*pb.Pair
	Text map[string]*pb.Text
	Bin  map[string]*pb.Bin
	Card map[string]*pb.Card
}

type ActualProtoData struct {
	Pairs []*pb.Pair
	Texts []*pb.Text
	Bins  []*pb.Bin
	Cards []*pb.Card
}
