package models

import (
	"database/sql"
	pb "github.com/EestiChameleon/gophkeeper/proto"
)

// User is a local struct for database interactions. Table gophkeeper_users.
type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Pair is a local struct for database interactions. Table gk_pair.
type Pair struct {
	ID        int          `json:"id"`
	UserID    int          `json:"user_id"`
	Title     string       `json:"title"`
	Login     string       `json:"login"`
	Pass      string       `json:"pass"`
	Comment   string       `json:"comment"`
	Version   uint32       `json:"version"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

// Text is a local struct for database interactions. Table gk_text.
type Text struct {
	ID        int          `json:"id"`
	UserID    int          `json:"user_id"`
	Title     string       `json:"title"`
	Body      string       `json:"body"`
	Comment   string       `json:"comment"`
	Version   uint32       `json:"version"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

// Bin is a local struct for database interactions. Table gk_bin.
type Bin struct {
	ID        int          `json:"id"`
	UserID    int          `json:"user_id"`
	Title     string       `json:"title"`
	Body      []byte       `json:"body"`
	Comment   string       `json:"comment"`
	Version   uint32       `json:"version"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

// Card is a local struct for database interactions. Table gk_card.
type Card struct {
	ID             int          `json:"id"`
	UserID         int          `json:"user_id"`
	Title          string       `json:"title"`
	Number         string       `json:"number"`
	ExpirationDate string       `json:"expiration_date"`
	Comment        string       `json:"comment"`
	Version        uint32       `json:"version"`
	DeletedAt      sql.NullTime `json:"deleted_at"`
}

// Vault is a local struct for client interactions. Mostly for easy and fast search.
type Vault struct {
	Pair map[string]*Pair `json:"pair"`
	Text map[string]*Text `json:"text"`
	Bin  map[string]*Bin  `json:"bin"`
	Card map[string]*Card `json:"card"`
}

// ActualData is a local struct for database interactions. Unites all data.
type ActualData struct {
	Pairs []*Pair `json:"pairs"`
	Texts []*Text `json:"texts"`
	Bins  []*Bin  `json:"bins"`
	Cards []*Card `json:"cards"`
}

// VaultProto is a structure for client local data operations.
type VaultProto struct {
	Pair map[string]*pb.Pair
	Text map[string]*pb.Text
	Bin  map[string]*pb.Bin
	Card map[string]*pb.Card
}

// ActualProtoData is a structure for gRPC requests operations.
type ActualProtoData struct {
	Pairs []*pb.Pair
	Texts []*pb.Text
	Bins  []*pb.Bin
	Cards []*pb.Card
}
