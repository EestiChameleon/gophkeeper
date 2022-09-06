package storage

import (
	"errors"
	"github.com/EestiChameleon/gophkeeper/models"
	"github.com/EestiChameleon/gophkeeper/server/cfg"
	migration "github.com/EestiChameleon/gophkeeper/server/migrations"
	"github.com/docker/distribution/context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"strconv"
	"sync"
)

var (
	ErrNotFound            = errors.New("no records found")
	ErrRecordAlreadyExists = errors.New("provided data already exists")
	ctxStorage             = context.Background()
	Vault                  *VaultStorage
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

// GetSingleValue returns a SINGLE value (!) from sql query (it can be number of rows affected, id of the new inserted row, etc...)
func GetSingleValue(funcName string, dest interface{}, args ...interface{}) (err error) {
	sqlQuery := "SELECT * FROM " + funcName + "("
	for i := range args {
		if i > 0 {
			sqlQuery += ", "
		}
		sqlQuery += "$" + strconv.Itoa(i+1)
	}
	sqlQuery += ");"
	if err = Vault.DB.QueryRow(ctxStorage, sqlQuery, args...).Scan(dest); err != nil {
		log.Println(err)
		return err
	}
	return
}

// GetOneRow returns a data ROW (1 row) from sql query
func GetOneRow(funcName string, dest interface{}, args ...interface{}) (err error) {
	sqlQuery := "SELECT * FROM " + funcName + "("
	for i := range args {
		if i > 0 {
			sqlQuery += ", "
		}
		sqlQuery += "$" + strconv.Itoa(i+1)
	}
	sqlQuery += ")"

	if err = pgxscan.Get(ctxStorage, Vault.DB, dest, sqlQuery, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		log.Println(err)
		return err
	}
	return
}

// GetAll returns a table with values from offset till limit params
func GetAll(funcName string, dest interface{}, args ...interface{}) (err error) {
	sqlQuery := "SELECT * FROM " + funcName + "("
	for i := range args {
		if i > 0 {
			sqlQuery += ", "
		}
		sqlQuery += "$" + strconv.Itoa(i+1)
	}
	sqlQuery += ")"
	if err = pgxscan.Select(ctxStorage, Vault.DB, dest, sqlQuery, args...); err != nil {
		log.Println(err)
		return err
	}
	return
}

func GetAllUserDataLastVersion(usrID int) (*models.ActualProtoData, error) {
	var err error
	data := new(models.ActualProtoData)
	if err = GetAll("pairs_all_last_version_by_user_id", &data.Pairs, usrID); err != nil {
		return nil, err
	}
	if err = GetAll("texts_all_last_version_by_user_id", &data.Texts, usrID); err != nil {
		return nil, err
	}
	if err = GetAll("bins_all_last_version_by_user_id", &data.Bins, usrID); err != nil {
		return nil, err
	}
	if err = GetAll("cards_all_last_version_by_user_id", &data.Cards, usrID); err != nil {
		return nil, err
	}

	return data, nil
}
