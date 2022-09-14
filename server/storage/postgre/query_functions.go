package postgre

import (
	"github.com/EestiChameleon/gophkeeper/models"
	"log"
)

type PostgreVault struct{}

// UserAdd inserts new user in database.
func (p *PostgreVault) UserAdd(login, pass string) (int, error) {
	var usrID int
	err := GetSingleValue(
		"INSERT INTO gophkeeper_users (login, password) VALUES ($1, $2) RETURNING id;",
		&usrID, login, pass)
	if err != nil {
		return -1, err
	}

	return usrID, nil
}

func (p *PostgreVault) UserLogin(log string) (*models.User, error) {
	u := new(models.User)
	if err := GetOneRow("SELECT id, login, password FROM gophkeeper_users WHERE login = $1;",
		u, log); err != nil {
		return nil, err
	}

	return u, nil
}

// PairByTitle provides pair data found in database by title and user id.
func (p *PostgreVault) PairByTitle(title string, usrID int) (*models.Pair, error) {
	data := new(models.Pair)
	err := GetOneRow(
		"SELECT id, user_id, title, login, pass, comment, version, deleted_at FROM gk_pair "+
			"WHERE title = $1 AND user_id = $2 AND deleted_at isnull ORDER BY version DESC LIMIT 1;",
		data, title, usrID)

	return data, err
}

// PairAdd inserts new pair data in database.
func (p *PostgreVault) PairAdd(uID int, title, login, pass, comment string, v uint32) error {
	var resultID int
	return GetSingleValue("INSERT INTO gk_pair (user_id, title, login, pass, comment, version) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;",
		&resultID,
		uID, title, login, pass, comment, v)
}

// PairDelete makes a soft delete of a pair data from database. Set deleted_at parameter to current_date.
func (p *PostgreVault) PairDelete(title string, uID int) error {
	affRows, err := ExecuteQuery(
		"UPDATE gk_pair SET deleted_at = current_timestamp WHERE title = $1 AND user_id = $2;",
		title, uID)
	log.Println("PairDelete affected rows:", affRows)
	return err
}

// TextByTitle provides text data found in database by title and user id.
func (p *PostgreVault) TextByTitle(title string, usrID int) (*models.Text, error) {
	data := new(models.Text)
	err := GetOneRow(
		"SELECT id, user_id, title, body, comment, version, deleted_at FROM gk_text "+
			"WHERE title = $1 AND user_id = $2 AND deleted_at isnull ORDER BY version DESC LIMIT 1;",
		data, title, usrID)

	return data, err
}

// TextAdd inserts new text data in database.
func (p *PostgreVault) TextAdd(uID int, title, body, comment string, v uint32) error {
	var resultID int
	return GetSingleValue(
		"INSERT INTO gk_text (user_id, title, body, comment, version) "+
			"VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		&resultID,
		uID, title, body, comment, v)
}

// TextDelete makes a soft delete of a text data from database. Set deleted_at parameter to current_date.
func (p *PostgreVault) TextDelete(title string, uID int) error {
	affRows, err := ExecuteQuery(
		"UPDATE gk_text SET deleted_at = current_timestamp "+
			"WHERE title = $1 AND user_id = $2;",
		title, uID)
	log.Println("TextDelete affected rows:", affRows)
	return err
}

// BinByTitle provides binary data found in database by title and user id.
func (p *PostgreVault) BinByTitle(title string, usrID int) (*models.Bin, error) {
	data := new(models.Bin)
	err := GetOneRow(
		"SELECT id, user_id, title, body, comment, version, deleted_at FROM gk_bin "+
			"WHERE title = $1 AND user_id = $2 AND deleted_at isnull ORDER BY version DESC LIMIT 1;",
		data, title, usrID)

	return data, err
}

// BinAdd inserts new binary data in database.
func (p *PostgreVault) BinAdd(uID int, title string, body []byte, comment string, v uint32) error {
	var resultID int
	return GetSingleValue("INSERT INTO gk_bin (user_id, title, body, comment, version) "+
		"VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		&resultID,
		uID, title, body, comment, v)
}

// BinDelete makes a soft delete of a binary data from database. Set deleted_at parameter to current_date.
func (p *PostgreVault) BinDelete(title string, uID int) error {
	affRows, err := ExecuteQuery(
		"UPDATE gk_bin SET deleted_at = current_timestamp WHERE title = $1 AND user_id = $2;",
		title, uID)
	log.Println("BinDelete affected rows:", affRows)
	return err
}

// CardByTitle provides card data found in database by title and user id.
func (p *PostgreVault) CardByTitle(title string, usrID int) (*models.Card, error) {
	data := new(models.Card)
	err := GetOneRow(
		"SELECT id, user_id, title, number, expiration_date, comment, version, deleted_at FROM gk_card "+
			"WHERE title = $1 AND user_id = $2 AND deleted_at isnull ORDER BY version DESC LIMIT 1;",
		data, title, usrID)

	return data, err
}

// CardAdd inserts new card data in database.
func (p *PostgreVault) CardAdd(uID int, title, number, expdate, comment string, v uint32) error {
	var resultID int
	return GetSingleValue(
		"INSERT INTO gk_card (user_id, title, number, expiration_date, comment, version) "+
			"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;",
		&resultID,
		uID, title, number, expdate, comment, v)
}

// CardDelete makes a soft delete of a card data from database. Set deleted_at parameter to current_date.
func (p *PostgreVault) CardDelete(title string, uID int) error {
	affRows, err := ExecuteQuery(
		"UPDATE gk_card SET deleted_at = current_timestamp WHERE title = $1 AND user_id = $2;",
		title, uID)
	log.Println("CardDelete affected rows:", affRows)
	return err
}
