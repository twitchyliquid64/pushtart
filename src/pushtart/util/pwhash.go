package util

import (
	"golang.org/x/crypto/bcrypt"
)

const salt = "dsDF*&(634)"

//HashPassword returns a hash of the given username+password+salt.
func HashPassword(username, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(salt+username+password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

//ComparePassHash returns nil if the given hash matches a given username and password. Otherwise, the returned error
//matches the semantics of bcrypt.CompareHashAndPassword.
func ComparePassHash(storedHash, username, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(salt+username+password))
}
