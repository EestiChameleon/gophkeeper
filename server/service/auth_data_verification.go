package service

import (
	"errors"
	"github.com/EestiChameleon/gophkeeper/models"
	"github.com/EestiChameleon/gophkeeper/server/ctxfunc"
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

func CheckAuthData(ld LoginData) (string, error) {
	u := new(models.User)
	if err := storage.GetOneRow("SELECT id, login, password FROM gophkeeper_users WHERE login = $1;",
		u, ld.Login); err != nil {
		return "", err
	}

	if EncryptPass(ld.Password) != u.Password {
		return "", ErrWrongAuthData
	}

	ctxfunc.SetUserID(u.ID)

	return JWTEncodeUserID(u.ID)
}
