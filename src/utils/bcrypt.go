package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func CreateHash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func MustHashed(hashed string, err error) (result string) {
	if err != nil {
		panic(err.Error())
	}
	return hashed
}

func ValidateHash(password, hashed string) (result bool) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
