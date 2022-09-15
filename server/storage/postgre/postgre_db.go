package postgre

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
)

var (
	ErrNotFound            = errors.New("no records found")
	ErrRecordAlreadyExists = errors.New("provided data already exists")
	db                     *pgxpool.Pool
)

// Run method initiates the DB connection and creates the gophkeeper tables.
func Run() (*PostgreVault, error) {
	//create tables if it doesn't exist
	if err := migration.InitMigration(); err != nil {
		return nil, err
	}

	// connect to DB
	conn, err := pgxpool.Connect(context.Background(), cfg.PostgreDatabaseURI)
	if err != nil {
		return nil, err
	}

	db = conn
	//Vault.MU = new(sync.Mutex)
	return &PostgreVault{}, nil
	// or may be init all and then synchronizes ?
}

func ShutDown() error {
	db.Close()
	return nil
}

//-------------------- DATABASE QUERIES--------------------

// ExecuteQuery is used for SQL queries that returns nothing. Like DELETE or UPDATE.
func ExecuteQuery(query string, args ...interface{}) (int, error) {
	rows, err := db.Exec(context.Background(), query, args...)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return int(rows.RowsAffected()), nil
}

// GetSingleValue returns a SINGLE value (!) from sql query (it can be number of rows affected, id of the new inserted row, etc...).
func GetSingleValue(query string, dest interface{}, args ...interface{}) (err error) {
	if err = db.QueryRow(context.Background(), query, args...).Scan(dest); err != nil {
		log.Println(err)
		return err
	}
	return
}

// GetOneRow returns a data ROW (1 row) from sql query.
func GetOneRow(query string, dest interface{}, args ...interface{}) (err error) {
	if err = pgxscan.Get(context.Background(), db, dest, query, args...); err != nil {
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
	if err = pgxscan.Select(context.Background(), db, dest, query, args...); err != nil {
		log.Println(err)
		return err
	}
	return
}

// getAllUserDataLastVersion returns all user's data found in database. Last version.
func getAllUserDataLastVersion(usrID int) (*models.ActualData, error) {
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

	return data, nil
}
