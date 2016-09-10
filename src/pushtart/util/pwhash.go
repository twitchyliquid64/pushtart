package util

import (
  "golang.org/x/crypto/bcrypt"
)

const salt = "dsDF*&(634)"

func HashPassword(username, password string)(string, error){
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(salt+username+password), bcrypt.DefaultCost)
  if err != nil {
      return "", err
  }
  return string(hashedPassword), nil
}


func ComparePassHash(storedHash, username, password string)error{
  return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(salt+username+password))
}
