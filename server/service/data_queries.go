package service

import (
	"github.com/EestiChameleon/gophkeeper/models"
	"github.com/EestiChameleon/gophkeeper/server/storage"
	"log"
)

// UserAdd inserts new user in database.
func UserAdd(login, pass string) (int, error) {
	var usrID int
	err := storage.GetSingleValue(
		"INSERT INTO gophkeeper_users (login, password) VALUES ($1, $2) RETURNING id;",
		&usrID, login, EncryptPass(pass))
	if err != nil {
		return -1, err
	}

	return usrID, nil
}

// PairByTitle provides pair data found in database by title and user id.
func PairByTitle(title string, usrID int) (*models.Pair, error) {
	data := new(models.Pair)
	err := storage.GetOneRow(
		"SELECT id, user_id, title, login, pass, comment, version, deleted_at FROM gk_pair "+
			"WHERE title = $1 AND user_id = $2 AND deleted_at isnull ORDER BY version DESC LIMIT 1;",
		data, title, usrID)

	return data, err
}

// PairAdd inserts new pair data in database.
func PairAdd(uID int, title, login, pass, comment string, v uint32) error {
	var resultID int
	return storage.GetSingleValue("INSERT INTO gk_pair (user_id, title, login, pass, comment, version) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;",
		&resultID,
		uID, title, login, pass, comment, v)
}

// PairDelete makes a soft delete of a pair data from database. Set deleted_at parameter to current_date.
func PairDelete(title string, uID int) error {
	affRows, err := storage.ExecuteQuery(
		"UPDATE gk_pair SET deleted_at = current_timestamp WHERE title = $1 AND user_id = $2;",
		title, uID)
	log.Println("PairDelete affected rows:", affRows)
	return err
}

// TextByTitle provides text data found in database by title and user id.
func TextByTitle(title string, usrID int) (*models.Text, error) {
	data := new(models.Text)
	err := storage.GetOneRow(
		"SELECT id, user_id, title, body, comment, version, deleted_at FROM gk_text "+
			"WHERE title = $1 AND user_id = $2 AND deleted_at isnull ORDER BY version DESC LIMIT 1;",
		data, title, usrID)

	return data, err
}

// TextAdd inserts new text data in database.
func TextAdd(uID int, title, body, comment string, v uint32) error {
	var resultID int
	return storage.GetSingleValue(
		"INSERT INTO gk_text (user_id, title, body, comment, version) "+
			"VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		&resultID,
		uID, title, body, comment, v)
}

// TextDelete makes a soft delete of a text data from database. Set deleted_at parameter to current_date.
func TextDelete(title string, uID int) error {
	affRows, err := storage.ExecuteQuery(
		"UPDATE gk_text SET deleted_at = current_timestamp "+
			"WHERE title = $1 AND user_id = $2;",
		title, uID)
	log.Println("TextDelete affected rows:", affRows)
	return err
}

// BinByTitle provides binary data found in database by title and user id.
func BinByTitle(title string, usrID int) (*models.Bin, error) {
	data := new(models.Bin)
	err := storage.GetOneRow(
		"SELECT id, user_id, title, body, comment, version, deleted_at FROM gk_bin "+
			"WHERE title = $1 AND user_id = $2 AND deleted_at isnull ORDER BY version DESC LIMIT 1;",
		data, title, usrID)

	return data, err
}

// BinAdd inserts new binary data in database.
func BinAdd(uID int, title string, body []byte, comment string, v uint32) error {
	var resultID int
	return storage.GetSingleValue("INSERT INTO gk_bin (user_id, title, body, comment, version) "+
		"VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		&resultID,
		uID, title, body, comment, v)
}

// BinDelete makes a soft delete of a binary data from database. Set deleted_at parameter to current_date.
func BinDelete(title string, uID int) error {
	affRows, err := storage.ExecuteQuery(
		"UPDATE gk_bin SET deleted_at = current_timestamp WHERE title = $1 AND user_id = $2;",
		title, uID)
	log.Println("BinDelete affected rows:", affRows)
	return err
}

// CardByTitle provides card data found in database by title and user id.
func CardByTitle(title string, usrID int) (*models.Card, error) {
	data := new(models.Card)
	err := storage.GetOneRow(
		"SELECT id, user_id, title, number, expiration_date, comment, version, deleted_at FROM gk_card "+
			"WHERE title = $1 AND user_id = $2 AND deleted_at isnull ORDER BY version DESC LIMIT 1;",
		data, title, usrID)

	return data, err
}

// CardAdd inserts new card data in database.
func CardAdd(uID int, title, number, expdate, comment string, v uint32) error {
	var resultID int
	return storage.GetSingleValue(
		"INSERT INTO gk_card (user_id, title, number, expiration_date, comment, version) "+
			"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;",
		&resultID,
		uID, title, number, expdate, comment, v)
}

// CardDelete makes a soft delete of a card data from database. Set deleted_at parameter to current_date.
func CardDelete(title string, uID int) error {
	affRows, err := storage.ExecuteQuery(
		"UPDATE gk_card SET deleted_at = current_timestamp WHERE title = $1 AND user_id = $2;",
		title, uID)
	log.Println("CardDelete affected rows:", affRows)
	return err
}
