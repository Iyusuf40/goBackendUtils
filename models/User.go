package models

import (
	"github.com/Iyusuf40/goBackendUtils/config"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     int    `json:"phone"`
	Password  string `json:"password"`
}

func (user *User) HashPassword() {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), config.UserPassowrdHashCost)
	user.Password = string(hash)
}

func (user *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
