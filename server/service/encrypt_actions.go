package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/EestiChameleon/gophkeeper/server/cfg"
	"github.com/robbert229/jwt"
	"log"
)

var (
	ErrInvalidToken = errors.New("failed to decode the provided Token")
)

// EncryptPass creates encrypted password based on md5 hash algorithm.
func EncryptPass(pass string) string {
	h := md5.New()
	h.Write([]byte(pass))
	return hex.EncodeToString(h.Sum(nil))
}

// JWTEncodeUserID creates JWT with userID encoded inside.
func JWTEncodeUserID(value interface{}) (string, error) {
	return JWTEncode("sub", value)
}

// JWTEncode creates JWT with encoded inside passed value.
func JWTEncode(key string, value interface{}) (string, error) {
	algorithm := jwt.HmacSha256(cfg.CryptoKey)

	claims := jwt.NewClaim()
	claims.Set(key, value)

	token, err := algorithm.Encode(claims)
	if err != nil {
		return ``, err
	}

	if err = algorithm.Validate(token); err != nil {
		return ``, err
	}

	return token, nil
}

// JWTDecodeUserID provides userID from JWT, if decoding is successful.
func JWTDecodeUserID(token string) (int, error) {
	value, err := JWTDecode(token, "sub")
	if err != nil {
		return -1, err
	}
	return int(value.(float64)), nil
}

// JWTDecode decodes the passed JWT and returns the interface value.
func JWTDecode(token, key string) (interface{}, error) {
	algorithm := jwt.HmacSha256(cfg.CryptoKey)

	if err := algorithm.Validate(token); err != nil {
		log.Println(err)
		return nil, ErrInvalidToken
	}

	claims, err := algorithm.Decode(token)
	if err != nil {
		log.Println(err)
		return nil, ErrInvalidToken
	}

	return claims.Get(key)
}
