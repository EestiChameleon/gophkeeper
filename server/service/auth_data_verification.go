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

func CheckAuthData(ld LoginData) (string, error) {
	u := new(models.User)
	if err := storage.GetOneRow("user_by_login", u, ld.Login); err != nil {
		return "", err
	}

	if EncryptPass(ld.Password) != u.Password {
		return "", ErrWrongAuthData
	}

	return JWTEncodeUserID(u.ID)
}
