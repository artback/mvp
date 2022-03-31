package pass

import (
	"golang.org/x/crypto/bcrypt"
)

func Compare(hashedPwd string, plainPwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)); err != nil {
		return false
	}
	return true
}
