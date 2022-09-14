package service

import (
	"errors"
	"github.com/EestiChameleon/gophkeeper/models"
	"github.com/EestiChameleon/gophkeeper/server/storage"
)

var (
	ErrWrongAuthData = errors.New("wrong authentication data")
)

// LoginData - структура данных логин/пароль пользователя
type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// CheckAuthData verifies the provided login&password values.
// If user found with such login and password - return JWT with encoded userID.
func CheckAuthData(ld LoginData) (string, error) {
	u := new(models.User)
	if err := storage.GetOneRow("SELECT id, login, password FROM gophkeeper_users WHERE login = $1;",
		u, ld.Login); err != nil {
		return "", err
	}

	if EncryptPass(ld.Password) != u.Password {
		return "", ErrWrongAuthData
	}

	return JWTEncodeUserID(u.ID)
}
