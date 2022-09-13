package storage

import (
	"errors"
	"github.com/EestiChameleon/gophkeeper/models"
	pb "github.com/EestiChameleon/gophkeeper/proto"
	"github.com/EestiChameleon/gophkeeper/server/cfg"
	migration "github.com/EestiChameleon/gophkeeper/server/migrations"
	"github.com/docker/distribution/context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"sync"
)

var (
	ErrNotFound            = errors.New("no records found")
	ErrRecordAlreadyExists = errors.New("provided data already exists")
	ctxStorage             = context.Background()
	Vault                  VaultStorage
)

type VaultStorage struct {
	DB     *pgxpool.Pool
	UserID string
	MU     *sync.Mutex
}

// Init method initializes a LocalMemory/DB storage.
func Init() error {
	conn, err := InitDBStorage()
	if err != nil {
		return err
	}

	Vault.DB = conn
	//Vault.MU = new(sync.Mutex)
	return nil
	// or may be init all and then synchronizes ?
}

//-------------------- DATABASE --------------------

// InitDBStorage initiates the DB connection and creates the shorten_pairs table.
// If there is an old shorten_pairs table, it's dropped.
func InitDBStorage() (conn *pgxpool.Pool, err error) {
	//create tables if it doesn't exist
	if err = migration.InitMigration(); err != nil {
		return nil, err
	}

	// connect to DB
	conn, err = pgxpool.Connect(context.Background(), cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// SyncVault synchronizes data between local storage and DB.
func (v *VaultStorage) SyncVault() error {
	return nil
}

// ExecuteQuery is used for SQL queries that returns nothing. Like DELETE or UPDATE.
func ExecuteQuery(query string, args ...interface{}) (int, error) {
	rows, err := Vault.DB.Exec(ctxStorage, query, args...)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return int(rows.RowsAffected()), nil
}

// GetSingleValue returns a SINGLE value (!) from sql query (it can be number of rows affected, id of the new inserted row, etc...).
func GetSingleValue(query string, dest interface{}, args ...interface{}) (err error) {
	if err = Vault.DB.QueryRow(ctxStorage, query, args...).Scan(dest); err != nil {
		log.Println(err)
		return err
	}
	return
}

// GetOneRow returns a data ROW (1 row) from sql query.
func GetOneRow(query string, dest interface{}, args ...interface{}) (err error) {
	if err = pgxscan.Get(ctxStorage, Vault.DB, dest, query, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		log.Println(err)
		return err
	}
	return
}

// GetAll returns a table with values from offset till limit params.
func GetAll(query string, dest interface{}, args ...interface{}) (err error) {
	if err = pgxscan.Select(ctxStorage, Vault.DB, dest, query, args...); err != nil {
		log.Println(err)
		return err
	}
	return
}

// GetAllUserDataLastVersion returns all user's data found in proto format.
func GetAllUserDataLastVersion(usrID int) (*models.ActualProtoData, error) {
	var err error
	data := new(models.ActualData)
	if err = GetAll("SELECT DISTINCT ON (title) title, login, pass, comment, version FROM gk_pair WHERE user_id = $1 AND deleted_at isnull ORDER BY title, version DESC;",
		&data.Pairs, usrID); err != nil {
		return nil, err
	}
	if err = GetAll(" SELECT DISTINCT ON (title) title, body, comment, version FROM gk_text WHERE user_id = $1 AND deleted_at isnull ORDER BY title, version DESC;",
		&data.Texts, usrID); err != nil {
		return nil, err
	}
	if err = GetAll("SELECT DISTINCT ON (title) title, body, comment, version FROM gk_bin WHERE user_id = $1 AND deleted_at isnull ORDER BY title, version DESC;",
		&data.Bins, usrID); err != nil {
		return nil, err
	}
	if err = GetAll("SELECT DISTINCT ON (title) title, number, expiration_date, comment, version FROM gk_card WHERE user_id = $1 AND deleted_at isnull ORDER BY title, version DESC;",
		&data.Cards, usrID); err != nil {
		return nil, err
	}

	return convertActualDataToProto(data), nil
}

// convertActualDataToProto converts local data structures, used for DB interactions, to gRPC proto structures.
func convertActualDataToProto(in *models.ActualData) *models.ActualProtoData {
	out := new(models.ActualProtoData)
	for _, v := range in.Pairs {
		out.Pairs = append(out.Pairs, convertPairToProto(v))
	}

	for _, v := range in.Texts {
		out.Texts = append(out.Texts, convertTextToProto(v))
	}

	for _, v := range in.Bins {
		out.Bins = append(out.Bins, convertBinToProto(v))
	}

	for _, v := range in.Cards {
		out.Cards = append(out.Cards, convertCardToProto(v))
	}

	return out
}

// convertPairToProto converts local Pair structure to proto Pair structure.
func convertPairToProto(in *models.Pair) *pb.Pair {
	return &pb.Pair{
		Title:   in.Title,
		Login:   in.Login,
		Pass:    in.Pass,
		Comment: in.Comment,
		Version: in.Version,
	}
}

// convertCardToProto converts local Card structure to proto Card structure.
func convertCardToProto(in *models.Card) *pb.Card {
	return &pb.Card{
		Title:   in.Title,
		Number:  in.Number,
		Expdate: in.ExpirationDate,
		Comment: in.Comment,
		Version: in.Version,
	}
}

// convertTextToProto converts local Text structure to proto Text structure.
func convertTextToProto(in *models.Text) *pb.Text {
	return &pb.Text{
		Title:   in.Title,
		Body:    in.Body,
		Comment: in.Comment,
		Version: in.Version,
	}
}

// convertBinToProto converts local Bin structure to proto Bin structure.
func convertBinToProto(in *models.Bin) *pb.Bin {
	return &pb.Bin{
		Title:   in.Title,
		Body:    in.Body,
		Comment: in.Comment,
		Version: in.Version,
	}
}
